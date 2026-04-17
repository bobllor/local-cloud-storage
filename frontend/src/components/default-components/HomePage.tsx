import type { JSX } from "react";
import { NavLink } from "react-router";

export default function HomePage(): JSX.Element{
    return (
        <>
            <div className="flex flex-col justify-center items-center">
                TEMPORARY: You are on the home page, not logged in.
                <NavLink to={"/login"} className={"border-2 w-30 flex justify-center items-center"}>
                    Login
                </NavLink>
            </div> 
        </>
    )
}