import React, { createContext, useState, type JSX, type SetStateAction } from "react";

// Represents the File data of the database.
// See file.go for the struct on the backend.
export type dbFile = {
    accountID: string
    fileName: string
    fileType: "dir" | "file"
    parentID?: string
    filePath: string
    fileSize: number
    modifiedOn: Date
    deletedOn?: Date
}

/**
 * In-memory map representing the cache of the file state. Each key
 * represents a parent ID with all of its files as its values.
 */
export type fileMap = Map<string, Array<dbFile>>;

export const FileContext = createContext<Context>({
    fileData: new Map<string, Array<dbFile>>(),
    setFileData: () => {},
});

export default function FileProvider({children}: {children: JSX.Element}): JSX.Element{
    const [fileData, setFileData] = useState<fileMap>(new Map<string, Array<dbFile>>());

    const data: Context = {
        fileData,
        setFileData,
    };

    return (
        <FileContext value={data}>
            {children}
        </FileContext>
    );
}

type Context = {
    fileData: fileMap,
    setFileData: React.Dispatch<SetStateAction<fileMap>>,
};