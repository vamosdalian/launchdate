import { apiFetch } from '../lib/api';
import type { LaunchBase, PublicLaunchBasePage } from '../types';

interface ListQueryOptions {
  page?: number;
  search?: string;
}

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

function buildListQuery({ page = 0, search = '' }: ListQueryOptions = {}) {
  const params = new URLSearchParams({ page: String(page) });
  const trimmedSearch = search.trim();
  if (trimmedSearch) {
    params.set('search', trimmedSearch);
  }
  return params.toString();
}

export async function fetchLaunchBases(options: ListQueryOptions = {}): Promise<PublicLaunchBasePage> {
  const response = await apiFetch<ApiResponse<PublicLaunchBasePage>>(`/api/v1/launch-bases?${buildListQuery(options)}`);
  return response.data;
}

export async function fetchLaunchBase(id: string): Promise<LaunchBase> {
  const response = await apiFetch<ApiResponse<LaunchBase>>(`/api/v1/launch-bases/${id}`);
  return response.data;
}
