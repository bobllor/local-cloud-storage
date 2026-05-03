import { useEffect } from "react";
import { validateSession } from "../../functions/fetchtils";
import { useNavigate } from "react-router";

/**
 * Hook used to validate the session ID from the cookie.
 * If it is a valid session, then a redirect to the authenticated 
 * storage page will occur. Otherwise, it will do nothing.
 */
export function useSessionValidation(){
    const navigate = useNavigate();

    useEffect(() => {
        validateSession().then(status => {
            // TODO: console dot log (real logging please)
            // temporary for now just to validate
            console.log(`Login: ${status}`);
            if(status){
                navigate("/storage"); 
            }
        });
    }, []);
}