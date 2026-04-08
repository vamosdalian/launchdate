import { apiClient } from './apiClient';
import type { Rocket, RocketList } from '@/types/rocket';
import type { LL2LauncherList, LL2LauncherFamilyList } from '@/types/ll2-responses';

export interface RocketFilters {
  fullName?: string;
  name?: string;
  variant?: string;
  sortBy?: 'full_name';
  sortOrder?: 'asc' | 'desc';
}

export const rocketService = {
  getProdRockets: async (
    limit = 50,
    offset = 0,
    filters: RocketFilters = {},
  ): Promise<RocketList> => {
    const params = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });

    const normalize = (value?: string) => {
      if (typeof value !== 'string') {
        return undefined;
      }
      const trimmed = value.trim();
      return trimmed.length > 0 ? trimmed : undefined;
    };

    const fullName = normalize(filters.fullName);
    const name = normalize(filters.name);
    const variant = normalize(filters.variant);

    if (fullName) params.set('full_name', fullName);
    if (name) params.set('name', name);
    if (variant) params.set('variant', variant);

    if (filters.sortBy === 'full_name') {
      params.set('sort_by', 'full_name');
    }

    if (filters.sortOrder && ['asc', 'desc'].includes(filters.sortOrder)) {
      params.set('sort_order', filters.sortOrder);
    }

    return apiClient.get<RocketList>(`/api/v1/data/rockets?${params.toString()}`);
  },

  updateRocket: async (id: string, data: Partial<Rocket>): Promise<void> => {
    await apiClient.put(`/api/v1/data/rockets/${id}`, data);
  },

  getLL2Launchers: async (limit = 50, offset = 0): Promise<LL2LauncherList> => {
    return apiClient.get<LL2LauncherList>(`/api/v1/ll2/launchers?limit=${limit}&offset=${offset}`);
  },

  getLL2LauncherFamilies: async (
    limit = 50,
    offset = 0
  ): Promise<LL2LauncherFamilyList> => {
    return apiClient.get<LL2LauncherFamilyList>(`/api/v1/ll2/launcher-families?limit=${limit}&offset=${offset}`);
  },

};
