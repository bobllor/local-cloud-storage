/**
 * Represents the response status of a request.
 *  - success: The request was successful and data was returned (if applicable)
 *  - error: The request was not successful, this can be due to a server-side error
 *  or bad data in the request
 */
export type ResponseStatus = "success" | "error";

/**
 * Represents a successful API response.
 */
export type ResponseApi<T> = {
    status: ResponseStatus
    output: T
    error?: ResponseError
}

/**
 * Represents an error response structure. This can be null if an
 * error does not occur.
 */
export type ResponseError = {
    code: number
    message: string
}