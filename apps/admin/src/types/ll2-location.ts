import type { LL2Image, LL2Country, LL2CelestialBodyDetailed } from "./ll2-common";
import type { LL2AgencyNormal } from "./ll2-agency";

export interface LL2PadSerializerNoLocation {
  id: number;
  url: string;
  active: boolean;
  agencies: LL2AgencyNormal[];
  name: string;
  image: LL2Image | null;
  description: string | null;
  info_url: string | null;
  wiki_url: string | null;
  map_url: string | null;
  latitude: number;
  longitude: number;
  country: LL2Country | null;
  map_image: string | null;
  total_launch_count: number;
  orbital_launch_attempt_count: number;
  fastest_turnaround: string | null;
}

export interface LL2Pad extends LL2PadSerializerNoLocation {
  location: LL2Location;
}

export interface LL2LocationResponse {
  count: number;
  next: string | null;
  previous: string | null;
  results: LL2LocationSerializerWithPads[];
}

export interface LL2Location {
  response_mode: string;
  id: number;
  name: string;
  url: string;
  celestial_body: LL2CelestialBodyDetailed | null;
  active: boolean;
  country: LL2Country | null;
  description: string | null;
  image: LL2Image | null;
  map_image: string | null;
  latitude: number;
  longitude: number;
  timezone_name: string;
  total_launch_count: number;
  total_landing_count: number;
}

export interface LL2LocationSerializerWithPads extends LL2Location {
  pads: LL2PadSerializerNoLocation[];
}
