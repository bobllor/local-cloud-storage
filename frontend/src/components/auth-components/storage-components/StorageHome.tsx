import { type JSX } from "react";
import { useLoaderData } from "react-router";
import { type User } from "../../../middleware";

export default function StorageHome(): JSX.Element{
    const loaderData = useLoaderData<User>();

    return (
        <>
            <div className="flex flex-col justify-center items-center">
                TEMPORARY: Hello {loaderData.username}. You are logged in.
                <button onClick={logout} className="border w-fit h-fit py-2 px-4">Logout</button>
            </div> 
        </>
    )
}

/**
 * Logouts the current validated user. This uses the session ID found
 * in the cookies.
 */
async function logout(): Promise<void>{

}