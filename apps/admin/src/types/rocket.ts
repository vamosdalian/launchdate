import type { LL2LauncherConfigNormal } from './ll2-launcher';

export interface Rocket {
  id: string;
  external_id: number;
  launch_image: string;
  main_image: string;
  thumb_image: string;
  image_list: string[];
}

export interface RocketSerializer extends Rocket {
  data: LL2LauncherConfigNormal;
}

export interface RocketList {
  count: number;
  rockets: RocketSerializer[];
}
