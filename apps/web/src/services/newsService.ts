import { apiFetch } from '../lib/api';
import type { News } from '../types';

export async function fetchNews(): Promise<News[]> {
  return apiFetch<News[]>('/api/v1/news');
}

export async function fetchNewsArticle(id: string): Promise<News> {
  return apiFetch<News>(`/api/v1/news/${id}`);
}
