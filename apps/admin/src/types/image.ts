export interface ImageThumb {
  key: string;
  url: string;
  width: number;
  height: number;
  size: number;
  content_type: string;
  upload_time: string;
}

export interface Image {
  id: number;
  key: string;
  url: string;
  name: string;
  width: number;
  height: number;
  size: number;
  content_type: string;
  upload_time: string;
  thumb_images?: ImageThumb[];
}

export interface ImageListResponse {
  count: number;
  images: Image[];
}

export interface GenerateThumbnailParams {
  id: number;
  width: number;
  height: number;
}
