import { apiClient } from './apiClient';
import type { News } from '@/types';

export const newsService = {
  getAll: async (): Promise<News[]> => {
    return apiClient.get<News[]>('/api/v1/news');
  },

  getById: async (id: number): Promise<News> => {
    return apiClient.get<News>(`/api/v1/news/${id}`);
  },

  create: async (news: Omit<News, 'id' | 'created_at' | 'updated_at'>): Promise<News> => {
    return apiClient.post<News>('/api/v1/news', news);
  },

  update: async (id: number, news: Partial<Omit<News, 'id' | 'created_at' | 'updated_at'>>): Promise<{ message: string }> => {
    return apiClient.put<{ message: string }>(`/api/v1/news/${id}`, news);
  },

  delete: async (id: number): Promise<void> => {
    return apiClient.delete<void>(`/api/v1/news/${id}`);
  },
};
