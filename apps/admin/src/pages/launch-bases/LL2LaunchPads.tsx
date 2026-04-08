import type { MouseEvent } from "react";
import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
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
import { cn, buildPaginationRange } from "@/lib/utils";
import type { LL2Pad } from "@/types/ll2-location";
import { RefreshCw } from "lucide-react";
import { launchBaseService, taskService } from "@/services";
import { toast } from "sonner";

interface LaunchBaseRow {
  id: string;
  name: string;
  location?: string;
  country?: string;
  latitude?: number;
  longitude?: number;
}

const mapLL2Pad = (pad: LL2Pad): LaunchBaseRow => ({
  id: `ll2-pad-${pad.id}`,
  name: pad.name,
  location: pad.location?.name || "N/A",
  country: pad.location?.country?.name || "N/A",
  latitude: typeof pad.latitude === "number" ? pad.latitude : undefined,
  longitude: typeof pad.longitude === "number" ? pad.longitude : undefined,
});

export default function LL2LaunchPads() {
  const [rows, setRows] = useState<LaunchBaseRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const perPage = 20;

  const fetchBases = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const { pads, count } = await launchBaseService.getLL2Pads(perPage, offset);
      setTotalCount(count);

      if (count === 0) {
        setRows([]);
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

      setRows(pads.map(mapLL2Pad));
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch launch pads");
      setRows([]);
    } finally {
      setLoading(false);
    }
  }, [perPage]);

  useEffect(() => {
    fetchBases(page);
  }, [page, fetchBases]);

  const handleSync = async () => {
    try {
      setSyncing(true);
      await taskService.startTask("pad");
      toast.success("Pad task started");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to start pad task");
    } finally {
      setSyncing(false);
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
                    href={`#/launch-bases/ll2-pads?page=${value}`}
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
        <h1 className="text-3xl font-bold">Launch Bases (LL2 Pads)</h1>
        <Button
          onClick={handleSync}
          disabled={syncing}
          variant="outline"
          className="w-[150px]"
        >
          <RefreshCw
            className={`h-4 w-4 mr-2 ${syncing ? "animate-spin" : ""}`}
          />
          {syncing ? "Starting..." : "Start Pad Task"}
        </Button>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-8 text-center">Loading launch pads...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">
              No LL2 pads found.
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Location</TableHead>
                    <TableHead>Country</TableHead>
                    <TableHead>Latitude</TableHead>
                    <TableHead>Longitude</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.location || "N/A"}</TableCell>
                      <TableCell>{row.country || "N/A"}</TableCell>
                      <TableCell>
                        {typeof row.latitude === "number" ? row.latitude.toFixed(4) : "N/A"}
                      </TableCell>
                      <TableCell>
                        {typeof row.longitude === "number" ? row.longitude.toFixed(4) : "N/A"}
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
    </div>
  );
}
