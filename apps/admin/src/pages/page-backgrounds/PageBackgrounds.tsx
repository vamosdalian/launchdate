import { useCallback, useEffect, useState } from 'react';
import { Loader2, Pencil, Upload, X } from 'lucide-react';
import { toast } from 'sonner';

import { ImageSelectionModal } from '@/components/ImageSelectionModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { pageBackgroundService } from '@/services';
import type { PageBackground } from '@/types/page-background';

export default function PageBackgrounds() {
  const [pageBackgrounds, setPageBackgrounds] = useState<PageBackground[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [selectedPage, setSelectedPage] = useState<PageBackground | null>(null);
  const [draftImage, setDraftImage] = useState('');
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isImageModalOpen, setIsImageModalOpen] = useState(false);

  const fetchPageBackgrounds = useCallback(async () => {
    try {
      setLoading(true);
      const data = await pageBackgroundService.getPageBackgrounds();
      setPageBackgrounds(data);
    } catch (error) {
      console.error(error);
      toast.error('Failed to load page backgrounds');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPageBackgrounds();
  }, [fetchPageBackgrounds]);

  const openEditor = (pageBackground: PageBackground) => {
    setSelectedPage(pageBackground);
    setDraftImage(pageBackground.background_image || '');
    setIsDialogOpen(true);
  };

  const handleSave = async () => {
    if (!selectedPage) {
      return;
    }

    try {
      setSaving(true);
      const updated = await pageBackgroundService.updatePageBackground(selectedPage.page_key, {
        background_image: draftImage.trim(),
      });

      setPageBackgrounds((current) =>
        current.map((item) => (item.page_key === updated.page_key ? updated : item)),
      );
      setIsDialogOpen(false);
      toast.success(`${selectedPage.display_name} background updated`);
    } catch (error) {
      console.error(error);
      toast.error('Failed to save page background');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Page Backgrounds</h1>
          <p className="text-muted-foreground">
            Configure the hero background image for the web home page and public list pages.
          </p>
        </div>
        <Button variant="outline" onClick={fetchPageBackgrounds} disabled={loading}>
          {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
          Refresh
        </Button>
      </div>

      {loading ? (
        <div className="flex justify-center py-16">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {pageBackgrounds.map((pageBackground) => (
            <Card key={pageBackground.page_key} className="overflow-hidden">
              <div className="aspect-[16/10] bg-muted relative">
                {pageBackground.background_image ? (
                  <img
                    src={pageBackground.background_image}
                    alt={pageBackground.display_name}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full items-center justify-center bg-muted text-sm text-muted-foreground">
                    No background configured
                  </div>
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent" />
                <div className="absolute bottom-0 left-0 right-0 p-4 text-white">
                  <div className="text-xs uppercase tracking-[0.2em] text-white/70">{pageBackground.page_key}</div>
                  <div className="text-xl font-semibold">{pageBackground.display_name}</div>
                </div>
              </div>
              <CardHeader>
                <CardTitle>{pageBackground.display_name}</CardTitle>
                <CardDescription>
                  {pageBackground.configured
                    ? `Last updated ${pageBackground.updated_at ? new Date(pageBackground.updated_at).toLocaleString() : 'recently'}`
                    : 'Uses the web fallback background until configured.'}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Button className="w-full" onClick={() => openEditor(pageBackground)}>
                  <Pencil className="mr-2 h-4 w-4" />
                  Edit Background
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="sm:max-w-[640px]">
          <DialogHeader>
            <DialogTitle>Edit Page Background</DialogTitle>
            <DialogDescription>
              {selectedPage ? `Update the hero background image for ${selectedPage.display_name}.` : 'Update the selected page background.'}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div className="space-y-2">
              <Label htmlFor="background-image">Background image URL</Label>
              <Input
                id="background-image"
                value={draftImage}
                onChange={(event) => setDraftImage(event.target.value)}
                placeholder="https://cdn.example.com/background.jpg"
              />
            </div>

            <div className="flex flex-wrap gap-2">
              <Button type="button" variant="outline" onClick={() => setIsImageModalOpen(true)}>
                <Upload className="mr-2 h-4 w-4" />
                Choose From Images
              </Button>
              <Button type="button" variant="ghost" onClick={() => setDraftImage('')}>
                <X className="mr-2 h-4 w-4" />
                Clear Selection
              </Button>
            </div>

            <div className="overflow-hidden rounded-lg border bg-muted">
              {draftImage ? (
                <img src={draftImage} alt="Background preview" className="aspect-[16/9] w-full object-cover" />
              ) : (
                <div className="flex aspect-[16/9] items-center justify-center text-sm text-muted-foreground">
                  No background selected
                </div>
              )}
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDialogOpen(false)} disabled={saving}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={saving}>
              {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
              Save
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ImageSelectionModal
        open={isImageModalOpen}
        onOpenChange={setIsImageModalOpen}
        onSelect={(imageUrl) => setDraftImage(imageUrl)}
      />
    </div>
  );
}