import type { LL2Image } from "./ll2-common";
import type { LL2AgencyNormal, LL2AgencyMini } from "./ll2-agency";

export interface LL2LauncherConfigList {
  response_mode: string;
  id: number;
  url: string;
  name: string;
  families: LL2LauncherConfigFamilyMini[];
  full_name: string;
  variant: string;
}

export interface LL2LauncherConfigNormal extends LL2LauncherConfigList {
  active: boolean;
  is_placeholder: boolean;
  manufacturer: LL2AgencyNormal;
  program: LL2ProgramNormal[];
  reusable: boolean;
  image: LL2Image | null;
  info_url: string | null;
  wiki_url: string | null;
}

export interface LL2LauncherConfigDetailed extends LL2LauncherConfigNormal {
  description: string | null;
  alias: string | null;
  min_stage: number | null;
  max_stage: number | null;
  length: number | null;
  diameter: number | null;
  maiden_flight: string | null;
  launch_cost: number | null;
  launch_mass: number | null;
  leo_capacity: number | null;
  gto_capacity: number | null;
  geo_capacity: number | null;
  sso_capacity: number | null;
  to_thrust: number | null;
  apogee: number | null;
  total_launch_count: number;
  consecutive_successful_launches: number;
  successful_launches: number;
  failed_launches: number;
  pending_launches: number;
  attempted_landings: number;
  successful_landings: number;
  failed_landings: number;
  consecutive_successful_landings: number;
  fastest_turnaround: string | null;
}

export interface LL2LauncherFamilyResponse {
  count: number;
  next: string | null;
  previous: string | null;
  results: LL2LauncherConfigFamilyDetailed[];
}

export interface LL2LauncherConfigFamilyMini {
  response_mode: string;
  id: number;
  name: string;
}

export interface LL2LauncherConfigFamilyNormal
  extends LL2LauncherConfigFamilyMini {
  manufacturer: LL2AgencyNormal[];
  parent: LL2LauncherConfigFamilyMini | null;
}

export interface LL2LauncherConfigFamilyDetailed
  extends LL2LauncherConfigFamilyNormal {
  description: string | null;
  active: boolean;
  maiden_flight: string | null;
  total_launch_count: number;
  consecutive_successful_launches: number;
  successful_launches: number;
  failed_launches: number;
  pending_launches: number;
  attempted_landings: number;
  successful_landings: number;
  failed_landings: number;
  consecutive_successful_landings: number;
}

export interface LL2ProgramType {
  id: number;
  name: string;
}

export interface LL2ProgramNormal {
  response_mode: string;
  id: number;
  url: string;
  name: string;
  image: LL2Image | null;
  info_url: string | null;
  wiki_url: string | null;
  description: string | null;
  agencies: LL2AgencyMini[];
  start_date: string | null;
  end_date: string | null;
  mission_patches: LL2MissionPatch[];
  type: LL2ProgramType | null;
}

export interface LL2MissionPatch {
  id: number;
  name: string;
  priority: number;
  image_url: string;
  agency: LL2AgencyMini;
  response_mode: string;
}
