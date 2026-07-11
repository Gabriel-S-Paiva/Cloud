CREATE TABLE Requests (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('Pending', 'Rejected')) DEFAULT 'Pending'
);

CREATE TABLE Users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('User', 'Admin')) DEFAULT 'User',
    quota INTEGER NOT NULL DEFAULT 1073741824,
    quota_used INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE Folders (
    id INTEGER PRIMARY KEY,
    display_name TEXT NOT NULL,
    owned_by INTEGER NOT NULL,
    parent_folder INTEGER,

    FOREIGN KEY (owned_by) REFERENCES Users(id),
    FOREIGN KEY (parent_folder) REFERENCES Folders(id)
);

CREATE TABLE Files (
    id INTEGER PRIMARY KEY,
    display_name TEXT NOT NULL,
    owned_by INTEGER NOT NULL,
    size INTEGER NOT NULL,
    uploaded_at INTEGER NOT NULL DEFAULT (unixepoch()),
    last_modified INTEGER NOT NULL DEFAULT (unixepoch()),
    parent_folder INTEGER,

    FOREIGN KEY (owned_by) REFERENCES Users(id),
    FOREIGN KEY (parent_folder) REFERENCES Folders(id)
);

CREATE TABLE Shares (
    id INTEGER PRIMARY KEY,
    file INTEGER,
    folder INTEGER,
    shared_with INTEGER NOT NULL,
    permissions TEXT NOT NULL CHECK(permissions IN ('View','Edit')),

    CHECK ((file IS NOT NULL AND folder IS NULL) OR (file IS NULL AND folder IS NOT NULL)),
    FOREIGN KEY (file) REFERENCES Files(id),
    FOREIGN KEY (folder) REFERENCES Folders(id),
    FOREIGN KEY (shared_with) REFERENCES Users(id)
);