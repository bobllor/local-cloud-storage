import type { JSX } from "react";
import type { File } from "../../../context/FileStore";
import { NavLink } from "react-router";

export default function FileListDisplay({files}: FileListDisplayProps): JSX.Element{
    return(
        <>
            <div>
                {
                    files.map((fileObj, i) => 
                        <div key={i}>
                            {
                                fileObj.fileType === "file" ?
                                <div>
                                    {fileObj.fileName}
                                </div>
                                :
                                <NavLink to={`folder/${fileObj.fileID}`}>
                                    {fileObj.fileName}
                                </NavLink>
                            }
                        </div>
                    )
                }
            </div> 
        </>
    )
}

type FileListDisplayProps = {
    files: Array<File>,
}