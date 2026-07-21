interface User {
    id: number
    username: string
    role: 'User' | 'Admin'
    quota: number
    quotaUsed: number
}

interface Folder {
    id: number,
    displayName: string
    ownedBy: number
    parentFolder: number|null
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

interface FolderContents {
    folders: Folder[]
    files: CloudFile[]
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

export type { User, Folder, CloudFile, FolderContents, RegisterRequest, PathSegment}