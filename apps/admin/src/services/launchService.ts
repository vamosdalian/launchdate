import { apiClient } from './apiClient';
import type { LaunchSerializer, LaunchList } from '@/types/launch';
import type { LL2LaunchList } from '@/types/ll2-responses';

export interface LaunchFilters {
  name?: string;
  status?: string;
  provider?: string;
  rocket?: string;
  mission?: string;
  sortBy?: 'time' | 'name';
  sortOrder?: 'asc' | 'desc';
}

export const launchService = {
  getProdLaunches: async (
    limit = 50,
    offset = 0,
    filters: LaunchFilters = {},
  ): Promise<LaunchList> => {
    const params = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });

    const normalizedText = (value?: string) => {
      if (typeof value !== 'string') {
        return undefined;
      }
      const trimmed = value.trim();
      return trimmed.length > 0 ? trimmed : undefined;
    };

    const nameFilter = normalizedText(filters.name);
    const statusFilter = normalizedText(filters.status);
    const providerFilter = normalizedText(filters.provider);
    const rocketFilter = normalizedText(filters.rocket);
    const missionFilter = normalizedText(filters.mission);

    if (nameFilter) params.set('name', nameFilter);
    if (statusFilter) params.set('status', statusFilter);
    if (providerFilter) params.set('launch_service_provider', providerFilter);
    if (rocketFilter) params.set('rocket', rocketFilter);
    if (missionFilter) params.set('mission', missionFilter);

    const sortBy = filters.sortBy && ['time', 'name'].includes(filters.sortBy) ? filters.sortBy : undefined;
    const sortOrder = filters.sortOrder && ['asc', 'desc'].includes(filters.sortOrder) ? filters.sortOrder : undefined;

    if (sortBy) params.set('sort_by', sortBy);
    if (sortOrder) params.set('sort_order', sortOrder);

    return apiClient.get<LaunchList>(`/api/v1/data/launches?${params.toString()}`);
  },

  getProdLaunch: async (id: string | number): Promise<LaunchSerializer> => {
    return apiClient.get<LaunchSerializer>(`/api/v1/data/launches/${id}`);
  },

  updateProdLaunch: async (id: string | number, data: { background_image?: string; image_list?: string[]; thumb_image?: string }): Promise<LaunchSerializer> => {
    return apiClient.put<LaunchSerializer>(`/api/v1/data/launches/${id}`, data);
  },

  getLL2Launches: async (limit = 50, offset = 0): Promise<LL2LaunchList> => {
    return apiClient.get<LL2LaunchList>(`/api/v1/ll2/launches?limit=${limit}&offset=${offset}`);
  },

};
