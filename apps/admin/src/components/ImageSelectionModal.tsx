import { useState, useEffect, useCallback, useRef } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Loader2, Upload } from "lucide-react";
import { imageService } from "@/services";
import { toast } from "sonner";
import type { Image } from "@/types/image";

interface ImageSelectionModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSelect: (imageUrl: string) => void;
}

export function ImageSelectionModal({ open, onOpenChange, onSelect }: ImageSelectionModalProps) {
  const [images, setImages] = useState<Image[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const limit = 12;
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchImages = useCallback(async (pageNumber: number) => {
    setLoading(true);
    try {
      const offset = (pageNumber - 1) * limit;
      const data = await imageService.getImages(limit, offset);
      setImages(data.images || []);
      setTotalCount(data.count || 0);
    } catch (error) {
      console.error(error);
      toast.error('Failed to load images');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (open) {
      fetchImages(page);
    }
  }, [open, page, fetchImages]);

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      setLoading(true);
      await imageService.uploadImage(file);
      toast.success("Image uploaded successfully");
      setPage(1);
      fetchImages(1);
    } catch (error) {
      console.error(error);
      toast.error("Failed to upload image");
    } finally {
      setLoading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const totalPages = Math.ceil(totalCount / limit);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[800px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Select Image</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="flex justify-end">
            <input
              type="file"
              ref={fileInputRef}
              className="hidden"
              accept="image/*"
              onChange={handleUpload}
            />
            <Button size="sm" onClick={() => fileInputRef.current?.click()}>
              <Upload className="mr-2 h-4 w-4" />
              Upload
            </Button>
          </div>
          {loading ? (
            <div className="flex justify-center p-8">
              <Loader2 className="h-8 w-8 animate-spin" />
            </div>
          ) : (
            <div className="grid grid-cols-3 sm:grid-cols-4 gap-4">
              {images.map((image) => (
                <div
                  key={image.id}
                  className="aspect-square relative group cursor-pointer border rounded-md overflow-hidden hover:border-primary"
                  onClick={() => {
                    onSelect(image.url);
                    onOpenChange(false);
                  }}
                >
                  <img
                    src={image.url}
                    alt={image.name}
                    className="object-cover w-full h-full"
                    loading="lazy"
                  />
                  <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors" />
                </div>
              ))}
            </div>
          )}
          {totalPages > 1 && (
            <div className="flex justify-center gap-2 mt-4">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page <= 1}
              >
                Previous
              </Button>
              <span className="flex items-center text-sm">
                Page {page} of {totalPages}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages}
              >
                Next
              </Button>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
