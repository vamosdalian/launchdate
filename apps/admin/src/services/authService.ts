import { apiFetch } from '@/lib/api';
import { API_BASE_URL } from '@/lib/config';
import { clearAuthState, getAuthUser, setAccessToken, setAuthUser } from '@/lib/authStore';
import type { AuthUser } from '@/types/auth';

export async function login(username: string, password: string): Promise<{ accessToken?: string; user?: AuthUser }> {
  const url = `${API_BASE_URL}/api/v1/auth/login`;
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ username, password })
  });

  if (!response.ok) {
    throw new Error('Login failed');
  }

  const data = await response.json();
  setAccessToken(data.accessToken ?? null);
  setAuthUser(data.user ?? null);
  return data;
}

export async function logout(): Promise<void> {
  try {
    await apiFetch('/api/v1/auth/logout', { method: 'POST' });
  } finally {
    clearAuthState();
    window.location.href = '/login';
  }
}

export async function getCurrentUser(forceRemote = false): Promise<AuthUser | null> {
  const cached = getAuthUser();
  if (cached && !forceRemote) {
    return cached;
  }

  try {
    const response = await apiFetch('/api/v1/auth/me', { method: 'GET' });
    if (!response.ok) {
      setAuthUser(null);
      return null;
    }
    const data = await response.json();
    const user: AuthUser | null = data?.user ?? null;
    setAuthUser(user);
    return user;
  } catch {
    setAuthUser(null);
    return null;
  }
}

export async function refreshAccessToken(): Promise<string | null> {
  try {
    const response = await apiFetch('/api/v1/auth/refresh', { method: 'POST' });
    if (!response.ok) {
      return null;
    }
    const data = await response.json();
    if (data?.accessToken) {
      setAccessToken(data.accessToken);
      return data.accessToken;
    }
    return null;
  } catch {
    return null;
  }
}
