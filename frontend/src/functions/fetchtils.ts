import type { ResponseApi } from "../response-types";
import { createUrl } from "../utils";

/**
 * Sends a request to validate the current session. This sends
 * the cookie to the backend.
 * 
 * @returns The session validation status
 */
export async function validateSession(): Promise<boolean>{
    const res = await getRequest(createUrl("/api/session"));
    const output: ResponseApi<boolean> = await res.json();

    return output.output;
}

/**
 * Sends a GET request to the given path and returns the response.
 * If args are used, it will send the data with args.
 * 
 * @param path
 * @returns Response promise
 */
export async function getRequest(path: string, body?: {}): Promise<Response>{
    const res = await fetch(path, {
        method: "GET",
        body: !body ? undefined : JSON.stringify(body),
        credentials: "include",
    })

    return res;
}
