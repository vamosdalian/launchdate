export interface Rocket {
  id: number;
  external_id?: number;
  name: string;
  description: string;
  height: number;
  diameter: number;
  mass: number;
  company_id?: number;
  company?: string;
  imageUrl: string;
  active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface RocketLaunchProvider {
  id: number;
  name: string;
  slug: string;
}

export interface RocketLaunchVehicle {
  id: number;
  name: string;
  company_id?: number;
  slug: string;
}

export interface RocketLaunchPadLocation {
  id: number;
  name: string;
  state: string;
  statename: string;
  country: string;
  slug: string;
}

export interface RocketLaunchPad {
  id: number;
  name: string;
  location?: RocketLaunchPadLocation;
}

export interface RocketLaunchMission {
  id: number;
  external_id?: number;
  name: string;
  description: string;
}

export interface RocketLaunchTag {
  id: number;
  text: string;
}

export interface Launch {
  id: number;
  external_id?: number;
  cospar_id: string;
  sort_date: string;
  name: string;
  launch_date: string;
  description?: string;
  provider?: RocketLaunchProvider;
  provider_id?: number;
  vehicle?: RocketLaunchVehicle;
  rocket_id?: number;
  pad?: RocketLaunchPad;
  launch_base_id?: number;
  missions?: RocketLaunchMission[];
  mission_description: string;
  launch_description: string;
  win_open?: string;
  t0?: string;
  win_close?: string;
  date_str: string;
  tags?: RocketLaunchTag[];
  slug: string;
  weather_summary: string;
  weather_temp?: number;
  weather_condition: string;
  weather_wind_mph?: number;
  weather_icon: string;
  weather_updated?: string;
  quicktext: string;
  suborbital: boolean;
  modified?: string;
  status: 'scheduled' | 'successful' | 'failed' | 'cancelled';
  created_at?: string;
  updated_at?: string;
}

export interface News {
  id: number;
  title: string;
  summary: string;
  content?: string;
  date: string;
  url: string;
  image_url: string;
  created_at?: string;
  updated_at?: string;
}

export interface LaunchBase {
  id: number;
  external_id?: number;
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
  external_id?: number;
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

export interface LL2LaunchStatus {
  id: number;
  name: string;
  abbrev: string;
}

export interface LL2LaunchServiceProvider {
  id: number;
  name: string;
}

export interface LL2LaunchRocketConfiguration {
  name?: string;
}

export interface LL2LaunchRocket {
  configuration?: LL2LaunchRocketConfiguration;
}

export interface LL2LaunchPadLocation {
  name?: string;
}

export interface LL2LaunchPad {
  name?: string;
  location?: LL2LaunchPadLocation;
}

export interface LL2Launch {
  id: number;
  name: string;
  net: string;
  status?: LL2LaunchStatus;
  launch_service_provider?: LL2LaunchServiceProvider;
  rocket?: LL2LaunchRocket;
  pad?: LL2LaunchPad;
}

type LL2KeyedValue = {
  id?: number | string;
  name?: string;
};

export interface LL2Agency {
  id: number;
  name: string;
  country_code?: string | LL2KeyedValue | null;
  type?: string | LL2KeyedValue | null;
  description?: string;
  launch_count?: number;
  successful_launches?: number;
}

export interface LL2PaginatedResult<T> {
  count: number;
  items: T[];
}

export interface LL2LocationPadReference {
  id?: number;
  name?: string;
}

export interface LL2Location {
  id: number;
  name: string;
  country_code?: string | null;
  pads?: LL2LocationPadReference[];
}

export interface LL2PadLocation {
  id?: number;
  name?: string;
  country_code?: string | null;
}

export interface LL2Pad {
  id: number;
  name: string;
  location?: LL2PadLocation;
  latitude?: number;
  longitude?: number;
  map_url?: string;
}

export interface LL2LauncherFamilyReference {
  id?: number;
  name?: string;
}

export interface LL2LauncherManufacturer {
  id?: number;
  name?: string;
}

export interface LL2Launcher {
  id: number;
  name: string;
  full_name?: string;
  variant?: string;
  family?: LL2LauncherFamilyReference;
  manufacturer?: LL2LauncherManufacturer;
}

export interface LL2LauncherFamilyOrbit {
  id?: number;
  name?: string;
}

export interface LL2LauncherFamily {
  id: number;
  name: string;
  description?: string;
  orbit?: LL2LauncherFamilyOrbit;
}
