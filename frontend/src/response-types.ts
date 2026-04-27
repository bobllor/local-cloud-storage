/**
 * Represents the response status of a request.
 *  - success: The request was successful and data was returned (if applicable)
 *  - error: The request was not successful, this can be due to a server-side error
 *  or bad data in the request
 */
export type ResponseStatus = "success" | "error";

/**
 * Represents a successful API response. The output can be of type T depending on
 * how the endpoint is configured. Output is only optional if an error occurs.
 */
export type ResponseApi<T> = {
    status: "success"
    output: T
    error: never
} | {
    status: "error"
    output: never
    error: ResponseError
}

/**
 * Represents an error response structure. This can be null if an
 * error does not occur.
 */
export type ResponseError = {
    code: number
    reason: ReasonCode
    message: string
}

/**
 * Codes for ReasonCode values used to determine what error has occurred
 * with the response return.
 */
export const REASON_CODES = {
    internalError: "INTERNAL_ERROR",
	badRequestData: "BAD_DATA",
	userAlreadyExists: "USER_ALREADY_EXISTS",
	badUsername: "BAD_USERNAME",
	badPassword:"BAD_PASSWORD",
	unauthorized: "UNAUTHORIZED",
}

/**
 * ReasonCode indicates what error type had occurred during the endpoint
 * call.
 */
export type ReasonCode = keyof typeof REASON_CODES