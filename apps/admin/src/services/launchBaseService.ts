import { apiClient } from './apiClient';
import type { LaunchBaseList } from '@/types/launch-base';
import type { LL2LocationList, LL2PadList } from '@/types/ll2-responses';

const normalizeFilterValue = (value?: string | null) => {
  if (typeof value !== 'string') {
    return undefined;
  }

  const trimmed = value.trim();
  return trimmed.length > 0 ? trimmed : undefined;
};

export type LaunchBaseFilters = {
  name?: string;
  celestialBody?: string;
  country?: string;
  sortBy?: 'name';
  sortOrder?: 'asc' | 'desc';
};

export const launchBaseService = {
  getProdLaunchBases: async (
    limit = 50,
    offset = 0,
    filters: LaunchBaseFilters = {},
  ): Promise<LaunchBaseList> => {
    const params = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });

    const name = normalizeFilterValue(filters.name);
    if (name) {
      params.set('name', name);
    }

    const celestialBody = normalizeFilterValue(filters.celestialBody);
    if (celestialBody) {
      params.set('celestial_body', celestialBody);
    }

    const country = normalizeFilterValue(filters.country);
    if (country) {
      params.set('country', country);
    }

    if (filters.sortBy === 'name') {
      params.set('sort_by', 'name');
    }

    if (filters.sortOrder) {
      params.set('sort_order', filters.sortOrder);
    }

    return apiClient.get<LaunchBaseList>(`/api/v1/data/launchbases?${params.toString()}`);
  },

  getLL2Locations: async (limit = 50, offset = 0): Promise<LL2LocationList> => {
    return apiClient.get<LL2LocationList>(`/api/v1/ll2/locations?limit=${limit}&offset=${offset}`);
  },

  getLL2Pads: async (limit = 50, offset = 0): Promise<LL2PadList> => {
    return apiClient.get<LL2PadList>(`/api/v1/ll2/pads?limit=${limit}&offset=${offset}`);
  },

};
