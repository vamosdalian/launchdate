import type { LL2Country, LL2Image, LL2SocialMediaLink } from "./ll2-common";

export interface LL2AgencyType {
  id: number;
  name: string;
}

export interface LL2AgencyMini {
  response_mode: string;
  id: number;
  url: string;
  name: string;
  abbrev: string;
  type: LL2AgencyType | null;
}

export interface LL2AgencyNormal extends LL2AgencyMini {
  featured: boolean;
  country: LL2Country[];
  description: string | null;
  administrator: string | null;
  founding_year: number | null;
  launchers: string;
  spacecraft: string;
  parent: string | null;
  image: LL2Image | null;
  logo: LL2Image | null;
  social_logo: LL2Image | null;
}

export interface LL2AgencyDetailed extends LL2AgencyNormal {
  total_launch_count: number;
  consecutive_successful_launches: number;
  successful_launches: number;
  failed_launches: number;
  pending_launches: number;
  consecutive_successful_landings: number;
  successful_landings: number;
  failed_landings: number;
  attempted_landings: number;
  successful_landings_spacecraft: number;
  failed_landings_spacecraft: number;
  attempted_landings_spacecraft: number;
  successful_landings_payload: number;
  failed_landings_payload: number;
  attempted_landings_payload: number;
  info_url: string | null;
  wiki_url: string | null;
  social_media_links: LL2SocialMediaLink[];
}
