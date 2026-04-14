import { apiFetch } from '../lib/api';
import type { PageBackground } from '../types';

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export async function fetchPageBackgrounds(): Promise<PageBackground[]> {
  const response = await apiFetch<ApiResponse<PageBackground[]>>('/api/v1/page-backgrounds');
  return response.data;
}