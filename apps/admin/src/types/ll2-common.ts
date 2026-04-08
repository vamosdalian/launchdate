export interface LL2Image {
  id: number;
  name: string;
  image_url: string;
  thumbnail_url: string;
  credit: string | null;
  license: LL2ImageLicense | null;
  single_use: boolean;
  variants: LL2ImageVariant[];
}

export interface LL2ImageLicense {
  id: number;
  name: string;
  priority: number;
  link: string | null;
}

export interface LL2ImageVariant {
  id: number;
  url: string;
  width: number;
  height: number;
}

export interface LL2Country {
  id: number;
  name: string;
  alpha2_code: string;
  alpha3_code: string;
  nationality_name: string;
  nationality_name_composed: string;
}

export interface LL2SocialMedia {
  id: number;
  name: string;
  url: string;
  logo: LL2Image | null;
}

export interface LL2SocialMediaLink {
  id: number;
  social_media: LL2SocialMedia;
  url: string;
}

export interface LL2InfoURLType {
  id: number;
  name: string;
}

export interface LL2Language {
  id: number;
  name: string;
  code: string;
}

export interface LL2InfoURL {
  priority: number;
  source: string;
  title: string;
  description: string;
  feature_image: string | null;
  url: string;
  type: LL2InfoURLType | null;
  language: LL2Language | null;
}

export interface LL2VidURLType {
  id: number;
  name: string;
}

export interface LL2VidURL {
  priority: number;
  source: string;
  publisher: string | null;
  title: string;
  description: string;
  feature_image: string | null;
  url: string;
  type: LL2VidURLType | null;
  language: LL2Language | null;
  start_time: string | null;
  end_time: string | null;
  live: boolean;
}

export interface LL2CelestialBodyType {
  id: number;
  name: string;
}

export interface LL2CelestialBodyMini {
  response_mode: string;
  id: number;
  name: string;
}

export interface LL2CelestialBodyDetailed extends LL2CelestialBodyMini {
  type: LL2CelestialBodyType | null;
  diameter: number | null;
  mass: number | null;
  gravity: number | null;
  length_of_day: string | null;
  atmosphere: boolean | null;
  image: LL2Image | null;
  description: string | null;
  wiki_url: string | null;
  total_attempted_landes: number;
  successful_launches: number;
  failed_launches: number;
  total_attempted_landings: number;
  successful_landings: number;
  failed_landings: number;
}
