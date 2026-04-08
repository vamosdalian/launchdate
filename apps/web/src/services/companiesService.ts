import { apiFetch } from '../lib/api';
import type { Company } from '../types';

export async function fetchCompanies(): Promise<Company[]> {
  return apiFetch<Company[]>('/api/v1/companies');
}

export async function fetchCompany(id: string): Promise<Company> {
  return apiFetch<Company>(`/api/v1/companies/${id}`);
}
