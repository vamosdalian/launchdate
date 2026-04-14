import { apiFetch } from '../lib/api';
import type { PublicRocketPage, RocketDetail } from '../types';

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

export async function fetchRockets(options: ListQueryOptions = {}): Promise<PublicRocketPage> {
  const response = await apiFetch<ApiResponse<PublicRocketPage>>(`/api/v1/rocket?${buildListQuery(options)}`);
  return response.data;
}

export async function fetchRocket(id: string): Promise<RocketDetail> {
  const response = await apiFetch<ApiResponse<RocketDetail>>(`/api/v1/rocket/${id}`);
  return response.data;
}
