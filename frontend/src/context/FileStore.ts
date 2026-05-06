import { create } from "zustand";
import { fetchApi } from "../functions/fetchtils";

type FileStore = {
    files: Record<string, Array<File>>,
    setFiles: (parentID?: string) => {},
}

/**
 * The type representing the File data of the database.
 * This does not include the file path.
 */
type File = {
    accountID: string
    fileName: string
    fileType: string
    fileID: string
    parentID: string
    modifiedOn: Date
    deletedOn?: Date
}

const ROOT_KEY = "root";

export const useFileStore = create<FileStore>((set, get) => ({
    files: {},
    /**
     * Sets the contents of the files based on the parentID. If the parentID already
     * has an entry, then this will do nothing.
     * @param parentID 
     */
    setFiles: async (parentID?: string) => {
        const key = parentID ? parentID : ROOT_KEY;
        const route = parentID ? `/api/storage/folder/${key}` : "/api/storage";
        const baseFiles = get().files;

        // will not update the state if it already exists
        if(key in baseFiles){
            return;
        }

        const newFiles = await fetchApi<Array<File>>(route);

        set(state => ({state, files: {key: newFiles.output}}));
        // TODO: remove this or something idk
        console.log(get().files);
    },
    updateFiles: () => {},
    /**
     * Retrieves the files based on the parentID.
     * @param parentID The parentID of the files, this can be null indicating it is the root folder
     * @returns The array of the files related to the parentID
     */
    getFiles: (parentID?: string) => {
        const files = get().files;
        const key = parentID ? parentID : ROOT_KEY;

        return files[key];
    }
}));
