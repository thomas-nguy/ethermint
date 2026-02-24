# How to Create a New Precompile in Ethermint

This guide walks through adding a **custom precompile**: a native contract at a fixed address that runs Go code instead of EVM bytecode. Ethermint merges go-ethereum’s default precompiles with your custom ones and gives them access to Cosmos `sdk.Context` when needed.

---

## 1. Implement `vm.PrecompiledContract`

Your type must satisfy the interface used by the EVM (from `github.com/ethereum/go-ethereum/core/vm`). The exact interface may vary slightly by go-ethereum version; typical methods are:

- **`Address() common.Address`** – Contract address (required by Ethermint’s keeper when registering).
- **`RequiredGas(input []byte) uint64`** – Gas cost for the given input (charged before `Run`).
- **`Run(input []byte, contract *vm.Contract) ([]byte, error)`** – Execute the precompile; return output or error (reverts the call on error).

**Example: minimal “identity” style precompile**

```go
package myprecompile

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const myPrecompileAddressHex = "0x0000000000000000000000000000000000000100" // pick an unused address

var myPrecompileAddress = common.HexToAddress(myPrecompileAddressHex)

type MyPrecompile struct{}

func (MyPrecompile) Address() common.Address {
	return myPrecompileAddress
}

func (MyPrecompile) RequiredGas(input []byte) uint64 {
	return 1000 // or compute from len(input), etc.
}

func (MyPrecompile) Run(input []byte, contract *vm.Contract) ([]byte, error) {
	// Echo input back, or do custom logic
	return input, nil
}
```

Use an address that does **not** conflict with default precompiles (e.g. 1–9 and any other addresses defined in your go-ethereum fork). Addresses like `0x0100...` are often used for custom precompiles.

---

## 2. (Optional) Use Cosmos context and native state

If your precompile must read/write Cosmos state or emit Cosmos events, it needs the **StateDB** that the EVM uses. In go-ethereum, `vm.Contract` gives access to the EVM and its `StateDB`. In Ethermint, the StateDB implementation can expose:

- **`Context() sdk.Context`** – So the precompile can read block height, KVStores, etc.
- **`ExecuteNativeAction(contract common.Address, converter EventConverter, action func(ctx sdk.Context) error) error`** – Run Cosmos-side logic (e.g. balance changes, module stores) in an isolated way; state is reverted if the action errors or the EVM call reverts.

**Example: precompile that uses Cosmos context**

```go
func (p MyPrecompile) Run(input []byte, contract *vm.Contract) ([]byte, error) {
	// Get StateDB; assert to Ethermint's type to access Context()
	db := contract.EVM().StateDB
	// If your app uses statedb.HookedStateDB or a wrapper, unwrap until you get one that implements Context().
	type contextDB interface {
		Context() sdk.Context
	}
	if cdb, ok := db.(contextDB); ok {
		ctx := cdb.Context()
		// Use ctx for reads: ctx.BlockHeight(), ctx.KVStore(storeKey), etc.
		_ = ctx
	}
	// For state-changing native logic, use ExecuteNativeAction if your StateDB exposes it.
	return input, nil
}
```

Only use `ExecuteNativeAction` (and the StateDB type that provides it) if your EVM is actually using Ethermint’s StateDB implementation that supports it.

---

## 3. Create a `CustomContractFn` (factory)

Ethermint’s keeper expects a **slice of functions** that produce a precompile for each (context, rules) pair. Each function has the type:

```go
type CustomContractFn func(sdk.Context, ethparams.Rules) vm.PrecompiledContract
```

Create a factory that returns your precompile:

```go
package myprecompile

import (
	ethparams "github.com/ethereum/go-ethereum/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
)

func NewMyPrecompileFn() evmkeeper.CustomContractFn {
	return func(ctx sdk.Context, rules ethparams.Rules) vm.PrecompiledContract {
		// You can use ctx/rules to enable/disable or configure the precompile per block.
		return MyPrecompile{}
	}
}
```

If you need to inject keepers or config, capture them in the closure:

```go
// types = e.g. github.com/evmos/ethermint/x/evm/types
func NewMyPrecompileFn(accountKeeper types.AccountKeeper) evmkeeper.CustomContractFn {
	return func(ctx sdk.Context, rules ethparams.Rules) vm.PrecompiledContract {
		return MyPrecompile{accountKeeper: accountKeeper}
	}
}
```

---

## 4. Register the precompile in the app

When constructing the EVM keeper, pass your factory (or factories) as the **custom precompiles** slice. In `evmd/app.go` the keeper is created with `nil` for that slice; replace it with your list:

**Before (snippet from `evmd/app.go`):**

```go
app.EvmKeeper = evmkeeper.NewKeeper(
	appCodec,
	keys[evmtypes.StoreKey], okeys[evmtypes.ObjectStoreKey], authtypes.NewModuleAddress(govtypes.ModuleName),
	app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.FeeMarketKeeper,
	tracer,
	evmSs,
	nil,  // <-- custom precompiles
	cast.ToUint64(appOpts.Get(server.FlagQueryGasLimit)),
)
```

**After:**

```go
import (
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	myprecompile "github.com/evmos/ethermint/x/myprecompile" // or your module path
)

// ...

app.EvmKeeper = evmkeeper.NewKeeper(
	appCodec,
	keys[evmtypes.StoreKey], okeys[evmtypes.ObjectStoreKey], authtypes.NewModuleAddress(govtypes.ModuleName),
	app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.FeeMarketKeeper,
	tracer,
	evmSs,
	[]evmkeeper.CustomContractFn{
		myprecompile.NewMyPrecompileFn(),
		// add more custom precompiles here
	},
	cast.ToUint64(appOpts.Get(server.FlagQueryGasLimit)),
)
```

After this, every new EVM instance will get both `vm.DefaultPrecompiles(cfg.Rules)` and your custom precompiles, and `vm.ActivePrecompiles(rules)` plus your custom addresses are passed to `StateDB.Prepare` for access-list handling.

---

## 5. Avoid address conflicts

- **Preinstalls**: If you use `AddPreinstalls` (contracts deployed at fixed addresses), the keeper checks that those addresses do **not** match any precompile address (including those from `customContractFns`). Pick precompile addresses that won’t clash with preinstalls.
- **Default precompiles**: Don’t reuse addresses 1–9 (and any others your go-ethereum fork defines). Use a distinct range (e.g. `0x0100`, `0x0101`, …) for custom precompiles.

---

## 6. Summary checklist

1. Implement a type with `Address() common.Address`, `RequiredGas(input []byte) uint64`, and `Run(input []byte, contract *vm.Contract) ([]byte, error)` so it satisfies `vm.PrecompiledContract`.
2. Choose a dedicated, non-conflicting address for your precompile.
3. Add a `CustomContractFn` that returns that precompile (optionally using `sdk.Context` and `ethparams.Rules`).
4. Pass that function in the `customContractFns` slice to `evmkeeper.NewKeeper` in `app.go`.
5. If you need Cosmos state or events, use the StateDB’s `Context()` and/or `ExecuteNativeAction` inside `Run` (with the StateDB type your app actually uses).

Contracts can then call your precompile by doing a normal `CALL`/`STATICCALL` to your chosen address with the ABI-encoded input; the EVM will dispatch to your `Run` and charge `RequiredGas`.
