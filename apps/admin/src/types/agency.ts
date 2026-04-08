import type { LL2AgencyNormal } from './ll2-agency';

export interface Agency {
  id: number;
  external_id: number;
  thumb_image?: string;
  images: string[];
  social_url: SocialUrl[];
}

export interface SocialUrl {
  name: string;
  url: string;
}

export interface AgencySerializer extends Agency {
  data: LL2AgencyNormal;
}

export interface AgencyList {
  count: number;
  agencies: AgencySerializer[];
}
