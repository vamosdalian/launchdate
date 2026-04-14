import type { LL2LocationSerializerWithPads } from './ll2-location';

export interface LaunchBase {
  id: string;
  external_id: number;
}

export interface LaunchBaseSerializer {
  id: string;
  data: LL2LocationSerializerWithPads;
}

export interface LaunchBaseList {
  count: number;
  launches: LaunchBaseSerializer[];
}
