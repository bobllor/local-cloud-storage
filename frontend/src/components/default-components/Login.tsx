import type React from "react";
import type { JSX } from "react";
import { NavLink } from "react-router";
import { createUrl } from "../../utils";

const inputFieldClass: string = "border-2 w-50 h-10"; 
const formInputNames = {
    username: "username",
    password: "password",
}

export default function Login(): JSX.Element{
    return (
        <> 
            <div>
                <form 
                method="POST"
                onSubmit={(e) => {
                    const validLogin: Promise<boolean> = login(e);
                    console.log(validLogin);
                }}
                className="flex flex-col gap-1 items-center">
                    <div className="flex flex-col gap-1 items-center justify-center">
                        <div className="flex flex-col">
                            <input name={formInputNames.username} type="text" className={inputFieldClass} />
                        </div>
                        <div className="flex flex-col">
                            <input name={formInputNames.password} type="password" className={inputFieldClass} />
                        </div>
                        <div className="flex w-full justify-center">
                            <button 
                            type="submit"
                            className="h-8 border-2 w-[40%]">
                                Login
                            </button>
                        </div>
                    </div>
                    <NavLink to="/register" className="h-8 border-2 w-[40%] flex items-center justify-center">
                        Register
                    </NavLink>
                    <NavLink to="/" className="h-8 border-2 w-[40%] flex items-center justify-center">
                        Home
                    </NavLink>
                </form>
            </div>
        </>
    )
}

/**
 * Login to the server.
 *
 * @param formEvent 
 * @returns Status of the login authentication
 */
async function login(formEvent: React.SubmitEvent<HTMLFormElement>): Promise<boolean>{
    formEvent.preventDefault();

    const formData: FormData = new FormData(formEvent.currentTarget);
    const userData = {username: "", password: ""}
    formData.forEach((v, k) => {
        if(k == formInputNames.username){
            userData.username = v.toString();
        }else{
            userData.password = v.toString();
        }
    })

    fetch(createUrl("/login"), {
        method: "POST",
        body: JSON.stringify(userData),
    }).then(res => {
        console.log(res);
        res.json().then(val => console.log(val))
    });

    return false
}