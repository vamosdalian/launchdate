import { clearAuthState, getAccessToken, setAccessToken } from './authStore';
import { API_BASE_URL } from './config';

type FetchOptions = RequestInit & { retry?: boolean };

const AUTH_EXEMPT_PATHS = new Set(['/api/v1/auth/refresh', '/api/v1/auth/logout', '/api/v1/auth/me']);

let isRefreshing = false;
let refreshQueue: Array<(success: boolean) => void> = [];

function normalizePath(input: RequestInfo): string | null {
  if (typeof input === 'string') {
    if (input.startsWith('http')) {
      try {
        return new URL(input).pathname;
      } catch {
        return null;
      }
    }
    const value = input.startsWith('/') ? input : `/${input}`;
    const endIndex = value.indexOf('?');
    return endIndex === -1 ? value : value.slice(0, endIndex);
  }
  if (input instanceof Request) {
    try {
      return new URL(input.url).pathname;
    } catch {
      return null;
    }
  }
  return null;
}

function shouldAttachAuth(input: RequestInfo): boolean {
  const path = normalizePath(input);
  if (!path) return true;
  return !AUTH_EXEMPT_PATHS.has(path);
}

function buildRequest(input: RequestInfo, init: FetchOptions = {}): { url: string; config: FetchOptions; authExempt: boolean } {
  if (typeof input !== 'string') {
    throw new Error('apiFetch expects string request paths');
  }

  const headers = new Headers(init.headers as HeadersInit | undefined);
  const authExempt = !shouldAttachAuth(input);
  if (!authExempt) {
    const token = getAccessToken();
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  return {
    url: `${API_BASE_URL}${input}`,
    config: {
      ...init,
      headers,
      credentials: 'include'
    },
    authExempt
  };
}

async function doRefresh(): Promise<void> {
  const url = `${API_BASE_URL}/api/v1/auth/refresh`;
  const res = await fetch(url, { method: 'POST', credentials: 'include' });
  if (!res.ok) throw new Error('refresh failed');
  const data = await res.json();
  if (!data?.accessToken) {
    throw new Error('refresh missing access token');
  }
  setAccessToken(data.accessToken);
}

async function waitForRefresh(): Promise<boolean> {
  return new Promise(resolve => {
    refreshQueue.push(resolve);
  });
}

function resolveQueue(success: boolean) {
  refreshQueue.forEach(cb => cb(success));
  refreshQueue = [];
}

export async function apiFetch(input: RequestInfo, init: FetchOptions = {}): Promise<Response> {
  const { url, config, authExempt } = buildRequest(input, init);
  const resp = await fetch(url, config);

  if (resp.status !== 401 || authExempt) return resp;

  if (!isRefreshing) {
    isRefreshing = true;
    try {
      await doRefresh();
      isRefreshing = false;
      resolveQueue(true);
      const retry = buildRequest(input, init);
      return fetch(retry.url, retry.config);
    } catch (err) {
      isRefreshing = false;
      resolveQueue(false);
      clearAuthState();
      window.location.href = '/login';
      throw err;
    }
  }

  const refreshSuccess = await waitForRefresh();
  if (refreshSuccess) {
    const retry = buildRequest(input, init);
    return fetch(retry.url, retry.config);
  }

  clearAuthState();
  throw new Error('authentication failed');
}
