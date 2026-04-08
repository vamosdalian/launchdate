import type { MouseEvent, FormEvent } from "react";
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
} from "@/components/ui/pagination.tsx";
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
import { cn, buildPaginationRange } from "@/lib/utils";
import type { LaunchBaseSerializer } from "@/types/launch-base";
import type { LaunchBaseFilters } from "@/services/launchBaseService";
import { RefreshCw, FileJson } from "lucide-react";
import { launchBaseService } from "@/services";
import { toast } from "sonner";

interface LaunchBaseRow {
  id: string;
  backendId?: number;
  name: string;
  location?: string;
  country?: string;
  coordinates?: string;
}

const mapProdBase = (base: LaunchBaseSerializer): LaunchBaseRow => {
  const data = base.data;
  const pads = data.pads || [];
  const firstPadCoord = pads.find(
    (pad) => typeof pad?.latitude === "number" && typeof pad?.longitude === "number"
  );
  const latitude = typeof firstPadCoord?.latitude === "number" ? firstPadCoord.latitude : data.latitude;
  const longitude = typeof firstPadCoord?.longitude === "number" ? firstPadCoord.longitude : data.longitude;
  const coordinates =
    typeof latitude === "number" && typeof longitude === "number"
      ? `${latitude.toFixed(4)}, ${longitude.toFixed(4)}`
      : undefined;

  return {
    id: String(base.id),
    backendId: base.id,
    name: data.name,
    location: pads.length > 0 ? pads[0]?.name ?? undefined : undefined,
    country: data.country?.name || undefined,
    coordinates,
  };
};

type LaunchBaseFilterState = {
  name: string;
  celestialBody: string;
  country: string;
  sortBy: "default" | "name";
  sortOrder: "asc" | "desc";
};

const defaultFilters: LaunchBaseFilterState = {
  name: "",
  celestialBody: "",
  country: "",
  sortBy: "default",
  sortOrder: "asc",
};

export default function ProdLaunchBases() {
  const [rows, setRows] = useState<LaunchBaseRow[]>([]);
  const [rawRows, setRawRows] = useState<unknown[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [viewingBase, setViewingBase] = useState<unknown | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [filterForm, setFilterForm] = useState<LaunchBaseFilterState>(defaultFilters);
  const [appliedFilters, setAppliedFilters] = useState<LaunchBaseFilterState>(defaultFilters);
  const perPage = 20;

  const updateFilterForm = <K extends keyof LaunchBaseFilterState>(key: K, value: LaunchBaseFilterState[K]) => {
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

  const fetchBases = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const filters: LaunchBaseFilters = {
        name: appliedFilters.name || undefined,
        celestialBody: appliedFilters.celestialBody || undefined,
        country: appliedFilters.country || undefined,
        sortBy: appliedFilters.sortBy === "name" ? "name" : undefined,
        sortOrder: appliedFilters.sortOrder,
      };

      const { launches, count } = await launchBaseService.getProdLaunchBases(perPage, offset, filters);
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

      setRows(launches.map(mapProdBase));
      setRawRows(launches);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch launch bases");
      setRows([]);
      setRawRows([]);
    } finally {
      setLoading(false);
    }
  }, [perPage, appliedFilters]);

  useEffect(() => {
    fetchBases(page);
  }, [page, fetchBases]);

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
                    href={`#/launch-bases/prod?page=${value}`}
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
        <h1 className="text-3xl font-bold">Launch Bases (Prod)</h1>
        <TooltipProvider>
          <div className="flex flex-wrap items-center gap-3">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  onClick={() => fetchBases(page)}
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
                placeholder="Celestial body"
                value={filterForm.celestialBody}
                onChange={(event) => updateFilterForm("celestialBody", event.target.value)}
              />
              <Input
                placeholder="Country"
                value={filterForm.country}
                onChange={(event) => updateFilterForm("country", event.target.value)}
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <p className="mb-2 text-sm font-medium">Sort by</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortBy}
                  onChange={(event) =>
                    updateFilterForm("sortBy", event.target.value as LaunchBaseFilterState["sortBy"])
                  }
                >
                  <option value="default">Default order</option>
                  <option value="name">Name</option>
                </select>
              </div>
              <div>
                <p className="mb-2 text-sm font-medium">Sort order</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortOrder}
                  onChange={(event) =>
                    updateFilterForm("sortOrder", event.target.value as LaunchBaseFilterState["sortOrder"])
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
            <div className="p-8 text-center">Loading launch bases...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">
              No launch bases found.
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Location</TableHead>
                    <TableHead>Country</TableHead>
                    <TableHead>Coordinates</TableHead>
                    <TableHead className="text-center">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row, index) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.location || "N/A"}</TableCell>
                      <TableCell>{row.country || "N/A"}</TableCell>
                      <TableCell>{row.coordinates || "N/A"}</TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end">
                          {typeof row.backendId === "number" ? (
                            <TooltipProvider>
                              <Tooltip>
                                <TooltipTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    onClick={() => {
                                      const rawBase = rawRows[index];
                                      setViewingBase(rawBase || null);
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
            <SheetTitle>Raw Launch Base Data</SheetTitle>
          </SheetHeader>
          <div className="mt-4 rounded-lg bg-muted p-4">
            <pre className="overflow-x-auto text-sm">
              {viewingBase ? JSON.stringify(viewingBase, null, 2) : "No data available."}
            </pre>
          </div>
        </SheetContent>
      </Sheet>
    </div>
  );
}
