export interface RocketListItem {
  id: string;
  name: string;
  thumb_image: string;
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
  agency_info?: PublicCompactAgency;
  launches: PublicCompactLaunch[]; 
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
  id: number;
  name: string;
  location: string;
  country: string;
  description: string;
  imageUrl: string;
  latitude: number;
  longitude: number;
  created_at?: string;
  updated_at?: string;
}

export interface Company {
  id: number;
  name: string;
  description: string;
  founded: number;
  founder: string;
  headquarters: string;
  employees: number;
  website: string;
  imageUrl: string;
  created_at?: string;
  updated_at?: string;
}

// New API Types
export interface PublicCompactLaunch {
  id: string;
  name: string;
  launch_time: string;
  status: number;
  thumb_image: string;
  rocket_name: string;
  agency_name: string;
  location: string;
}

export interface PublicLaunchList {
  count: number;
  launches: PublicCompactLaunch[];
}

export interface PublicCompactRocket {
  id: string;
  name: string;
  thumb_image: string;
}

export interface PublicCompactAgency {
  id: string;
  name: string;
  thumb_image: string;
}

export interface PublicCompactLocation {
  id: string;
  name: string;
  lat: number;
  lon: number;
}

export interface TimelineEvent {
  relative_time: string;
  abbrev: string;
  description: string;
}

export interface Mission {
  id: string;
  name: string;
  description: string;
}

export interface PublicLaunchDetail {
  id: string;
  name: string;
  launch_time: string;
  status: number;
  background_image: string;
  image_list: string[];
  rocket_info: PublicCompactRocket;
  agency_info: PublicCompactAgency;
  location_info: PublicCompactLocation;
  mission_info: Mission[];
  timeline_event: TimelineEvent[];
}
