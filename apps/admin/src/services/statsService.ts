import { apiClient } from './apiClient';

interface LL2Stats {
  launches: number;
  agencies: number;
  launchers: number;
  launcher_families: number;
  locations: number;
  pads: number;
}

interface ApiStats {
  rocket: number;
  launch: number;
  agency: number;
  launch_base: number;
  ll2: LL2Stats;
}

export interface DashboardStats {
  totalRockets: number;
  upcomingLaunches: number;
  newsArticles: number;
  launchBases: number;
  companies: number;
  ll2: LL2Stats;
}

export const statsService = {
  getStats: async (): Promise<DashboardStats> => {
    const apiData = await apiClient.get<ApiStats>('/api/v1/data/stats');

    // The new API doesn't provide newsArticles, so we default it to 0.
    return {
      totalRockets: apiData.rocket,
      upcomingLaunches: apiData.launch,
      companies: apiData.agency,
      launchBases: apiData.launch_base,
      newsArticles: 0,
      ll2: apiData.ll2,
    };
  },
};
