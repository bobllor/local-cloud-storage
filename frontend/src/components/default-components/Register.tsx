import type { JSX } from "react";

const registerInputNames = {
    username: "username",
    password: "default-password",
    confirmPassword: "confirm-password",
}

export default function Register(): JSX.Element{
    return (
        <>
            <form 
            onSubmit={e => {e.preventDefault(); console.log("submitted")}}
            className="flex flex-col justify-center items-center gap-1">
                <input name={registerInputNames.username} className="user-input-field" />
                <input name={registerInputNames.password} className="user-input-field" type="password" />
                <input name={registerInputNames.confirmPassword} className="user-input-field" type="password" />
                <button type="submit" className="border w-20 h-10">Register</button>
            </form>
        </>
    )
}