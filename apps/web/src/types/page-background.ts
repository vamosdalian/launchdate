export type PageBackgroundKey = 'home' | 'launches' | 'rockets' | 'launch-bases' | 'companies';

export interface PageBackground {
  id?: string;
  page_key: PageBackgroundKey;
  display_name: string;
  background_image: string;
  configured: boolean;
  updated_at?: string;
}