/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ProtobufAny {
  typeUrl?: string;

  /** @format byte */
  value?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
 * AccessTuple is the element type of an access list.
 */
export interface V1Alpha1AccessTuple {
  address?: string;
  storageKeys?: string[];
}

/**
* Log represents an protobuf compatible Ethereum Log that defines a contract
log event. These events are generated by the LOG opcode and stored/indexed by
the node.
*/
export interface V1Alpha1Log {
  address?: string;

  /** list of topics provided by the contract. */
  topics?: string[];

  /** @format byte */
  data?: string;

  /** @format uint64 */
  blockNumber?: string;
  txHash?: string;

  /** @format uint64 */
  txIndex?: string;
  blockHash?: string;

  /** @format uint64 */
  index?: string;

  /**
   * The Removed field is true if this log was reverted due to a chain
   * reorganisation. You must pay attention to this field if you receive logs
   * through a filter query.
   */
  removed?: boolean;
}

/**
 * MsgEthereumTxResponse defines the Msg/EthereumTx response type.
 */
export interface V1Alpha1MsgEthereumTxResponse {
  /**
   * contract_address contains the ethereum address of the created contract (if
   * any). If the state transition is an evm.Call, the contract address will be
   * empty.
   */
  contractAddress?: string;

  /** @format byte */
  bloom?: string;

  /**
   * tx_logs contains the transaction hash and the proto-compatible ethereum
   * logs.
   */
  txLogs?: V1Alpha1TransactionLogs;

  /**
   * ret defines the bytes from the execution.
   * @format byte
   */
  ret?: string;
  reverted?: boolean;
}

export interface V1Alpha1Params {
  /**
   * evm_denom represents the token denomination used to run the EVM state
   * transitions.
   */
  evmDenom?: string;
  enableCreate?: boolean;
  enableCall?: boolean;
  extraEips?: string[];
}

/**
 * QueryAccountResponse is the response type for the Query/Account RPC method.
 */
export interface V1Alpha1QueryAccountResponse {
  /** balance is the balance of the EVM denomination. */
  balance?: string;

  /**
   * code_hash is the code bytes from the EOA.
   * @format byte
   */
  codeHash?: string;

  /**
   * nonce is the account's sequence number.
   * @format uint64
   */
  nonce?: string;
}

/**
 * QueryBalanceResponse is the response type for the Query/Balance RPC method.
 */
export interface V1Alpha1QueryBalanceResponse {
  /** balance is the balance of the EVM denomination. */
  balance?: string;
}

/**
* QueryBlockBloomResponse is the response type for the Query/BlockBloom RPC
method.
*/
export interface V1Alpha1QueryBlockBloomResponse {
  /**
   * bloom represents bloom filter for the given block hash.
   * @format byte
   */
  bloom?: string;
}

/**
 * QueryTxLogs is the response type for the Query/BlockLogs RPC method.
 */
export interface V1Alpha1QueryBlockLogsResponse {
  /** logs represents the ethereum logs generated at the given block hash. */
  txLogs?: V1Alpha1TransactionLogs[];
}

/**
* QueryCodeResponse is the response type for the Query/Code RPC
method.
*/
export interface V1Alpha1QueryCodeResponse {
  /**
   * code represents the code bytes from an ethereum address.
   * @format byte
   */
  code?: string;
}

/**
 * QueryCosmosAccountResponse is the response type for the Query/CosmosAccount RPC method.
 */
export interface V1Alpha1QueryCosmosAccountResponse {
  /** cosmos_address is the cosmos address of the account. */
  cosmosAddress?: string;

  /**
   * sequence is the account's sequence number.
   * @format uint64
   */
  sequence?: string;

  /** @format uint64 */
  accountNumber?: string;
}

/**
 * QueryParamsResponse defines the response type for querying x/evm parameters.
 */
export interface V1Alpha1QueryParamsResponse {
  /** params define the evm module parameters. */
  params?: V1Alpha1Params;
}

export interface V1Alpha1QueryStaticCallResponse {
  /** @format byte */
  data?: string;
}

/**
* QueryStorageResponse is the response type for the Query/Storage RPC
method.
*/
export interface V1Alpha1QueryStorageResponse {
  /** key defines the storage state value hash associated with the given key. */
  value?: string;
}

/**
 * QueryTxLogs is the response type for the Query/TxLogs RPC method.
 */
export interface V1Alpha1QueryTxLogsResponse {
  /** logs represents the ethereum logs generated from the given transaction. */
  logs?: V1Alpha1Log[];
}

/**
 * QueryTxReceiptResponse is the response type for the Query/TxReceipt RPC method.
 */
export interface V1Alpha1QueryTxReceiptResponse {
  /** receipt represents the ethereum receipt for the given transaction. */
  receipt?: V1Alpha1TxReceipt;
}

/**
 * QueryTxReceiptsByBlockHashResponse is the response type for the Query/TxReceiptsByBlockHash RPC method.
 */
export interface V1Alpha1QueryTxReceiptsByBlockHashResponse {
  receipts?: V1Alpha1TxReceipt[];
}

/**
 * QueryTxReceiptsByBlockHeightResponse is the response type for the Query/TxReceiptsByBlockHeight RPC method.
 */
export interface V1Alpha1QueryTxReceiptsByBlockHeightResponse {
  receipts?: V1Alpha1TxReceipt[];
}

/**
* TransactionLogs define the logs generated from a transaction execution
with a given hash. It it used for import/export data as transactions are not
persisted on blockchain state after an upgrade.
*/
export interface V1Alpha1TransactionLogs {
  hash?: string;
  logs?: V1Alpha1Log[];
}

/**
* TxData implements the Ethereum transaction data structure. It is used
solely as intended in Ethereum abiding by the protocol.
*/
export interface V1Alpha1TxData {
  /** @format byte */
  chainId?: string;

  /**
   * nonce corresponds to the account nonce (transaction sequence).
   * @format uint64
   */
  nonce?: string;

  /**
   * price defines the unsigned integer value of the gas price in bytes.
   * @format byte
   */
  gasPrice?: string;

  /**
   * gas defines the gas limit defined for the transaction.
   * @format uint64
   */
  gas?: string;
  to?: string;

  /**
   * value defines the unsigned integer value of the transaction amount.
   * @format byte
   */
  value?: string;

  /**
   * input defines the data payload bytes of the transaction.
   * @format byte
   */
  input?: string;
  accesses?: V1Alpha1AccessTuple[];

  /** @format byte */
  v?: string;

  /** @format byte */
  r?: string;

  /** @format byte */
  s?: string;
}

/**
 * TxReceipt defines the receipt type stored in KV for each EVM transaction.
 */
export interface V1Alpha1TxReceipt {
  hash?: string;
  from?: string;

  /**
   * TxData implements the Ethereum transaction data structure. It is used
   * solely as intended in Ethereum abiding by the protocol.
   */
  data?: V1Alpha1TxData;

  /** TxResult stores results of Tx execution. */
  result?: V1Alpha1TxResult;

  /** @format uint64 */
  index?: string;

  /** @format uint64 */
  blockHeight?: string;
  blockHash?: string;
}

/**
 * TxResult stores results of Tx execution.
 */
export interface V1Alpha1TxResult {
  /**
   * contract_address contains the ethereum address of the created contract (if
   * any). If the state transition is an evm.Call, the contract address will be
   * empty.
   */
  contractAddress?: string;

  /** @format byte */
  bloom?: string;

  /**
   * tx_logs contains the transaction hash and the proto-compatible ethereum
   * logs.
   */
  txLogs?: V1Alpha1TransactionLogs;

  /**
   * ret defines the bytes from the execution.
   * @format byte
   */
  ret?: string;
  reverted?: boolean;

  /** @format uint64 */
  gasUsed?: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: keyof Omit<Body, "body" | "bodyUsed">;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType) => RequestParams | void;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "";
  private securityData: SecurityDataType = null as any;
  private securityWorker: null | ApiConfig<SecurityDataType>["securityWorker"] = null;
  private abortControllers = new Map<CancelToken, AbortController>();

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType) => {
    this.securityData = data;
  };

  private addQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];

    return (
      encodeURIComponent(key) +
      "=" +
      encodeURIComponent(Array.isArray(value) ? value.join(",") : typeof value === "number" ? value : `${value}`)
    );
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) =>
        typeof query[key] === "object" && !Array.isArray(query[key])
          ? this.toQueryString(query[key] as QueryParamsType)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((data, key) => {
        data.append(key, input[key]);
        return data;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  private mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format = "json",
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams = (secure && this.securityWorker && this.securityWorker(this.securityData)) || {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];

    return fetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      signal: cancelToken ? this.createAbortSignal(cancelToken) : void 0,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = (null as unknown) as T;
      r.error = (null as unknown) as E;

      const data = await response[format]()
        .then((data) => {
          if (r.ok) {
            r.data = data;
          } else {
            r.error = data;
          }
          return r;
        })
        .catch((e) => {
          r.error = e;
          return r;
        });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title ethermint/evm/v1alpha1/tx.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryAccount
   * @summary Account queries an Ethereum account.
   * @request GET:/ethermint/evm/v1alpha1/account/{address}
   */
  queryAccount = (address: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryAccountResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/account/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Query
 * @name QueryBalance
 * @summary Balance queries the balance of a the EVM denomination for a single
EthAccount.
 * @request GET:/ethermint/evm/v1alpha1/balances/{address}
 */
  queryBalance = (address: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryBalanceResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/balances/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryBlockBloom
   * @summary BlockBloom queries the block bloom filter bytes at a given height.
   * @request GET:/ethermint/evm/v1alpha1/block_bloom
   */
  queryBlockBloom = (query?: { height?: string }, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryBlockBloomResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/block_bloom`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryBlockLogs
   * @summary BlockLogs queries all the ethereum logs for a given block hash.
   * @request GET:/ethermint/evm/v1alpha1/block_logs/{hash}
   */
  queryBlockLogs = (hash: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryBlockLogsResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/block_logs/${hash}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryCode
   * @summary Code queries the balance of all coins for a single account.
   * @request GET:/ethermint/evm/v1alpha1/codes/{address}
   */
  queryCode = (address: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryCodeResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/codes/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryCosmosAccount
   * @summary Account queries an Ethereum account's Cosmos Address.
   * @request GET:/ethermint/evm/v1alpha1/cosmos_account/{address}
   */
  queryCosmosAccount = (address: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryCosmosAccountResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/cosmos_account/${address}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Params queries the parameters of x/evm module.
   * @request GET:/ethermint/evm/v1alpha1/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<V1Alpha1QueryParamsResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStaticCall
   * @summary StaticCall queries the static call value of x/evm module.
   * @request GET:/ethermint/evm/v1alpha1/static_call
   */
  queryStaticCall = (query?: { address?: string; input?: string }, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryStaticCallResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/static_call`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryStorage
   * @summary Storage queries the balance of all coins for a single account.
   * @request GET:/ethermint/evm/v1alpha1/storage/{address}/{key}
   */
  queryStorage = (address: string, key: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryStorageResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/storage/${address}/${key}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTxLogs
   * @summary TxLogs queries ethereum logs from a transaction.
   * @request GET:/ethermint/evm/v1alpha1/tx_logs/{hash}
   */
  queryTxLogs = (hash: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryTxLogsResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/tx_logs/${hash}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTxReceipt
   * @summary TxReceipt queries a receipt by a transaction hash.
   * @request GET:/ethermint/evm/v1alpha1/tx_receipt/{hash}
   */
  queryTxReceipt = (hash: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryTxReceiptResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/tx_receipt/${hash}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTxReceiptsByBlockHeight
   * @summary TxReceiptsByBlockHeight queries tx receipts by a block height.
   * @request GET:/ethermint/evm/v1alpha1/tx_receipts_block/{height}
   */
  queryTxReceiptsByBlockHeight = (height: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryTxReceiptsByBlockHeightResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/tx_receipts_block/${height}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTxReceiptsByBlockHash
   * @summary TxReceiptsByBlockHash queries tx receipts by a block hash.
   * @request GET:/ethermint/evm/v1alpha1/tx_receipts_block_hash/{hash}
   */
  queryTxReceiptsByBlockHash = (hash: string, params: RequestParams = {}) =>
    this.request<V1Alpha1QueryTxReceiptsByBlockHashResponse, RpcStatus>({
      path: `/ethermint/evm/v1alpha1/tx_receipts_block_hash/${hash}`,
      method: "GET",
      format: "json",
      ...params,
    });
}
