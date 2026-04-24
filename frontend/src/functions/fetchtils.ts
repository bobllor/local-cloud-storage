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
 * @param path The path of the request URL
 * @param method The method used on the request
 * @param body A BodyInit object used to pass data to the backend
 * @returns Response promise
 */
export async function request(path: string, method: string = "GET", body?: BodyInit): Promise<Response>{
    const res = await fetch(path, {
        method: method,
        body: !body ? undefined : JSON.stringify(body),
        credentials: "include",
    })

    return res;
}
