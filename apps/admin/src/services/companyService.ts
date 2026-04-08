import { apiClient } from './apiClient';
import type { AgencyList } from '@/types/agency';
import type { LL2AgencyList } from '@/types/ll2-responses';

const normalizeFilterValue = (value?: string | null) => {
  if (typeof value !== 'string') {
    return undefined;
  }

  const trimmed = value.trim();
  return trimmed.length > 0 ? trimmed : undefined;
};

export type CompanyFilters = {
  name?: string;
  type?: string;
  country?: string;
  sortBy?: 'founding_year';
  sortOrder?: 'asc' | 'desc';
};

export const companyService = {
  getProdAgencies: async (
    limit = 50,
    offset = 0,
    filters: CompanyFilters = {},
  ): Promise<AgencyList> => {
    const params = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });

    const name = normalizeFilterValue(filters.name);
    if (name) {
      params.set('name', name);
    }

    const type = normalizeFilterValue(filters.type);
    if (type) {
      params.set('type', type);
    }

    const country = normalizeFilterValue(filters.country);
    if (country) {
      params.set('country', country);
    }

    if (filters.sortBy === 'founding_year') {
      params.set('sort_by', 'founding_year');
    }

    if (filters.sortOrder) {
      params.set('sort_order', filters.sortOrder);
    }

    return apiClient.get<AgencyList>(`/api/v1/data/agencies?${params.toString()}`);
  },

  updateAgency: async (id: number | string, data: Partial<AgencyList['agencies'][0]>) => {
    return apiClient.put<{ message: string }>(`/api/v1/data/agencies/${id}`, data);
  },

  getLL2Agencies: async (limit = 50, offset = 0): Promise<LL2AgencyList> => {
    return apiClient.get<LL2AgencyList>(`/api/v1/ll2/angecies?limit=${limit}&offset=${offset}`);
  },

};
