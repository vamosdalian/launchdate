import { apiFetch } from '../lib/api';
import type { Company, PublicCompanyPage } from '../types';

interface ListQueryOptions {
  page?: number;
  search?: string;
  homepageOnly?: boolean;
}

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

function buildListQuery({ page = 0, search = '', homepageOnly = false }: ListQueryOptions = {}) {
  const params = new URLSearchParams({ page: String(page) });
  const trimmedSearch = search.trim();
  if (trimmedSearch) {
    params.set('search', trimmedSearch);
  }
  if (homepageOnly) {
    params.set('homepage_only', 'true');
  }
  return params.toString();
}

export async function fetchCompanies(options: ListQueryOptions = {}): Promise<PublicCompanyPage> {
  const response = await apiFetch<ApiResponse<PublicCompanyPage>>(`/api/v1/companies?${buildListQuery(options)}`);
  return response.data;
}

export async function fetchCompany(id: string): Promise<Company> {
  const response = await apiFetch<ApiResponse<Company>>(`/api/v1/companies/${id}`);
  return response.data;
}
