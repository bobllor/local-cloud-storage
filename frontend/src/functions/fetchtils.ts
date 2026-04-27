import type { ResponseApi } from "../response-types";
import { createUrl } from "../utils";

/**
 * Sends a request to validate the current session. This sends
 * the cookie to the backend.
 * 
 * @returns The session validation status
 */
export async function validateSession(): Promise<boolean>{
    const res = await fetchApi<boolean>("/api/session");

    if(res.status == "error"){
        return false;
    }

    return res.output;
}

/**
 * Sends a request to the given path and returns its response.
 * If args are used, it will send the data with args.
 * 
 * @param path The non-base request URL, this can include the forward slash
 * @param method The method used on the request
 * @param body Any object that is being sent to the backend
 * @returns ResponseApi promise of type T
 */
export async function fetchApi<T>(path: string, method: method = "GET", data?: {}): Promise<ResponseApi<T>>{
    const res = await fetch(createUrl(path), {
        method: method,
        body: !data ? undefined : JSON.stringify(data),
        credentials: "include",
    });

    const r: ResponseApi<T> = await res.json()

    return r;
}

export type method = "GET" | "POST";