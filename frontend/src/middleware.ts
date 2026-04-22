import { createContext, redirect, type LoaderFunction, type MiddlewareFunction } from "react-router";
import { createUrl } from "./utils";
import type { ResponseApi } from "./response-types";

export const userContext = createContext<User|null>(null);

type User = {
    account_id: string
    username: string
    created_on: Date
    active: boolean
}

/**
 * Sends a request to the backend for authentication using the cookie
 * in the request header. This validates if the session ID is correct
 * and the user is valid.
 * 
 * Once valid, it will set the user context and it can be used for all routes
 * afterwards.
 * 
 * Authentication failures will redirect to the login page.
 */
export const authMiddleware: MiddlewareFunction = async ({context}) => {
    const user = await getUser();

    if(!user){
        throw redirect("/login")
    }

    // TODO: log res not in console.log
    context.set(userContext, user);
}

/**
 * Sends a request to the backend for authentication using the cookie
 * in the request header. This validates if the session ID is correct
 * and the user is valid.
 * 
 * Upon a successful authentication it will redirect to the storage,
 * bypassing the login page.
 */
export const loginMiddleware: MiddlewareFunction = async () => {
    const res: ResponseApi<User> = await fetch(createUrl("/api/user"), {
        method: "GET",
        credentials: "include",
    }).then(val => val.json());

    if(res.status == "success"){
        throw redirect("/storage");
    }
}

export const getUserContext: LoaderFunction = async ({context}) => {
    const user = context.get(userContext);

    return user;
}

export async function getUser(): Promise<User|null>{
    const res: ResponseApi<User> = await fetch(createUrl("/api/user"), {
        method: "GET",
        credentials: "include",
    }).then(val => val.json());

    if (res.status == "error"){
        return null;
    }

    return res.output;
}