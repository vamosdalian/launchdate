import { apiFetch } from '../lib/api';
import type { RocketDetail, RocketListItem } from '../types';

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

interface RocketListResponse {
  count: number;
  rockets: RocketListItem[];
}

export async function fetchRockets(): Promise<RocketListItem[]> {
  const response = await apiFetch<ApiResponse<RocketListResponse>>('/api/v1/rocket');
  return response.data.rockets;
}

export async function fetchRocket(id: string): Promise<RocketDetail> {
  const response = await apiFetch<ApiResponse<RocketDetail>>(`/api/v1/rocket/${id}`);
  return response.data;
}
