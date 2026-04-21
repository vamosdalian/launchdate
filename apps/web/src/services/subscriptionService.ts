import { apiPost } from '../lib/api';

export interface SubscribeRequest {
  email: string;
}

export interface SubscribeResponse {
  status: 'subscribed';
}

interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export async function subscribeEmail(payload: SubscribeRequest): Promise<SubscribeResponse> {
  const response = await apiPost<ApiResponse<SubscribeResponse>>('/api/v1/subscriptions', payload);
  return response.data;
}
