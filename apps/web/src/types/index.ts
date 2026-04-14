export interface RocketListItem {
  id: string;
  name: string;
  thumb_image: string;
}

export type PublicLaunchStatus =
  | 'scheduled'
  | 'success'
  | 'failure'
  | 'cancelled'
  | 'delayed'
  | 'in_flight'
  | 'unknown';

export interface PublicCompactRocket {
  id: string;
  name: string;
  image_url: string;
  thumb_image: string;
}

export interface PublicCompactCompany {
  id: string;
  name: string;
  image_url: string;
}

export interface PublicCompactLaunchBase {
  id: string;
  name: string;
  location: string;
  country: string;
  latitude: number;
  longitude: number;
}

export interface PublicLaunchSummary {
  id: string;
  name: string;
  launch_time: string;
  status: PublicLaunchStatus;
  status_label: string;
  thumb_image: string;
  background_image: string;
  rocket: PublicCompactRocket;
  company: PublicCompactCompany;
  launch_base: PublicCompactLaunchBase;
}

export interface RocketDetail {
  id: string;
  name: string;
  description: string;
  active: boolean;
  reusable: boolean;
  launch_image: string;
  main_image: string;
  image_list: string[];
  company: PublicCompactCompany;
  launches: PublicLaunchSummary[];
  launch_cost: number;
  diameter: number;
  length: number;
  liftoff_thrust: number;
  launch_mass: number;
  leo_capacity: number;
  gto_capacity: number;
  geo_capacity: number;
  sso_capacity: number;
  total_launches: number;
  success_launches: number;
  failure_launches: number;
  total_landings: number;
  success_landings: number;
  failure_landings: number;
}

export type Rocket = RocketDetail; // Alias for backward compatibility in RocketDetail.tsx

export interface Launch {
  id: number;
  name: string;
  // Legacy date field for backward compatibility
  date?: string;
  rocket_id?: number;
  rocket?: string;
  launch_base_id?: number;
  launchBase?: string;
  status: 'scheduled' | 'successful' | 'failed' | 'cancelled';
  description?: string;
  created_at?: string;
  updated_at?: string;
  
  // New fields from API changes
  cospar_id?: string;
  sort_date?: string;
  slug?: string;
  modified?: string;
  
  // Launch window fields
  // Note: API supports both window_* and win_* variants for compatibility
  // Prefer window_* fields, but support win_* for backward compatibility
  window_open?: string;
  win_open?: string;
  t0?: string;
  window_close?: string;
  win_close?: string;
  date_str?: string;
  
  // Provider information (nested)
  provider_id?: number;
  provider?: {
    id: number;
    name: string;
    slug: string;
  };
  
  // Vehicle information (nested)
  vehicle?: {
    id: number;
    name: string;
    company_id?: number;
    slug: string;
  };
  
  // Pad & Location (nested)
  pad?: {
    id: number;
    name: string;
    location?: {
      id: number;
      name: string;
      state: string;
      statename: string;
      country: string;
      slug: string;
    };
  };
  
  // Mission details
  mission_description?: string;
  launch_description?: string;
  missions?: Array<{
    id: number;
    name: string;
    description: string;
  }>;
  
  // Weather information
  weather_summary?: string;
  weather_temp?: number;
  weather_condition?: string;
  weather_wind_mph?: number;
  weather_icon?: string;
  weather_updated?: string;
  
  // Additional metadata
  tags?: Array<{
    id: number;
    text: string;
  }>;
  quicktext?: string;
  suborbital?: boolean;
}

export interface News {
  id: number;
  title: string;
  summary: string;
  content?: string; // Markdown content for the article
  date: string;
  url: string;
  imageUrl: string;
  created_at?: string;
  updated_at?: string;
}

export interface LaunchBase {
  id: string;
  name: string;
  location: string;
  country: string;
  description: string;
  image_url: string;
  latitude: number;
  longitude: number;
  launches?: PublicLaunchSummary[];
  stats?: {
    launch_count: number;
    upcoming_launch_count: number;
    successful_launches: number;
    failed_launches: number;
    success_rate: number;
  };
  created_at?: string;
  updated_at?: string;
}

export interface Company {
  id: string;
  name: string;
  description: string;
  founded: number;
  founder: string;
  headquarters: string;
  employees: number;
  website: string;
  image_url: string;
  rockets?: RocketListItem[];
  launches?: PublicLaunchSummary[];
  stats?: {
    rocket_count: number;
    launch_count: number;
    successful_launches: number;
    failed_launches: number;
    pending_launches: number;
  };
  created_at?: string;
  updated_at?: string;
}

export type { PageBackground, PageBackgroundKey } from './page-background';

export interface PublicLaunchPage {
  count: number;
  launches: PublicLaunchSummary[];
}

export interface PublicRocketPage {
  count: number;
  rockets: RocketListItem[];
}

export interface PublicCompanyPage {
  count: number;
  companies: Company[];
}

export interface PublicLaunchBasePage {
  count: number;
  launch_bases: LaunchBase[];
}

export interface TimelineEvent {
  relative_time: string;
  abbrev: string;
  description: string;
}

export interface Mission {
  name: string;
  description: string;
}

export interface PublicLaunchView {
  id: string;
  name: string;
  launch_time: string;
  status: PublicLaunchStatus;
  status_label: string;
  background_image: string;
  image_list: string[];
  rocket: PublicCompactRocket;
  company: PublicCompactCompany;
  launch_base: PublicCompactLaunchBase;
  missions: Mission[];
  timeline: TimelineEvent[];
}
