import type React from "react";
import { useState, type JSX } from "react";
import { fetchApi } from "../../functions/fetchtils";
import { NavLink, useNavigate } from "react-router";
import { useSessionValidation } from "./hooks";

const registerInputNames = {
    username: "username",
    password: "password",
    confirmPassword: "confirmPassword",
}

export default function Register(): JSX.Element{
    useSessionValidation();

    const [inputErrors, setInputErrors] = useState<inputObject>({username: "", password: ""});
    const navigate = useNavigate();

    /**
     * Sets the error message for the inputObject state.
     * @param errType The field of the error type
     * @param msg The message of the error
     */
    const setError = (errType: "username" | "password", msg: string) => {
        setInputErrors(prev => ({...prev, [errType]: msg}));
    }
    
    /**
     * Registers the user. Upon successful registration, the session ID is written to the cookie
     * and it will redirect to the main storage page.
     * @param event The HTML form submission event
     */
    const registerUser = async (event: React.SubmitEvent<HTMLFormElement>): Promise<void> => {
        event.preventDefault();

        const data = new FormData(event.currentTarget);

        const username = data.get(registerInputNames.username);
        const password = data.get(registerInputNames.password);
        const confirmPassword = data.get(registerInputNames.confirmPassword);

        let error = false;

        const userMsg = validateUsername(username);
        if(userMsg !== null){
            setError("username", userMsg);
            error = true;
        }else{
            setError("username", "");
        }
        const pwMsg = validatePassword(password, confirmPassword);
        if(pwMsg !== null){
            setError("password", pwMsg);
            error = true;
        }else{
            setError("password", "");
        }

        if(error){
            return;
        }

        const reqObj: Record<string, any> = {};

        data.forEach((v, k) => {
            reqObj[k] = v;
        })

        const res = await fetchApi<boolean>("/api/register", "POST", reqObj);

        if(res.status === "error"){
            const errMsg = res.error.message;
            if(res.error.reason === "BAD_PASSWORD"){
                setError("password", errMsg);
            }else if(res.error.reason === "BAD_USERNAME" || res.error.reason === "USER_ALREADY_EXISTS"){
                setError("username", errMsg);
            }else{
               // TODO: log this, any other reason is not expected- something went wrong with the backend
               // TODO: make a notification/popup for this
               console.error("Invalid server response");
            }

            return;
        }

        if(res.status === "success"){
            navigate("/storage");
        }
    }

    return (
        <>
            <form 
            onSubmit={registerUser}
            className="flex flex-col justify-center items-center gap-1">
                <div className="flex flex-col justify-center items-center w-60">
                    <input name={registerInputNames.username} className={`user-input-field ${inputErrors.username && "border-red-500"}`} />
                    {inputErrors.username && 
                    <span className="text-wrap text-sm">
                        {inputErrors.username}
                    </span>
                    }
                </div>
                <div className="flex flex-col justify-center items-center w-60">
                    <input 
                        name={registerInputNames.password} 
                        className={`user-input-field ${inputErrors.password && "border-red-500"}`} 
                        type="password" 
                    />
                    {inputErrors.password && 
                    <span className="text-wrap text-sm">
                        {inputErrors.password}
                    </span>
                    }
                </div>
                <input name={registerInputNames.confirmPassword} className={`user-input-field`} type="password" />
                <button type="submit" className="border w-20 h-10">Register</button>
                <div className="flex gap-2 justify-center items-center">
                    <span>
                        Already have an account?
                    </span>
                    <NavLink className="" to={"/login"}>
                        Sign in
                    </NavLink>
                </div>
            </form>
        </>
    )
}

type inputObject = {
    username: string
    password: string
}

function validateUsername(username: FormDataEntryValue | null): string | null{
    if(!username || username.toString() == ""){
        return "Username cannot be empty";
    }

    const strUser = username.toString();
    // NOTE: matches the same as the backend
    // TODO: look into some file that can be shared with backend + frontend
    const minStrLength = 6;
    const maxStrLength = 32;

    if(strUser.length > maxStrLength || strUser.length < minStrLength){
        return `Username must be between ${minStrLength} to ${maxStrLength} characters long`;
    }
    const alphaOnlyReg = /[A-Za-z]/;
    const doublePeriodReg = /(.*[..]{2}.*)/;
    const normalUserReg = `^([A-Za-z0-9.]+)$`

    if(!strUser[0].match(alphaOnlyReg)){
        return "Username must start with a letter";
    }
    if(!strUser[strUser.length - 1].match(/[A-Za-z0-9]/)){
        return "Username must end with an alphanumeric character";
    }
    if(strUser.match(doublePeriodReg) || !strUser.match(normalUserReg)){
        return "Username must only consist of alphanumeric characters and single periods";
    }

    return null;
}

/**
 * Validates the password. If a string is returned, then it indicates an error.
 * A valid password will result in a null value.
 * 
 * @param pw 
 * @param confirmPw 
 * @returns 
 */
function validatePassword(pw: FormDataEntryValue | null, confirmPw: FormDataEntryValue | null): string | null{
    if(!pw || pw.toString() == ""){
        return "Password cannot be empty";
    }
    const minLength = 8;
    const maxLength = 64;

    const pwStr = pw.toString()

    if(pwStr.length > maxLength || pwStr.length < minLength){
        return `Password must be between ${minLength} and ${maxLength} characters long`
    }

    if(pw !== confirmPw){
        return "Passwords do not match";
    }

    return null;
}