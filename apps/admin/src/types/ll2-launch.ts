import type {
  LL2Image,
  LL2InfoURL,
  LL2VidURL,
  LL2CelestialBodyMini,
} from "./ll2-common";
import type { LL2AgencyDetailed } from "./ll2-agency";
import type { LL2Pad } from "./ll2-location";
import type { LL2LauncherConfigList, LL2ProgramNormal } from "./ll2-launcher";
import type { LL2AgencyMini } from "./ll2-agency";

export interface LL2LaunchDetailed extends LL2LaunchNormal {
  flightclub_url: string | null;
  updates: LL2Update[];
  info_urls: LL2InfoURL[];
  vid_urls: LL2VidURL[];
  timeline: LL2TimelineEvent[];
  pad_turnaround: string | null;
  mission_patches: LL2MissionPatch[];
}

export interface LL2Update {
  id: number;
  profile_image: string | null;
  comment: string | null;
  info_url: string | null;
  created_by: string | null;
  created_on: string | null;
}

export interface LL2TimelineEvent {
  relative_time: string | null;
  type: LL2TimelineEventType | null;
}

export interface LL2TimelineEventType {
  id: number;
  abbrev: string;
  description: string | null;
}

export interface LL2LaunchNormal extends LL2LaunchBasic {
  probability: number | null;
  weather_concerns: string | null;
  failreason: string | null;
  hashtag: string | null;
  launch_service_provider: LL2AgencyMini;
  rocket: LL2RocketNormal;
  mission: LL2Mission | null;
  pad: LL2Pad;
  webcast_live: boolean;
  program: LL2ProgramNormal[];
  orbital_launch_attempt_count: number;
  location_launch_attempt_count: number;
  pad_launch_attempt_count: number;
  agency_launch_attempt_count: number;
  orbital_launch_attempt_count_year: number;
  location_launch_attempt_count_year: number;
  pad_launch_attempt_count_year: number;
  agency_launch_attempt_count_year: number;
}

export interface LL2LaunchBasic {
  id: string;
  url: string;
  name: string;
  response_mode: string;
  slug: string;
  launch_designator: string | null;
  status: LL2Status;
  last_updated: string;
  net: string;
  net_precision: LL2NetPrecision | null;
  window_end: string;
  window_start: string;
  image: LL2Image | null;
  infographic: string | null;
}

export interface LL2Status {
  id: number;
  name: string;
  abbrev: string;
  description: string;
}

export interface LL2NetPrecision {
  id: number;
  name: string;
  abbrev: string;
  description: string;
}

export interface LL2RocketNormal {
  id: number;
  configuration: LL2LauncherConfigList;
}

export interface LL2Mission {
  id: number;
  name: string;
  description: string;
  type: string;
  image: LL2Image | null;
  orbit: LL2Orbit | null;
  agencies: LL2AgencyDetailed[];
  info_urls: LL2InfoURL[];
  vid_urls: LL2VidURL[];
}

export interface LL2Orbit {
  id: number;
  name: string;
  abbrev: string;
  celestial_body: LL2CelestialBodyMini | null;
}

export interface LL2MissionPatch {
  id: number;
  name: string;
  priority: number;
  image_url: string;
  agency: LL2AgencyMini;
  response_mode: string;
}
