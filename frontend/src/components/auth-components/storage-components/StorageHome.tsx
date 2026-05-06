import { useEffect, type JSX } from "react";
import { useLoaderData, useNavigate } from "react-router";
import { type User } from "../../../middleware";
import { fetchApi } from "../../../functions/fetchtils";
import { useFileStore } from "../../../context/FileStore";

export default function StorageHome(): JSX.Element{
    const loaderData = useLoaderData<User>();
    const {setFiles} = useFileStore();
    const navigate = useNavigate();

    useEffect(() => {
        setFiles();
    }, []);

    /**
     * Logouts the current validated user. This uses the session ID found
     * in the cookies.
     * 
     * If successful, it will logout the current user, invalidate the session in the
     * cookie, and redirect back to the home page.
     */
    async function logout(): Promise<void>{
        try{
            const res = await fetchApi<boolean>("/api/logout", "POST");

            if(res.output){
                navigate("/");
            }
        }catch(err){
            // TODO: add error popup here
            console.error(err);
        }
    }

    return (
        <>
            <div className="flex flex-col justify-center items-center">
                TEMPORARY: Hello {loaderData.username}. You are logged in.
                <button onClick={logout} className="border w-fit h-fit py-2 px-4">Logout</button>
            </div> 
        </>
    )
}