import type { AuthUser } from '@/types/auth';

let accessToken: string | null = null;
let currentUser: AuthUser | null = null;

export function setAccessToken(token: string | null): void {
  accessToken = token ?? null;
}

export function getAccessToken(): string | null {
  return accessToken;
}

export function setAuthUser(user: AuthUser | null): void {
  currentUser = user ?? null;
}

export function getAuthUser(): AuthUser | null {
  return currentUser;
}

export function clearAuthState(): void {
  accessToken = null;
  currentUser = null;
}
