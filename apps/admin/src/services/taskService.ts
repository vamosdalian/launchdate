import { apiClient } from './apiClient';

export type TaskType =
  | 'launch'
  | 'agency'
  | 'launcher'
  | 'launcher_family'
  | 'pad'
  | 'location'
  | 'upcoming';

export type TaskAction = 'pause' | 'resume' | 'cancel';
export type TaskStatus = 'idle' | 'running' | 'paused' | 'failed';

export interface TaskInfo {
  type: string;
  status: TaskStatus;
  progress?: Record<string, unknown>;
  last_error?: string;
  started_at: string;
  updated_at: string;
}

export const taskService = {
  getCurrentTask: async (): Promise<TaskInfo | null> => {
    const data = await apiClient.get<TaskInfo | null>('/api/v1/task');
    return data ?? null;
  },

  startTask: async (type: TaskType): Promise<void> => {
    await apiClient.post<unknown>('/api/v1/task', { type });
  },

  actionTask: async (action: TaskAction): Promise<void> => {
    await apiClient.post<unknown>('/api/v1/task/action', { action });
  },
};
