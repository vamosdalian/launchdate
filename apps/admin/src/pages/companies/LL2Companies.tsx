import { useState, useEffect, useCallback } from 'react';
import type { MouseEvent } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination.tsx';
import { cn, buildPaginationRange } from '@/lib/utils';
import type { LL2AgencyDetailed } from '@/types/ll2-agency';
import { RefreshCw } from 'lucide-react';
import { companyService, taskService } from '@/services';
import { toast } from 'sonner';

interface CompanyRow {
  id: string;
  name: string;
  countryCode?: string;
}

const mapLL2Agency = (agency: LL2AgencyDetailed): CompanyRow => ({
  id: String(agency.id),
  name: agency.name,
  countryCode: agency.country?.map(c => c.alpha2_code).join(', ') || undefined,
});

export default function LL2Companies() {
  const [rows, setRows] = useState<CompanyRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingLL2, setLoadingLL2] = useState(false);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const perPage = 20;

  const fetchCompanies = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const { agencies, count } = await companyService.getLL2Agencies(perPage, offset);
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

      setRows(agencies.map(mapLL2Agency));
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to fetch companies');
      setRows([]);
    } finally {
      setLoading(false);
    }
  }, [perPage]);

  useEffect(() => {
    fetchCompanies(page);
  }, [page, fetchCompanies]);

  const handleLL2Load = async () => {
    try {
      setLoadingLL2(true);
      await taskService.startTask('agency');
      toast.success('Agency task started');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to start agency task');
    } finally {
      setLoadingLL2(false);
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
                className={cn(isFirst && 'pointer-events-none opacity-50')}
                aria-disabled={isFirst}
              />
            </PaginationItem>
            {range.map((value, index) => {
              if (value === 'ellipsis') {
                return (
                  <PaginationItem key={`ellipsis-${index}`}>
                    <PaginationEllipsis />
                  </PaginationItem>
                );
              }
              return (
                <PaginationItem key={value}>
                  <PaginationLink
                    href={`#/companies/ll2?page=${value}`}
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
                className={cn(isLast && 'pointer-events-none opacity-50')}
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
        <h1 className="text-3xl font-bold">Companies (LL2)</h1>
        <div className="flex h-9 items-center gap-3">
          <Button
            onClick={handleLL2Load}
            disabled={loadingLL2}
            variant="outline"
            className="w-[150px]"
          >
            <RefreshCw
              className={`h-4 w-4 mr-2 ${loadingLL2 ? "animate-spin" : ""}`}
            />
            {loadingLL2 ? "Starting..." : "Start Agency Task"}
          </Button>
        </div>
      </div>
      <Card>
        <CardContent className="p-0">
          {loading ? (
            <div className="p-8 text-center">Loading companies...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">
              No companies found.
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Country</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.countryCode || 'N/A'}</TableCell>
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
