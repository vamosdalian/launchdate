export interface PageBackground {
  id?: string;
  page_key: string;
  display_name: string;
  background_image: string;
  configured: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface UpdatePageBackgroundPayload {
  background_image: string;
}