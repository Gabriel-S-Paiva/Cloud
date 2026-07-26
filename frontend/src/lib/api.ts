import { goto } from '$app/navigation';
import type { User, Folder, CloudFile, FolderContents, RegisterRequest } from '$lib/types';

const API_BASE = 'http://localhost:8080';

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include',
  });

  if (response.status === 401) {
    goto('/login');
    throw new Error('Session expired. Please log in again.');
  }

  const text = await response.text();

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    if (text) {
      try {
        message = JSON.parse(text).error ?? message;
      } catch {
        // body wasn't JSON — fall back to the generic message
      }
    }
    throw new Error(message);
  }

  return text ? JSON.parse(text) : (undefined as T);
}

export const endpoints = {
    register: (username: string, password: string) =>
        request<void>('/users', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        }),
    login: (username: string, password: string) => 
        request<void>('/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        }),
    logout: () => 
        request<void>('/logout', {
            method: 'POST'
        }),
    getMe: () => 
        request<User>('/users/me', {
            method: 'GET'
        }),
    createFolder: (displayName: string, parentFolder: number|null = null) => 
        request<{ id: number }>('/folders', {
            method: 'POST',
            body: JSON.stringify({ displayName, parentFolder })
        }),
    getFolder: (id: number) =>
        request<Folder>(`/folders/${id}`, {
            method: 'GET'
        }),
    getFolderContent: (id: number) => 
        request<FolderContents>(`/folders/${id}/content`, {
            method: 'GET'
        }),
    updateFolder: (id: number, displayName: string|null = null, parentFolder: number|null = null) =>
        request<void>(`/folders/${id}`, {
            method: 'PATCH',
            body: JSON.stringify({ displayName, parentFolder })
        }),
    deleteFolder: (id: number) =>
        request<void>(`/folders/${id}`, {
            method:'DELETE'
        }),
    createFile: (displayName: string, parentFolder: number|null = null, size: number, contentType: string) => 
        request<{ id: number }>('/files', {
            method:'POST',
            body: JSON.stringify({ displayName, parentFolder, size, contentType})
        }),
    uploadChunk: (id: number, chunk: Blob) =>
        request<void>(`/files/${id}/chunk`, {
            method: 'POST',
            body: chunk,
            headers: { 'Content-Type': 'application/octet-stream' }
        }),
    getFile: (id: number) =>
        request<CloudFile>(`/files/${id}`, {
            method: 'GET'
        }),
    getFileContent: async (id: number, download = false): Promise<Blob> => {
        const response = await fetch(`${API_BASE}/files/${id}/content${download ? '?download=true' : ''}`, {
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Failed to fetch file content');
            return response.blob();
        },
    updateFile: (id: number, displayName: string|null = null, parentFolder: number|null = null) =>
        request<void>(`/files/${id}`, {
            method: 'PATCH',
            body: JSON.stringify({displayName, parentFolder})
        }),
    deleteFile: (id: number) =>
        request<void>(`/files/${id}`, {
            method: 'DELETE'
        }),
    createShare: (fileId: number|null, folderId: number|null = null, sharedWith: number, permission: 'Edit'|'View') =>
        request<{ id: number }>('/shares', {
            method: 'POST',
            body: JSON.stringify({ fileId, folderId, sharedWith, permission })
        }),
    getShareIncoming: () => 
        request<FolderContents>('/shares/incoming', {
            method: 'GET'
        }),
    getShareOutgoing: () =>
        request<FolderContents>('/shares/outgoing', {
            method: 'GET'
        }),
    updateShare: (id: number, permission: 'Edit'|'View') => 
        request<void>(`/shares/${id}`, {
            method: 'PATCH',
            body: JSON.stringify({ permission })
        }),
    deleteShare: (id: number) =>
        request<void>(`/shares/${id}`, {
            method:'DELETE'
        }),
    getRegisterRequests: () =>
        request<RegisterRequest[]>('/users/requests', {
            method:'GET'
        }),
    aproveRequest: (id: number) =>
        request<{ id: number }>(`/users/requests/${id}/aprove`, {
            method:'POST'
        }),
    rejectRequest: (id: number) =>
        request<void>(`/users/requests/${id}/reject`, {
            method: 'POST'
        })
};