import type { JSX } from "react";
import { NavLink } from "react-router";

export default function StorageHome(): JSX.Element{
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