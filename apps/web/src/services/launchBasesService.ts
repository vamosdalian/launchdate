import { apiFetch } from '../lib/api';
import type { LaunchBase } from '../types';

export async function fetchLaunchBases(): Promise<LaunchBase[]> {
  return apiFetch<LaunchBase[]>('/api/v1/launch-bases');
}

export async function fetchLaunchBase(id: string): Promise<LaunchBase> {
  return apiFetch<LaunchBase>(`/api/v1/launch-bases/${id}`);
}
