import { SERVER_BASE_URL } from "./env";

/**
 * Creates the URL from the base URL in the .env file. This requires the .env file
 * located in the root of the React files.
 * If the given path does not include the forward slash, it will be given one.
 *  
 * @param path The string path of the url, this is the non-base path
 * @returns The combined URL
 */
export function createUrl(path: string): string{
    if(path.length > 0 && path[0] !== "/"){
        path = "/" + path;
    }

    var url: string = SERVER_BASE_URL + path;

    return url;
}