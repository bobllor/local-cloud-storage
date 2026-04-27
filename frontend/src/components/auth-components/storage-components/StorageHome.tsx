import { type JSX } from "react";
import { NavLink, useLoaderData } from "react-router";
import { type User } from "../../../middleware";

export default function StorageHome(): JSX.Element{
    const loaderData = useLoaderData<User>();

    return (
        <>
            <div className="flex flex-col justify-center items-center">
                TEMPORARY: Hello {loaderData.username}. You are logged in.
                <NavLink to={"/"}>
                    Home
                </NavLink>
            </div> 
        </>
    )
}