import { type JSX } from "react";
import { NavLink, useLoaderData } from "react-router";

export default function StorageHome(): JSX.Element{
    const loaderData = useLoaderData();

    return (
        <>
            <div>
                TEMPORARY: You are logged in.
                <NavLink to={"/"}>
                    Home
                </NavLink>
            </div> 
        </>
    )
}