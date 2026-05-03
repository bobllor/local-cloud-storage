import type React from "react";
import { useState, type JSX } from "react";
import { NavLink, useNavigate } from "react-router";
import { HiOutlineXMark } from "react-icons/hi2";
import { fetchApi } from "../../functions/fetchtils";
import { useSessionValidation } from "./hooks";

const formInputNames = {
    username: "username",
    password: "password",
}

export default function Login(): JSX.Element{
    const [loginStatus, setLoginStatus]= useState<null|boolean>(null);
    const navigate = useNavigate();

    useSessionValidation();

    return (
        <> 
            <div>
                <form 
                method="POST"
                onSubmit={async (e) => {
                    login(e).then(status => {
                        setLoginStatus(status);

                        if(status){
                            navigate("/storage");
                        }
                    });
                }}
                className="flex flex-col gap-1 items-center">
                    {loginStatus != null && !loginStatus &&
                        <div className="flex border-2 h-1 items-center justify-center bg-red-700/60">
                            <div className="px-5">
                                <div className="flex gap-2 justify-between items-center">
                                    <p>
                                        Incorrect username or password
                                    </p>
                                    <button 
                                    onClick={e => handleErrorClose(e, setLoginStatus)}
                                    className="flex items-center justify-center hover:bg-gray-400/70 p-1">
                                        <HiOutlineXMark />
                                    </button>
                                </div>
                            </div>
                        </div>
                    }
                    <div className="flex flex-col gap-1 items-center justify-center">
                        <div className="flex flex-col">
                            <input name={formInputNames.username} type="text" className={"user-input-field"} />
                        </div>
                        <div className="flex flex-col">
                            <input name={formInputNames.password} type="password" className={"user-input-field"} />
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
 * Login to the server. If the login was successful, it will return the status.
 *
 * @param formEvent 
 * @returns Status of the login authentication
 */
async function login(formEvent: React.SubmitEvent<HTMLFormElement>): Promise<boolean>{
    formEvent.preventDefault();

    const formData: FormData = new FormData(formEvent.currentTarget);
    const reqObj: Record<string, any> = {};

    formData.forEach((v, k) => {
        reqObj[k] = v;
    });

    const res = await fetchApi<boolean>("/api/login", "POST", reqObj);

    if(res.status == "error"){
        return false;
    }

    return res.output;
}

/**
 * Handles closing the error login popup.
 * @param event 
 */
function handleErrorClose(
    event: React.MouseEvent<HTMLButtonElement>, 
    setStatus: React.Dispatch<React.SetStateAction<null|boolean>>
): void{
    event.preventDefault();

    setStatus(null);
}