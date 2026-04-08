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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn, buildPaginationRange } from "@/lib/utils";
import type { LL2LocationSerializerWithPads } from "@/types/ll2-location";
import { RefreshCw } from "lucide-react";
import { launchBaseService, taskService } from "@/services";
import { toast } from "sonner";

interface LaunchBaseRow {
  id: string;
  name: string;
  country?: string;
  padsCount?: number;
  padNames?: string;
}

const mapLL2Location = (location: LL2LocationSerializerWithPads): LaunchBaseRow => {
  const pads = Array.isArray(location.pads) ? location.pads : [];
  const padNames = pads
    .map((pad) => pad?.name)
    .filter((name): name is string => Boolean(name))
    .join(", ");

  return {
    id: `ll2-location-${location.id}`,
    name: location.name,
    country: location.country?.name || "N/A",
    padsCount: pads.length,
    padNames: padNames || "N/A",
  };
};

export default function LL2Locations() {
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
      const { locations, count } = await launchBaseService.getLL2Locations(perPage, offset);
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

      setRows(locations.map(mapLL2Location));
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch launch bases");
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
      await taskService.startTask("location");
      toast.success("Location task started");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to start location task");
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
                    href={`#/launch-bases/ll2?page=${value}`}
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
        <h1 className="text-3xl font-bold">Launch Bases (LL2)</h1>
        <TooltipProvider>
          <div className="flex h-9 items-center gap-3">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  onClick={handleSync}
                  disabled={syncing || loading}
                  variant="outline"
                  size="icon"
                >
                  <RefreshCw className={`h-4 w-4 ${syncing ? "animate-spin" : ""}`} />
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>Start location task</p>
              </TooltipContent>
            </Tooltip>
          </div>
        </TooltipProvider>
      </div>

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
                    <TableHead>Country</TableHead>
                    <TableHead>Pads Count</TableHead>
                    <TableHead>Pads</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.country}</TableCell>
                      <TableCell>{row.padsCount}</TableCell>
                      <TableCell className="max-w-md truncate" title={row.padNames}>
                        {row.padNames}
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
