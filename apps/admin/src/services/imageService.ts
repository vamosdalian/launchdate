import { apiClient } from './apiClient';
import type { ApiEnvelope } from './apiClient';
import { apiFetch } from '@/lib/api';
import type { ImageListResponse, GenerateThumbnailParams } from '@/types/image';

export const uploadImage = async (file: File): Promise<{ key: string }> => {
  const formData = new FormData();
  formData.append('file', file);

  const response = await apiFetch('/api/v1/images', {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
      throw new Error(`Upload failed: ${response.statusText}`);
  }
  const result = await response.json() as ApiEnvelope<{ key: string }>;
  
  if (result.code !== 0) {
    throw new Error(result.message || 'Request failed');
  }
  return result.data;
};

export const getImages = async (limit: number = 10, offset: number = 0): Promise<ImageListResponse> => {
  const params = new URLSearchParams({
    limit: limit.toString(),
    offset: offset.toString(),
  });
  return apiClient.get<ImageListResponse>(`/api/v1/images?${params.toString()}`);
};

export const deleteImage = async (key: string): Promise<void> => {
  await apiClient.delete<void>(`/api/v1/images/${key}`);
};

export const generateThumbnail = async (params: GenerateThumbnailParams): Promise<void> => {
  await apiClient.post<void>('/api/v1/images/thumb', params);
};
