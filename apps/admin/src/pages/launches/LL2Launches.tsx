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
import { RefreshCw } from "lucide-react";
import { launchService, taskService } from "@/services";
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

interface Ll2LaunchRow {
  id: string;
  name: string;
  date?: string;
  provider: string;
  vehicle: string;
  location: string;
  status: string;
}

export default function LL2Launches() {
  const [launches, setLaunches] = useState<Ll2LaunchRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState(false);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const perPage = 20;

  const fetchLaunches = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const { launches: items, count } = await launchService.getLL2Launches(perPage, offset);
      setTotalCount(count);

      if (count === 0) {
        setLaunches([]);
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

      const rows: Ll2LaunchRow[] = items.map((item) => ({
        id: `ll2-${item.id}`,
        name: item.name,
        date: item.net,
        provider: item.launch_service_provider?.name || "N/A",
        vehicle: item.rocket?.configuration?.name || "N/A",
        location: item.pad?.location?.name || item.pad?.name || "N/A",
        status: item.status?.abbrev || item.status?.name || "Unknown",
      }));

      setLaunches(rows);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to fetch launches"
      );
      setLaunches([]);
    } finally {
      setLoading(false);
    }
  }, [perPage]);

  useEffect(() => {
    fetchLaunches(page);
  }, [page, fetchLaunches]);

  const handleSync = async () => {
    try {
      setSyncing(true);
      await taskService.startTask("launch");
      toast.success("Launch task started");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to start launch task"
      );
    } finally {
      setSyncing(false);
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
                    href={`#/launches/ll2?page=${value}`}
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
        <h1 className="text-3xl font-bold">Launches (LL2)</h1>
        <Button
          onClick={handleSync}
          disabled={syncing}
          variant="outline"
          className="w-[150px]"
        >
          <RefreshCw
            className={`h-4 w-4 mr-2 ${syncing ? "animate-spin" : ""}`}
          />
          {syncing ? "Starting..." : "Start Launch Task"}
        </Button>
      </div>
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
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {launches.map((launch) => (
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
