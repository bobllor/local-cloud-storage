CREATE DATABASE IF NOT EXISTS MasterLocalCloudStorage;
USE MasterLocalCloudStorage;

CREATE TABLE IF NOT EXISTS UserAccount(
    AccountID varchar(50),
    Username varchar(30) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    CreatedOn DATETIME,
    SeenOn DATETIME,
    PRIMARY KEY (AccountID)
);

CREATE TABLE IF NOT EXISTS File(
    FileOwnerID varchar(50) NOT NULL,
    FileName varchar(255),
    FileType varchar(9),
    FileID varchar(50),
    ParentID varchar(255),
    FilePath varchar(5120),
    FileSize int,
    ModifiedOn DATETIME,
    DeletedOn DATETIME,
    PRIMARY KEY (FileID),
    CONSTRAINT FK_UserAccount
        FOREIGN KEY (FileOwnerID) 
        REFERENCES UserAccount(AccountID)
        ON DELETE CASCADE
);