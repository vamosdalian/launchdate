import type { LL2LaunchNormal } from './ll2-launch';

export interface Launch {
  id: string;
  external_id: string;
  background_image: string;
  image_list: string[];
  thumb_image: string;
}

export interface LaunchSerializer extends Launch {
  data: LL2LaunchNormal;
}

export interface LaunchList {
  count: number;
  launches: LaunchSerializer[];
}
