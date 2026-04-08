import { apiFetch } from '../lib/api';
import type { PublicLaunchList, PublicLaunchDetail } from '../types';

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export async function fetchRocketLaunches(page: number = 0): Promise<PublicLaunchList> {
  const response = await apiFetch<ApiResponse<PublicLaunchList>>(`/api/v1/launch?page=${page}`);
  return response.data;
}

export async function fetchRocketLaunch(id: string): Promise<PublicLaunchDetail> {
  const response = await apiFetch<ApiResponse<PublicLaunchDetail>>(`/api/v1/launch/${id}`);
  return response.data;
}
