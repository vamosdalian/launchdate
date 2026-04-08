import type { LL2LaunchNormal } from './ll2-launch';
import type { LL2AgencyDetailed } from './ll2-agency';
import type { LL2LauncherConfigNormal, LL2LauncherConfigFamilyDetailed } from './ll2-launcher';
import type { LL2LocationSerializerWithPads, LL2Pad } from './ll2-location';

export interface LL2LaunchList {
  count: number;
  launches: LL2LaunchNormal[];
}

export interface LL2AgencyList {
  count: number;
  agencies: LL2AgencyDetailed[];
}

export interface LL2LauncherList {
  count: number;
  launchers: LL2LauncherConfigNormal[];
}

export interface LL2LauncherFamilyList {
  count: number;
  families: LL2LauncherConfigFamilyDetailed[];
}

export interface LL2LocationList {
  count: number;
  locations: LL2LocationSerializerWithPads[];
}

export interface LL2PadList {
  count: number;
  pads: LL2Pad[];
}
