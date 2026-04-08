import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { cn } from "@/lib/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import type { RocketSerializer } from "@/types/rocket";
import type { RocketFilters } from "@/services/rocketService";
import { RefreshCw, FileJson, Pencil, ArrowUp, ArrowDown, Trash2, Plus } from "lucide-react";
import { rocketService } from "@/services";
import { toast } from "sonner";
import { buildPaginationRange } from "@/lib/utils";
import type { MouseEvent, FormEvent } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ImageSelectionModal } from "@/components/ImageSelectionModal";

interface RocketRow {
  id: string;
  backendId?: string;
  name: string;
  variant?: string;
  family?: string;
  manufacturer?: string;
  description?: string;
  source: "prod";
}

const mapProdRocket = (rocket: RocketSerializer): RocketRow => {
  const data = rocket.data;
  const familyNames = data.families?.map((f) => f.name).join(", ");
  const manufacturerName = data.manufacturer?.name;

  return {
    id: String(rocket.id),
    backendId: rocket.id,
    name: data.name || data.full_name || "Unknown",
    variant: data.variant,
    family: familyNames,
    manufacturer: manufacturerName,
    description: (data as { description?: string }).description,
    source: "prod",
  };
};

type RocketFilterState = {
  fullName: string;
  name: string;
  variant: string;
  sortBy: "default" | "full_name";
  sortOrder: "asc" | "desc";
};

const defaultFilters: RocketFilterState = {
  fullName: "",
  name: "",
  variant: "",
  sortBy: "default",
  sortOrder: "asc",
};

export default function ProdRockets() {
  const [rows, setRows] = useState<RocketRow[]>([]);
  const [rawRows, setRawRows] = useState<unknown[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [viewingRocket, setViewingRocket] = useState<unknown | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [filterForm, setFilterForm] = useState<RocketFilterState>(defaultFilters);
  const [appliedFilters, setAppliedFilters] = useState<RocketFilterState>(defaultFilters);
  const [editingRocket, setEditingRocket] = useState<RocketSerializer | null>(null);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [imageModalOpen, setImageModalOpen] = useState(false);
  const [activeImageField, setActiveImageField] = useState<"launch_image" | "main_image" | "thumb_image" | "image_list" | null>(null);
  const perPage = 20;

  const updateFilterForm = <K extends keyof RocketFilterState>(key: K, value: RocketFilterState[K]) => {
    setFilterForm((prev) => ({ ...prev, [key]: value }));
  };

  const handleFilterSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setAppliedFilters({ ...filterForm });
    setPage(1);
  };

  const handleResetFilters = () => {
    setFilterForm(defaultFilters);
    setAppliedFilters(defaultFilters);
    setPage(1);
  };

  const handleEdit = (rocket: RocketSerializer) => {
    setEditingRocket({ ...rocket });
    setIsEditOpen(true);
  };

  const handleSave = async () => {
    if (!editingRocket) return;
    try {
      await rocketService.updateRocket(editingRocket.id, {
        launch_image: editingRocket.launch_image,
        main_image: editingRocket.main_image,
        thumb_image: editingRocket.thumb_image,
        image_list: editingRocket.image_list,
      });
      toast.success("Rocket updated successfully");
      setIsEditOpen(false);
      fetchRockets(page);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update rocket");
    }
  };

  const handleImageSelect = (imageUrl: string) => {
    if (editingRocket && activeImageField) {
      if (activeImageField === "image_list") {
        setEditingRocket({
          ...editingRocket,
          image_list: [...(editingRocket.image_list || []), imageUrl],
        });
      } else {
        setEditingRocket({ ...editingRocket, [activeImageField]: imageUrl });
      }
      setImageModalOpen(false);
    }
  };

  const fetchRockets = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const filters: RocketFilters = {
        fullName: appliedFilters.fullName || undefined,
        name: appliedFilters.name || undefined,
        variant: appliedFilters.variant || undefined,
        sortBy: appliedFilters.sortBy === "full_name" ? "full_name" : undefined,
        sortOrder: appliedFilters.sortOrder,
      };

      const { rockets, count } = await rocketService.getProdRockets(perPage, offset, filters);
      setTotalCount(count);

      if (count === 0) {
        setRows([]);
        setRawRows([]);
        if (pageNumber !== 1) {
          setPage(1);
        }
        return;
      }

      const totalPages = Math.max(1, Math.ceil(count / perPage));
      if (pageNumber > totalPages) {
        setPage(totalPages);
        return;
      }

      setRows(rockets.map(mapProdRocket));
      setRawRows(rockets);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch rockets");
      setRows([]);
      setRawRows([]);
    } finally {
      setLoading(false);
    }
  }, [perPage, appliedFilters]);

  useEffect(() => {
    fetchRockets(page);
  }, [page, fetchRockets]);

  const totalPages = totalCount > 0 ? Math.ceil(totalCount / perPage) : 1;

  const handlePageChange = (nextPage: number) => {
    if (nextPage < 1 || nextPage > totalPages || nextPage === page) {
      return;
    }
    setPage(nextPage);
  };

  const renderPagination = () => {
    if (loading || totalCount <= perPage) {
      return null;
    }

    const range = buildPaginationRange(page, totalPages);
    const isFirst = page === 1;
    const isLast = page === totalPages;

    return (
      <div className="border-t px-4 py-4">
        <Pagination className="justify-end">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                href="#"
                onClick={(event: MouseEvent<HTMLAnchorElement>) => {
                  event.preventDefault();
                  handlePageChange(page - 1);
                }}
                className={cn(isFirst && "pointer-events-none opacity-50")}
                aria-disabled={isFirst}
              />
            </PaginationItem>
            {range.map((value, index) => {
              if (value === "ellipsis") {
                return (
                  <PaginationItem key={`ellipsis-${index}`}>
                    <PaginationEllipsis />
                  </PaginationItem>
                );
              }

              return (
                <PaginationItem key={value}>
                  <PaginationLink
                    href={`#/rockets/prod?page=${value}`}
                    isActive={value === page}
                    onClick={(event: MouseEvent<HTMLAnchorElement>) => {
                      event.preventDefault();
                      handlePageChange(value);
                    }}
                  >
                    {value}
                  </PaginationLink>
                </PaginationItem>
              );
            })}
            <PaginationItem>
              <PaginationNext
                href="#"
                onClick={(event: MouseEvent<HTMLAnchorElement>) => {
                  event.preventDefault();
                  handlePageChange(page + 1);
                }}
                className={cn(isLast && "pointer-events-none opacity-50")}
                aria-disabled={isLast}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
    );
  };

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="text-3xl font-bold">Rockets (Prod)</h1>
        <TooltipProvider>
          <div className="flex flex-wrap items-center gap-3">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  onClick={() => fetchRockets(page)}
                  disabled={loading}
                  variant="outline"
                  size="icon"
                >
                  <RefreshCw className={`h-4 w-4 ${loading ? "animate-spin" : ""}`} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Refresh data</p>
              </TooltipContent>
            </Tooltip>
          </div>
        </TooltipProvider>
      </div>

      <Card>
        <CardContent className="p-4">
          <form className="space-y-4" onSubmit={handleFilterSubmit}>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              <Input
                placeholder="Full name"
                value={filterForm.fullName}
                onChange={(event) => updateFilterForm("fullName", event.target.value)}
              />
              <Input
                placeholder="Name"
                value={filterForm.name}
                onChange={(event) => updateFilterForm("name", event.target.value)}
              />
              <Input
                placeholder="Variant"
                value={filterForm.variant}
                onChange={(event) => updateFilterForm("variant", event.target.value)}
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <p className="mb-2 text-sm font-medium">Sort by</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortBy}
                  onChange={(event) =>
                    updateFilterForm("sortBy", event.target.value as RocketFilterState["sortBy"])
                  }
                >
                  <option value="default">Default</option>
                  <option value="full_name">Full name</option>
                </select>
              </div>
              <div>
                <p className="mb-2 text-sm font-medium">Sort order</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortOrder}
                  onChange={(event) =>
                    updateFilterForm("sortOrder", event.target.value as RocketFilterState["sortOrder"])
                  }
                >
                  <option value="asc">Ascending</option>
                  <option value="desc">Descending</option>
                </select>
              </div>
            </div>
            <div className="flex flex-wrap justify-end gap-2">
              <Button type="button" variant="outline" onClick={handleResetFilters}>
                Reset
              </Button>
              <Button type="submit">Apply Filters</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-8 text-center">Loading rockets...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">No rockets found.</div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                      <TableHead>Name</TableHead>
                      <TableHead>Variant</TableHead>
                      <TableHead>Family</TableHead>
                      <TableHead>Manufacturer</TableHead>
                      <TableHead>Description</TableHead>
                    <TableHead className="text-center">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row, index) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.variant || "N/A"}</TableCell>
                      <TableCell>{row.family || "N/A"}</TableCell>
                      <TableCell>{row.manufacturer || "N/A"}</TableCell>
                      <TableCell className="max-w-xl truncate" title={row.description || "N/A"}>
                        {row.description || "N/A"}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          {row.backendId ? (
                            <>
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      onClick={() => {
                                        const rawRocket = rawRows[index] as RocketSerializer;
                                        handleEdit(rawRocket);
                                      }}
                                    >
                                      <Pencil className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p>Edit rocket</p>
                                  </TooltipContent>
                                </Tooltip>
                              </TooltipProvider>
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      onClick={() => {
                                        const rawRocket = rawRows[index];
                                        setViewingRocket(rawRocket || null);
                                        setIsSheetOpen(true);
                                      }}
                                    >
                                    <FileJson className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p>View raw data</p>
                                  </TooltipContent>
                                </Tooltip>
                              </TooltipProvider>
                            </>
                          ) : (
                            <span className="text-sm text-muted-foreground">—</span>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              {renderPagination()}
            </>
          )}
        </CardContent>
      </Card>

      <Sheet open={isSheetOpen} onOpenChange={setIsSheetOpen}>
        <SheetContent className="sm:max-w-2xl overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Raw Rocket Data</SheetTitle>
          </SheetHeader>
          <div className="mt-4 rounded-lg bg-muted p-4">
            <pre className="overflow-x-auto text-sm">
              {viewingRocket ? JSON.stringify(viewingRocket, null, 2) : "No data available."}
            </pre>
          </div>
        </SheetContent>
      </Sheet>

      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className="sm:max-w-[800px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Rocket</DialogTitle>
          </DialogHeader>
          {editingRocket && (
            <div className="grid gap-6 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label className="text-right">ID</Label>
                <Input value={editingRocket.id || ''} disabled className="col-span-3" />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label className="text-right">Name</Label>
                <Input value={editingRocket.data.full_name || editingRocket.data.name || ''} disabled className="col-span-3" />
              </div>

              <div className="grid grid-cols-4 items-start gap-4">
                <Label className="text-right pt-2">Launch Image</Label>
                <div className="col-span-3 space-y-2">
                  <div className="relative group w-full max-w-md aspect-video bg-muted rounded-lg overflow-hidden border">
                    {editingRocket.launch_image ? (
                      <img
                        src={editingRocket.launch_image}
                        alt="Launch"
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="flex items-center justify-center w-full h-full text-muted-foreground">
                        No image selected
                      </div>
                    )}
                    <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <Button
                        variant="secondary"
                        onClick={() => {
                          setActiveImageField("launch_image");
                          setImageModalOpen(true);
                        }}
                      >
                        Change
                      </Button>
                    </div>
                  </div>
                  <Input
                    value={editingRocket.launch_image || ""}
                    onChange={(e) =>
                      setEditingRocket({ ...editingRocket, launch_image: e.target.value })
                    }
                    placeholder="Launch Image URL"
                  />
                </div>
              </div>

              <div className="grid grid-cols-4 items-start gap-4">
                <Label className="text-right pt-2">Main Image</Label>
                <div className="col-span-3 space-y-2">
                  <div className="relative group w-full max-w-md aspect-video bg-muted rounded-lg overflow-hidden border">
                    {editingRocket.main_image ? (
                      <img
                        src={editingRocket.main_image}
                        alt="Main"
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="flex items-center justify-center w-full h-full text-muted-foreground">
                        No image selected
                      </div>
                    )}
                    <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <Button
                        variant="secondary"
                        onClick={() => {
                          setActiveImageField("main_image");
                          setImageModalOpen(true);
                        }}
                      >
                        Change
                      </Button>
                    </div>
                  </div>
                  <Input
                    value={editingRocket.main_image || ""}
                    onChange={(e) =>
                      setEditingRocket({ ...editingRocket, main_image: e.target.value })
                    }
                    placeholder="Main Image URL"
                  />
                </div>
              </div>

              <div className="grid grid-cols-4 items-start gap-4">
                <Label className="text-right pt-2">Thumbnail Image</Label>
                <div className="col-span-3 space-y-2">
                  <div className="relative group w-full max-w-md aspect-video bg-muted rounded-lg overflow-hidden border">
                    {editingRocket.thumb_image ? (
                      <img
                        src={editingRocket.thumb_image}
                        alt="Thumbnail"
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="flex items-center justify-center w-full h-full text-muted-foreground">
                        No image selected
                      </div>
                    )}
                    <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <Button
                        variant="secondary"
                        onClick={() => {
                          setActiveImageField("thumb_image");
                          setImageModalOpen(true);
                        }}
                      >
                        Change
                      </Button>
                    </div>
                  </div>
                  <Input
                    value={editingRocket.thumb_image || ""}
                    onChange={(e) =>
                      setEditingRocket({ ...editingRocket, thumb_image: e.target.value })
                    }
                    placeholder="Thumbnail Image URL"
                  />
                </div>
              </div>

              <div className="grid grid-cols-4 items-start gap-4">
                <Label className="text-right pt-2">Image List</Label>
                <div className="col-span-3 space-y-4">
                  <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                    {(editingRocket.image_list || []).map((url, index) => (
                      <div key={index} className="relative group aspect-square bg-muted rounded-lg overflow-hidden border">
                        <img src={url} alt={`Image ${index + 1}`} className="w-full h-full object-cover" />
                        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-2">
                          <div className="flex gap-2">
                            <Button
                              variant="secondary"
                              size="icon"
                              className="h-8 w-8"
                              disabled={index === 0}
                              onClick={() => {
                                const newList = [...(editingRocket.image_list || [])];
                                [newList[index - 1], newList[index]] = [newList[index], newList[index - 1]];
                                setEditingRocket({ ...editingRocket, image_list: newList });
                              }}
                            >
                              <ArrowUp className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="secondary"
                              size="icon"
                              className="h-8 w-8"
                              disabled={index === (editingRocket.image_list || []).length - 1}
                              onClick={() => {
                                const newList = [...(editingRocket.image_list || [])];
                                [newList[index + 1], newList[index]] = [newList[index], newList[index + 1]];
                                setEditingRocket({ ...editingRocket, image_list: newList });
                              }}
                            >
                              <ArrowDown className="h-4 w-4" />
                            </Button>
                          </div>
                          <Button
                            variant="destructive"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() => {
                              const newList = (editingRocket.image_list || []).filter((_, i) => i !== index);
                              setEditingRocket({ ...editingRocket, image_list: newList });
                            }}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    ))}
                    <Button
                      variant="outline"
                      className="aspect-square flex flex-col items-center justify-center gap-2 h-full"
                      onClick={() => {
                        setActiveImageField("image_list");
                        setImageModalOpen(true);
                      }}
                    >
                      <Plus className="h-8 w-8" />
                      <span>Add Image</span>
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          )}
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleSave}>Save changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ImageSelectionModal
        open={imageModalOpen}
        onOpenChange={setImageModalOpen}
        onSelect={handleImageSelect}
      />
    </div>
  );
}
