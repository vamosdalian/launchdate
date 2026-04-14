import { apiClient } from './apiClient';
import type { PageBackground, UpdatePageBackgroundPayload } from '@/types/page-background';

export const pageBackgroundService = {
  getPageBackgrounds(): Promise<PageBackground[]> {
    return apiClient.get<PageBackground[]>('/api/v1/data/page-backgrounds');
  },

  updatePageBackground(pageKey: string, payload: UpdatePageBackgroundPayload): Promise<PageBackground> {
    return apiClient.put<PageBackground>(`/api/v1/data/page-backgrounds/${pageKey}`, payload);
  },
};