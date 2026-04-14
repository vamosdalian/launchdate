// API configuration and base URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

interface ApiEnvelope {
  code: number;
  message: string;
  data?: unknown;
}

function isApiEnvelope(value: unknown): value is ApiEnvelope {
  return (
    typeof value === 'object' &&
    value !== null &&
    'code' in value &&
    typeof (value as { code: unknown }).code === 'number' &&
    'message' in value &&
    typeof (value as { message: unknown }).message === 'string'
  );
}

// Generic API fetch wrapper
async function apiFetch<T>(endpoint: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`);

  const payload: unknown = await response.json();

  if (!response.ok) {
    if (isApiEnvelope(payload)) {
      throw new Error(payload.message);
    }

    throw new Error(`API request failed: ${response.status} ${response.statusText}`);
  }

  if (isApiEnvelope(payload) && payload.code !== 0) {
    throw new Error(payload.message);
  }

  return payload as T;
}

export { API_BASE_URL, apiFetch };
