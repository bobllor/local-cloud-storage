import type { ResponseApi } from "../response-types";
import { createUrl } from "../utils";

/**
 * Sends a request to validate the current session. This sends
 * the cookie to the backend.
 * 
 * @returns The session validation status
 */
export async function validateSession(): Promise<boolean>{
    const res = await request(createUrl("/api/session"));
    const output: ResponseApi<boolean> = await res.json();

    return output.output;
}

/**
 * Sends a request to the given path and returns the response.
 * If args are used, it will send the data with args.
 * 
 * @param path The non-base request URL, this can include the forward slash
 * @param method The method used on the request
 * @param body Any object that is being sent to the backend
 * @returns Response promise
 */
export async function request(path: string, method: method = "GET", data?: {}): Promise<Response>{
    const res = await fetch(createUrl(path), {
        method: method,
        body: !data ? undefined : JSON.stringify(data),
        credentials: "include",
    })

    return res;
}

export type method = "GET" | "POST";