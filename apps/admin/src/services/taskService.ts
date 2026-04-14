import { apiClient } from './apiClient';

export type TaskType =
  | 'launch'
  | 'agency'
  | 'launcher'
  | 'launcher_family'
  | 'pad'
  | 'location'
  | 'update';

export type TaskAction = 'pause' | 'resume' | 'cancel';
export type TaskStatus = 'idle' | 'running' | 'paused' | 'completed' | 'canceled' | 'failed';

export interface TaskProgress {
  current_count?: number;
  total_count?: number;
  overlap_seconds?: number;
  watermark_last_updated?: string;
  current_window_start?: string;
  current_window_end?: string;
  current_offset?: number;
  next_run_at?: string;
  last_success_at?: string;
  [key: string]: unknown;
}

export interface TaskInfo {
  type: string;
  status: TaskStatus;
  progress?: TaskProgress;
  last_error?: string;
  started_at: string;
  updated_at: string;
  finished_at?: string;
}

export const taskService = {
  getCurrentTask: async (): Promise<TaskInfo | null> => {
    const data = await apiClient.get<TaskInfo | null>('/api/v1/task');
    return data ?? null;
  },

  getTaskHistory: async (limit = 10): Promise<TaskInfo[]> => {
    const data = await apiClient.get<TaskInfo[]>(`/api/v1/task/history?limit=${limit}`);
    return data ?? [];
  },

  startTask: async (type: TaskType): Promise<TaskInfo | null> => {
    return apiClient.post<TaskInfo | null>('/api/v1/task', { type });
  },

  actionTask: async (action: TaskAction): Promise<TaskInfo | null> => {
    return apiClient.post<TaskInfo | null>('/api/v1/task/action', { action });
  },
};
