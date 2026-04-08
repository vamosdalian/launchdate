import { useState, useEffect, useCallback } from "react";
import { Button } from "@/components/ui/button";
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
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import type { LL2LauncherConfigNormal } from "@/types/ll2-launcher";
import { RefreshCw } from "lucide-react";
import { rocketService, taskService } from "@/services";
import { toast } from "sonner";
import { buildPaginationRange } from "@/lib/utils";
import type { MouseEvent } from "react";

interface RocketRow {
  id: string;
  name: string;
  variant?: string;
  family?: string;
  manufacturer?: string;
  source: "ll2-launchers";
}

const mapLL2Launcher = (launcher: LL2LauncherConfigNormal): RocketRow => ({
  id: `ll2-launcher-${launcher.id}`,
  name: launcher.full_name || launcher.name,
  variant: launcher.variant || "N/A",
  family: launcher.families?.map((f) => f.name).join(", ") || "N/A",
  manufacturer: launcher.manufacturer?.name || "N/A",
  source: "ll2-launchers",
});

export default function LL2Launchers() {
  const [rows, setRows] = useState<RocketRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const perPage = 20;

  const fetchRockets = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const { launchers, count } = await rocketService.getLL2Launchers(perPage, offset);
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

      setRows(launchers.map(mapLL2Launcher));
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to fetch rockets");
      setRows([]);
    } finally {
      setLoading(false);
    }
  }, [perPage]);

  useEffect(() => {
    fetchRockets(page);
  }, [page, fetchRockets]);

  const handleSync = async () => {
    try {
      setSyncing(true);
      await taskService.startTask("launcher");
      toast.success("Launcher task started");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to start launcher task");
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
                    href={`#/rockets/ll2-launchers?page=${value}`}
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
        <h1 className="text-3xl font-bold">Rockets (LL2 Launchers)</h1>
        <div className="flex flex-wrap items-center gap-3">
          <Button
            onClick={handleSync}
            disabled={syncing}
            variant="outline"
            className="w-[150px]"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${syncing ? "animate-spin" : ""}`} />
            {syncing ? "Starting..." : "Start Launcher Task"}
          </Button>
        </div>
      </div>

      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-8 text-center">Loading rockets...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">No LL2 launchers found.</div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Variant</TableHead>
                    <TableHead>Family</TableHead>
                    <TableHead>Manufacturer</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.variant || "N/A"}</TableCell>
                      <TableCell>{row.family || "N/A"}</TableCell>
                      <TableCell>{row.manufacturer || "N/A"}</TableCell>
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
