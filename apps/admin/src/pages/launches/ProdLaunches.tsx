import type { MouseEvent, FormEvent } from "react";
import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Trash2, RefreshCw, FileJson, Pencil, Plus, ArrowUp, ArrowDown } from "lucide-react";
import { launchService } from "@/services";
import type { LaunchFilters } from "@/services/launchService";
import type { LaunchSerializer } from "@/types/launch";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination.tsx";
import { cn, buildPaginationRange } from "@/lib/utils";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ImageSelectionModal } from "@/components/ImageSelectionModal";

interface ProdLaunchRow {
  id: string;
  externalId: string;
  name: string;
  date?: string;
  provider: string;
  vehicle: string;
  location: string;
  status: string;
}

const mapProdLaunch = (launch: LaunchSerializer): ProdLaunchRow => ({
  id: launch.id,
  externalId: launch.external_id,
  name: launch.data.name || "Unknown",
  date: launch.data.net,
  provider: launch.data.launch_service_provider?.name || "N/A",
  vehicle: launch.data.rocket?.configuration?.name || "N/A",
  location: launch.data.pad?.location?.name || launch.data.pad?.name || "N/A",
  status: launch.data.status?.name || launch.data.status?.abbrev || "Unknown",
});

type LaunchFilterState = Required<Pick<LaunchFilters, "sortBy" | "sortOrder">> & {
  name: string;
  status: string;
  provider: string;
  rocket: string;
  mission: string;
};

const defaultFilters: LaunchFilterState = {
  name: "",
  status: "",
  provider: "",
  rocket: "",
  mission: "",
  sortBy: "time",
  sortOrder: "asc",
};

export default function ProdLaunches() {
  const [launches, setLaunches] = useState<ProdLaunchRow[]>([]);
  const [rawLaunches, setRawLaunches] = useState<LaunchSerializer[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [viewingLaunch, setViewingLaunch] = useState<LaunchSerializer | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [editingLaunch, setEditingLaunch] = useState<ProdLaunchRow | null>(null);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isImageSelectOpen, setIsImageSelectOpen] = useState(false);
  const [imageSelectMode, setImageSelectMode] = useState<'main' | 'list' | 'thumb'>('main');
  const [editForm, setEditForm] = useState<{
    background_image: string;
    image_list: string[];
    thumb_image: string;
  }>({ background_image: "", image_list: [], thumb_image: "" });
  const [filterForm, setFilterForm] = useState<LaunchFilterState>(defaultFilters);
  const [appliedFilters, setAppliedFilters] = useState<LaunchFilterState>(defaultFilters);
  const perPage = 20;

  const updateFilterForm = <K extends keyof LaunchFilterState>(key: K, value: LaunchFilterState[K]) => {
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

  const fetchLaunches = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const { launches: items, count } = await launchService.getProdLaunches(perPage, offset, appliedFilters);
      setTotalCount(count);

      if (count === 0) {
        setLaunches([]);
        setRawLaunches([]);
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

      setLaunches(items.map(mapProdLaunch));
      setRawLaunches(items);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to fetch launches"
      );
      setLaunches([]);
      setRawLaunches([]);
    } finally {
      setLoading(false);
    }
  }, [perPage, appliedFilters]);

  useEffect(() => {
    fetchLaunches(page);
  }, [page, fetchLaunches]);

  const handleEditClick = (launch: ProdLaunchRow, rawLaunch: LaunchSerializer) => {
    setEditingLaunch(launch);
    setEditForm({
      background_image: rawLaunch.background_image || "",
      image_list: rawLaunch.image_list || [],
      thumb_image: rawLaunch.thumb_image || "",
    });
    setIsEditOpen(true);
  };

  const handleUpdateLaunch = async () => {
    if (!editingLaunch || !editingLaunch.id) return;

    try {
      await launchService.updateProdLaunch(editingLaunch.id, {
        background_image: editForm.background_image,
        image_list: editForm.image_list,
        thumb_image: editForm.thumb_image,
      });

      toast.success("Launch updated successfully");
      setIsEditOpen(false);
      fetchLaunches(page);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to update launch"
      );
    }
  };

  const handleImageSelect = (imageUrl: string) => {
    if (imageSelectMode === 'main') {
      setEditForm(prev => ({ ...prev, background_image: imageUrl }));
    } else if (imageSelectMode === 'thumb') {
      setEditForm(prev => ({ ...prev, thumb_image: imageUrl }));
    } else {
      setEditForm(prev => ({ ...prev, image_list: [...prev.image_list, imageUrl] }));
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) {
      return "N/A";
    }
    return new Date(dateString).toLocaleString();
  };

  const getStatusColor = (status: string) => {
    const normalized = status.toLowerCase();
    switch (normalized) {
      case "scheduled":
      case "go":
        return "bg-blue-100 text-blue-800";
      case "successful":
      case "success":
        return "bg-green-100 text-green-800";
      case "failed":
      case "failure":
      case "no go":
        return "bg-red-100 text-red-800";
      case "cancelled":
      case "tbd":
        return "bg-gray-100 text-gray-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

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
                    href={`#/launches/prod?page=${value}`}
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
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Launches (Prod)</h1>
        <TooltipProvider>
          <div className="flex items-center gap-3">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  onClick={() => fetchLaunches(page)}
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
                placeholder="Name"
                value={filterForm.name}
                onChange={(event) => updateFilterForm("name", event.target.value)}
              />
              <Input
                placeholder="Status"
                value={filterForm.status}
                onChange={(event) => updateFilterForm("status", event.target.value)}
              />
              <Input
                placeholder="Provider"
                value={filterForm.provider}
                onChange={(event) => updateFilterForm("provider", event.target.value)}
              />
              <Input
                placeholder="Rocket"
                value={filterForm.rocket}
                onChange={(event) => updateFilterForm("rocket", event.target.value)}
              />
              <Input
                placeholder="Mission"
                value={filterForm.mission}
                onChange={(event) => updateFilterForm("mission", event.target.value)}
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <p className="mb-2 text-sm font-medium">Sort by</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortBy}
                  onChange={(event) =>
                    updateFilterForm("sortBy", event.target.value as LaunchFilterState["sortBy"])
                  }
                >
                  <option value="time">Launch time</option>
                  <option value="name">Name</option>
                </select>
              </div>
              <div>
                <p className="mb-2 text-sm font-medium">Sort order</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortOrder}
                  onChange={(event) =>
                    updateFilterForm(
                      "sortOrder",
                      event.target.value as LaunchFilterState["sortOrder"]
                    )
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
            <div className="p-8 text-center">Loading launches...</div>
          ) : launches.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">
              No launches found.
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Date</TableHead>
                    <TableHead>Provider</TableHead>
                    <TableHead>Vehicle</TableHead>
                    <TableHead>Location</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-center">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {launches.map((launch, index) => (
                    <TableRow key={launch.id}>
                      <TableCell className="font-medium">
                        {launch.name}
                      </TableCell>
                      <TableCell>{formatDate(launch.date)}</TableCell>
                      <TableCell>{launch.provider}</TableCell>
                      <TableCell>{launch.vehicle}</TableCell>
                      <TableCell>{launch.location}</TableCell>
                      <TableCell>
                        <span
                          className={`px-2 py-1 rounded-full text-xs ${getStatusColor(
                            launch.status
                          )}`}
                        >
                          {launch.status}
                        </span>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-1">
                          {launch.id && (
                            <>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => {
                                  const rawLaunch = rawLaunches[index];
                                  setViewingLaunch(rawLaunch || null);
                                  setIsSheetOpen(true);
                                }}
                                title="View raw data"
                              >
                                <FileJson className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleEditClick(launch, rawLaunches[index])}
                                title="Edit launch"
                              >
                                <Pencil className="h-4 w-4" />
                              </Button>

                            </>
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
            <SheetTitle>Raw Launch Data</SheetTitle>
          </SheetHeader>
          <div className="mt-4 rounded-lg bg-muted p-4">
            <pre className="overflow-x-auto text-sm">
              {viewingLaunch ? JSON.stringify(viewingLaunch, null, 2) : "No data available."}
            </pre>
          </div>
        </SheetContent>
      </Sheet>

      <ImageSelectionModal
        open={isImageSelectOpen}
        onOpenChange={setIsImageSelectOpen}
        onSelect={handleImageSelect}
      />

      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className="sm:max-w-[800px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Launch</DialogTitle>
            <DialogDescription>
              Make changes to the launch here. Click update when you're done.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-6 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">ID</Label>
              <Input value={editingLaunch?.id || ''} disabled className="col-span-3" />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">External ID</Label>
              <Input value={editingLaunch?.externalId || ''} disabled className="col-span-3" />
            </div>
            
            <div className="grid grid-cols-4 items-start gap-4">
              <Label className="text-right pt-2">Main Image</Label>
              <div className="col-span-3 space-y-2">
                <div className="relative group w-full max-w-md aspect-video bg-muted rounded-lg overflow-hidden border">
                  {editForm.background_image ? (
                    <img
                      src={editForm.background_image}
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
                        setImageSelectMode('main');
                        setIsImageSelectOpen(true);
                      }}
                    >
                      Change
                    </Button>
                  </div>
                </div>
                <Input
                  value={editForm.background_image}
                  onChange={(e) => setEditForm({ ...editForm, background_image: e.target.value })}
                  placeholder="Image URL"
                />
              </div>
            </div>

            <div className="grid grid-cols-4 items-start gap-4">
              <Label className="text-right pt-2">Thumb Image</Label>
              <div className="col-span-3 space-y-2">
                <div className="relative group w-full max-w-md aspect-video bg-muted rounded-lg overflow-hidden border">
                  {editForm.thumb_image ? (
                    <img
                      src={editForm.thumb_image}
                      alt="Thumb"
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
                        setImageSelectMode('thumb');
                        setIsImageSelectOpen(true);
                      }}
                    >
                      Change
                    </Button>
                  </div>
                </div>
                <Input
                  value={editForm.thumb_image}
                  onChange={(e) => setEditForm({ ...editForm, thumb_image: e.target.value })}
                  placeholder="Image URL"
                />
              </div>
            </div>

            <div className="grid grid-cols-4 items-start gap-4">
              <Label className="text-right pt-2">Image List</Label>
              <div className="col-span-3 space-y-4">
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                  {editForm.image_list.map((url, index) => (
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
                              const newList = [...editForm.image_list];
                              [newList[index - 1], newList[index]] = [newList[index], newList[index - 1]];
                              setEditForm({ ...editForm, image_list: newList });
                            }}
                          >
                            <ArrowUp className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="secondary"
                            size="icon"
                            className="h-8 w-8"
                            disabled={index === editForm.image_list.length - 1}
                            onClick={() => {
                              const newList = [...editForm.image_list];
                              [newList[index + 1], newList[index]] = [newList[index], newList[index + 1]];
                              setEditForm({ ...editForm, image_list: newList });
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
                            const newList = editForm.image_list.filter((_, i) => i !== index);
                            setEditForm({ ...editForm, image_list: newList });
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
                      setImageSelectMode('list');
                      setIsImageSelectOpen(true);
                    }}
                  >
                    <Plus className="h-8 w-8" />
                    <span>Add Image</span>
                  </Button>
                </div>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleUpdateLaunch}>Update</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
