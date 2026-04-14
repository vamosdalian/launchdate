import { apiFetch } from '../lib/api';
import type { PublicLaunchPage, PublicLaunchView } from '../types';

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

export async function fetchRocketLaunches(options: ListQueryOptions = {}): Promise<PublicLaunchPage> {
  const response = await apiFetch<ApiResponse<PublicLaunchPage>>(`/api/v1/launch?${buildListQuery(options)}`);
  return response.data;
}

export async function fetchRocketLaunch(id: string): Promise<PublicLaunchView> {
  const response = await apiFetch<ApiResponse<PublicLaunchView>>(`/api/v1/launch/${id}`);
  return response.data;
}
