import { apiFetch } from '@/lib/api';

export interface ApiError {
  message: string;
  status?: number;
}

export interface ApiEnvelope<T> {
  code: number;
  message: string;
  data: T;
}

class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers = new Headers(options.headers as HeadersInit | undefined);
    const method = (options.method ?? 'GET').toUpperCase();
    if (method !== 'GET' && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json');
    }

    const config: RequestInit = {
      ...options,
      headers,
    };

    try {
      const response = await apiFetch(endpoint, config);

      if (!response.ok) {
        const error: ApiError = {
          message: `HTTP error! status: ${response.status}`,
          status: response.status,
        };
        throw error;
      }

      // Handle empty responses (e.g., 204 No Content for DELETE)
      const contentType = response.headers.get('content-type');
      if (response.status === 204 || !contentType?.includes('application/json')) {
        return {} as T;
      }

      const result = await response.json() as ApiEnvelope<T>;

      if (result.code !== 0) {
        throw new Error(result.message || 'Request failed');
      }

      return result.data;
    } catch (error) {
      // Re-throw ApiError instances
      if (error && typeof error === 'object' && 'status' in error) {
        throw error;
      }
      // Wrap other errors in ApiError format
      const apiError: ApiError = {
        message: error instanceof Error ? error.message : 'Network error',
      };
      throw apiError;
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}
export const apiClient = new ApiClient();
