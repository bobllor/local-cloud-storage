import { SERVER_BASE_URL } from "./env";

/**
 * Creates the URL from the given base URL in the .env file. This requires the .env file
 * located in the root folder of the vite files.
 *  
 * @param path The string path of the url, this must start with a forward slash.
 * @returns The combined URL
 */
export function createUrl(path: string): string{
    var url: string = SERVER_BASE_URL + path;

    return url;
}