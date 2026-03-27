CREATE DATABASE IF NOT EXISTS TestLocalCloudStorage;
USE TestLocalCloudStorage;

CREATE TABLE IF NOT EXISTS UserAccount(
    AccountID varchar(50),
    Username varchar(30) NOT NULL UNIQUE,
    PasswordHash varchar(255) NOT NULL,
    CreatedDate DATE,
    SeenDate DATE,
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
    ModifiedDate DATE,
    DeletedDate DATE,
    PRIMARY KEY (FileID),
    CONSTRAINT FK_UserAccount
        FOREIGN KEY (FileOwnerID) 
        REFERENCES UserAccount(AccountID)
        ON DELETE CASCADE
);

INSERT INTO UserAccount
    VALUES
    ("89672a64-f3ff-490c-8f2d-7e5cf5d4aa70", 
    "test.username", 
    "3eb72b4431dff57dd10e76d0921d1787", 
    NOW(), 
    NULL);

INSERT INTO File
    VALUES
    ("89672a64-f3ff-490c-8f2d-7e5cf5d4aa70",
    "test1.txt",
    "file",
    "randomfileidhere",
    NULL,
    "/path/to/file",
    1234,
    NOW(),
    NULL
    );