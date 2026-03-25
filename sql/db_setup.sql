CREATE DATABASE IF NOT EXISTS MasterLocalCloudStorage;
USE MasterLocalCloudStorage;

CREATE TABLE IF NOT EXISTS Files(
    AccountID varchar(255) NOT NULL,
    FileName varchar(255),
    FileType varchar(255),
    FileID varchar(255) PRIMARY KEY,
    ParentID varchar(255),
    FilePath varchar(5120),
    FileSize int,
    ModifiedDate DATE,
    DeletedTime DATE
);

SELECT * FROM Files;