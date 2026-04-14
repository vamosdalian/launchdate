import React, { useEffect, useState, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';
import { Loader2, Trash2, Image as ImageIcon, Upload, Crop } from 'lucide-react';

import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

import { imageService } from '@/services';
import type { Image } from '@/types/image';
import { cn } from '@/lib/utils';

interface ImageCardProps {
  image: Image;
  onDelete: (key: string) => Promise<void>;
  onGenerateThumbnail: (id: string, width: number, height: number) => Promise<void>;
}

const ImageCard = ({ image, onDelete, onGenerateThumbnail }: ImageCardProps) => {
  const [width, setWidth] = useState(100);
  const [height, setHeight] = useState(100);
  const [isThumbOpen, setIsThumbOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  return (
    <Card className="overflow-hidden">
      <div className="aspect-square relative group bg-muted">
        <img
          src={image.url}
          alt={image.name}
          className="object-cover w-full h-full transition-transform group-hover:scale-105"
          loading="lazy"
        />
        <div className={cn(
          "absolute inset-0 bg-black/40 transition-opacity flex items-center justify-center gap-2",
          "opacity-0 group-hover:opacity-100 has-[[data-state=open]]:opacity-100"
        )}>
          <Button variant="secondary" size="icon" onClick={() => window.open(image.url, '_blank')} title="View Image">
            <ImageIcon className="h-4 w-4" />
          </Button>
          
          <Popover open={isThumbOpen} onOpenChange={setIsThumbOpen}>
            <PopoverTrigger asChild>
              <Button variant="secondary" size="icon" title="Generate Thumbnail">
                <Crop className="h-4 w-4" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80">
              <div className="grid gap-4">
                <div className="space-y-2">
                  <h4 className="font-medium leading-none">Thumbnail Size</h4>
                  <p className="text-sm text-muted-foreground">
                    Enter dimensions for the thumbnail.
                  </p>
                </div>
                <div className="grid gap-2">
                  <div className="grid grid-cols-3 items-center gap-4">
                    <Label htmlFor="width">Width</Label>
                    <Input
                      id="width"
                      type="number"
                      value={width}
                      onChange={(e) => setWidth(Number(e.target.value))}
                      className="col-span-2 h-8"
                    />
                  </div>
                  <div className="grid grid-cols-3 items-center gap-4">
                    <Label htmlFor="height">Height</Label>
                    <Input
                      id="height"
                      type="number"
                      value={height}
                      onChange={(e) => setHeight(Number(e.target.value))}
                      className="col-span-2 h-8"
                    />
                  </div>
                </div>
                <Button onClick={() => {
                  onGenerateThumbnail(image.id, width, height);
                  setIsThumbOpen(false);
                }}>Generate</Button>
              </div>
            </PopoverContent>
          </Popover>

          <Popover open={isDeleteOpen} onOpenChange={setIsDeleteOpen}>
            <PopoverTrigger asChild>
              <Button variant="destructive" size="icon" title="Delete Image">
                <Trash2 className="h-4 w-4" />
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-60">
                <div className="grid gap-4">
                  <div className="space-y-2">
                    <h4 className="font-medium leading-none">Confirm Delete</h4>
                    <p className="text-sm text-muted-foreground">
                      Are you sure you want to delete this image?
                    </p>
                  </div>
                  <div className="flex justify-end gap-2">
                    <Button variant="outline" size="sm" onClick={() => setIsDeleteOpen(false)}>Cancel</Button>
                    <Button variant="destructive" size="sm" onClick={() => {
                      onDelete(image.key);
                      setIsDeleteOpen(false);
                    }}>Delete</Button>
                  </div>
                </div>
            </PopoverContent>
          </Popover>
        </div>
      </div>
      <CardContent className="p-3">
        <div className="text-sm font-medium truncate" title={image.name}>{image.name}</div>
        <div className="text-xs text-muted-foreground flex justify-between mt-1">
          <span>{image.width}x{image.height}</span>
          <span>{(image.size / 1024).toFixed(1)} KB</span>
        </div>
      </CardContent>
    </Card>
  );
};

export default function Images() {
  const [images, setImages] = useState<Image[]>([]);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);
  const [totalCount, setTotalCount] = useState(0);
  const [searchParams, setSearchParams] = useSearchParams();
  
  const page = parseInt(searchParams.get('page') || '1');
  const limit = 20;
  const offset = (page - 1) * limit;

  const fetchImages = useCallback(async () => {
    setLoading(true);
    try {
      const data = await imageService.getImages(limit, offset);
      setImages(data.images || []);
      setTotalCount(data.count || 0);
    } catch (error) {
      console.error(error);
      toast.error('Failed to load images');
    } finally {
      setLoading(false);
    }
  }, [limit, offset]);

  useEffect(() => {
    fetchImages();
  }, [fetchImages]);

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    try {
      await imageService.uploadImage(file);
      toast.success('Image uploaded successfully');
      fetchImages();
    } catch (error) {
      console.error(error);
      toast.error('Failed to upload image');
    } finally {
      setUploading(false);
      // Reset input
      e.target.value = '';
    }
  };

  const handleDelete = async (key: string) => {
    try {
      await imageService.deleteImage(key);
      toast.success('Image deleted');
      fetchImages();
    } catch (error) {
      console.error(error);
      toast.error('Failed to delete image');
    }
  };

  const handleGenerateThumbnail = async (id: string, width: number, height: number) => {
    try {
      await imageService.generateThumbnail({ id, width, height });
      toast.success('Thumbnail generation started');
    } catch (error) {
      console.error(error);
      toast.error('Failed to generate thumbnail');
    }
  };

  const totalPages = Math.ceil(totalCount / limit);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Images</h1>
        <div className="flex items-center gap-2">
          <label className={cn(buttonVariants({ variant: 'default' }), "cursor-pointer", uploading && "opacity-50 pointer-events-none")}>
            {uploading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Upload className="mr-2 h-4 w-4" />}
            Upload Image
            <input type="file" className="hidden" accept="image/*" onChange={handleUpload} disabled={uploading} />
          </label>
        </div>
      </div>

      {loading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
            {images?.map((image) => (
              <ImageCard 
                key={image.id} 
                image={image} 
                onDelete={handleDelete} 
                onGenerateThumbnail={handleGenerateThumbnail} 
              />
            ))}
          </div>

          {totalPages > 1 && (
            <Pagination className="mt-6">
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious 
                    href="#" 
                    onClick={(e) => {
                      e.preventDefault();
                      if (page > 1) setSearchParams({ page: String(page - 1) });
                    }}
                    className={page <= 1 ? 'pointer-events-none opacity-50' : ''}
                  />
                </PaginationItem>
                <PaginationItem>
                  <span className="px-4 text-sm">Page {page} of {totalPages}</span>
                </PaginationItem>
                <PaginationItem>
                  <PaginationNext 
                    href="#" 
                    onClick={(e) => {
                      e.preventDefault();
                      if (page < totalPages) setSearchParams({ page: String(page + 1) });
                    }}
                    className={page >= totalPages ? 'pointer-events-none opacity-50' : ''}
                  />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          )}
        </>
      )}
    </div>
  );
}
