interface User {
    id: number
    username: string
    role: 'User' | 'Admin'
    quota: number
    quotaUsed: number
    rootFolderId: number
}

interface UserSummary {
    id: number
    username: string
}

interface Folder {
    id: number,
    displayName: string
    ownedBy: number
    parentFolder: number|null
}

interface SharedFolder extends Folder {
    shareId: number
    sharedWith: string
    ownedByUsername: string
    permissions: string
}

interface CloudFile {
    id: number
    displayName: string
    ownedBy: number
    size: number
    bytesReceived: number
    status: 'Uploading' | 'Complete'
    contentType: string
    uploadedAt: number
    lastModified: number
    parentFolder: number|null
}

interface SharedFile extends CloudFile {
    shareId: number
    sharedWith: string
    ownedByUsername: string
    permissions: string
}

interface FolderContents {
    folders: Folder[]
    files: CloudFile[]
}

interface SharedContents {
    folders: SharedFolder[]
    files: SharedFile[]
}

interface RegisterRequest {
    id: number
    username: string
    status: 'Pending' | 'Rejected'
}

interface PathSegment {
  id: number;
  displayName: string;
}

export type { User, Folder, CloudFile, FolderContents, RegisterRequest, PathSegment, SharedContents, SharedFile, SharedFolder, UserSummary }