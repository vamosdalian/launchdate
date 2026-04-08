import type { MouseEvent, FormEvent } from 'react';
import { useState, useEffect, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
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
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ImageSelectionModal } from "@/components/ImageSelectionModal";
import { cn, buildPaginationRange } from '@/lib/utils';
import type { AgencySerializer } from '@/types/agency';
import type { CompanyFilters } from '@/services/companyService';
import { RefreshCw, FileJson, Pencil, Trash2, Plus, ArrowUp, ArrowDown } from 'lucide-react';
import { companyService } from '@/services';
import { toast } from 'sonner';

interface CompanyRow {
  id: string;
  backendId?: number;
  name: string;
  founder?: string;
  founded?: number;
  headquarters?: string;
  employees?: number;
  website?: string;
}

const mapProdCompany = (company: AgencySerializer): CompanyRow => {
  const data = company.data;
  const website = company.social_url?.find(s => s.name === 'Website' || s.name === 'Homepage')?.url;
  
  return {
    id: String(company.id),
    backendId: company.id,
    name: data.name || 'Unknown',
    founder: data.administrator || undefined,
    founded: data.founding_year || undefined,
    headquarters: undefined, // Not available in LL2AgencyNormal directly
    employees: undefined, // Not available in LL2AgencyNormal
    website: website,
  };
};

type CompanyFilterState = {
  name: string;
  type: string;
  country: string;
  sortBy: 'name' | 'founding_year';
  sortOrder: 'asc' | 'desc';
};

const defaultFilters: CompanyFilterState = {
  name: '',
  type: '',
  country: '',
  sortBy: 'name',
  sortOrder: 'asc',
};

export default function ProdCompanies() {
  const [rows, setRows] = useState<CompanyRow[]>([]);
  const [rawRows, setRawRows] = useState<unknown[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [viewingCompany, setViewingCompany] = useState<unknown | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [editingCompany, setEditingCompany] = useState<CompanyRow | null>(null);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isImageSelectOpen, setIsImageSelectOpen] = useState(false);
  const [imageSelectMode, setImageSelectMode] = useState<'thumb' | 'list'>('thumb');
  const [editForm, setEditForm] = useState<{
    thumb_image: string;
    images: string[];
    social_url: { name: string; url: string }[];
  }>({ thumb_image: "", images: [], social_url: [] });
  const [filterForm, setFilterForm] = useState<CompanyFilterState>(defaultFilters);
  const [appliedFilters, setAppliedFilters] = useState<CompanyFilterState>(defaultFilters);
  const perPage = 20;

  const updateFilterForm = <K extends keyof CompanyFilterState>(key: K, value: CompanyFilterState[K]) => {
    setFilterForm((prev) => ({ ...prev, [key]: value }));
  };

  const handleEditClick = (company: CompanyRow, rawCompany: AgencySerializer) => {
    setEditingCompany(company);
    setEditForm({
      thumb_image: rawCompany.thumb_image || "",
      images: rawCompany.images || [],
      social_url: rawCompany.social_url || [],
    });
    setIsEditOpen(true);
  };

  const handleUpdateCompany = async () => {
    if (!editingCompany || !editingCompany.id) return;

    try {
      await companyService.updateAgency(editingCompany.id, {
        thumb_image: editForm.thumb_image,
        images: editForm.images,
        social_url: editForm.social_url,
      });

      toast.success("Company updated successfully");
      setIsEditOpen(false);
      fetchCompanies(page);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to update company"
      );
    }
  };

  const handleImageSelect = (imageUrl: string) => {
    if (imageSelectMode === 'thumb') {
      setEditForm(prev => ({ ...prev, thumb_image: imageUrl }));
    } else {
      setEditForm(prev => ({ ...prev, images: [...prev.images, imageUrl] }));
    }
    setIsImageSelectOpen(false);
  };

  const handleSocialUrlChange = (index: number, field: 'name' | 'url', value: string) => {
    const newSocialUrl = [...editForm.social_url];
    newSocialUrl[index] = { ...newSocialUrl[index], [field]: value };
    setEditForm(prev => ({ ...prev, social_url: newSocialUrl }));
  };

  const addSocialUrl = () => {
    setEditForm(prev => ({
      ...prev,
      social_url: [...prev.social_url, { name: '', url: '' }]
    }));
  };

  const removeSocialUrl = (index: number) => {
    setEditForm(prev => ({
      ...prev,
      social_url: prev.social_url.filter((_, i) => i !== index)
    }));
  };

  const removeImage = (index: number) => {
    setEditForm(prev => ({
      ...prev,
      images: prev.images.filter((_, i) => i !== index)
    }));
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

  const fetchCompanies = useCallback(async (pageNumber: number) => {
    try {
      setLoading(true);
      const offset = (pageNumber - 1) * perPage;
      const filters: CompanyFilters = {
        name: appliedFilters.name || undefined,
        type: appliedFilters.type || undefined,
        country: appliedFilters.country || undefined,
        sortBy: appliedFilters.sortBy === 'founding_year' ? 'founding_year' : undefined,
        sortOrder: appliedFilters.sortOrder,
      };

      const { agencies, count } = await companyService.getProdAgencies(perPage, offset, filters);
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

      setRows(agencies.map(mapProdCompany));
      setRawRows(agencies);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to fetch companies');
      setRows([]);
      setRawRows([]);
    } finally {
      setLoading(false);
    }
  }, [perPage, appliedFilters]);

  useEffect(() => {
    fetchCompanies(page);
  }, [page, fetchCompanies]);

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
                    href={`#/companies/prod?page=${value}`}
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
      <div className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="text-3xl font-bold">Companies (Prod)</h1>
        <TooltipProvider>
          <div className="flex flex-wrap items-center gap-3">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  onClick={() => fetchCompanies(page)}
                  disabled={loading}
                  variant="outline"
                  size="icon"
                >
                  <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
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
                onChange={(event) => updateFilterForm('name', event.target.value)}
              />
              <Input
                placeholder="Type"
                value={filterForm.type}
                onChange={(event) => updateFilterForm('type', event.target.value)}
              />
              <Input
                placeholder="Country"
                value={filterForm.country}
                onChange={(event) => updateFilterForm('country', event.target.value)}
              />
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <p className="mb-2 text-sm font-medium">Sort by</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortBy}
                  onChange={(event) =>
                    updateFilterForm('sortBy', event.target.value as CompanyFilterState['sortBy'])
                  }
                >
                  <option value="name">Name</option>
                  <option value="founding_year">Founding Year</option>
                </select>
              </div>
              <div>
                <p className="mb-2 text-sm font-medium">Sort order</p>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  value={filterForm.sortOrder}
                  onChange={(event) =>
                    updateFilterForm('sortOrder', event.target.value as CompanyFilterState['sortOrder'])
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
            <div className="p-8 text-center">Loading companies...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">No companies found.</div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Founder</TableHead>
                    <TableHead>Founded</TableHead>
                    <TableHead>Headquarters</TableHead>
                    <TableHead>Employees</TableHead>
                    <TableHead>Website</TableHead>
                    <TableHead className="text-center">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row, index) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell>{row.founder || 'N/A'}</TableCell>
                      <TableCell>{row.founded || 'N/A'}</TableCell>
                      <TableCell>{row.headquarters || 'N/A'}</TableCell>
                      <TableCell>{row.employees || 'N/A'}</TableCell>
                      <TableCell>
                        {row.website ? (
                          <a
                            href={row.website}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-500 hover:underline"
                          >
                            Link
                          </a>
                        ) : (
                          'N/A'
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          {typeof row.backendId === 'number' ? (
                            <>
                              <TooltipProvider>
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      onClick={() => {
                                        const rawCompany = rawRows[index] as AgencySerializer;
                                        handleEditClick(row, rawCompany);
                                      }}
                                    >
                                      <Pencil className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p>Edit company</p>
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
                                        const rawCompany = rawRows[index];
                                        setViewingCompany(rawCompany || null);
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
            <SheetTitle>Raw Company Data</SheetTitle>
          </SheetHeader>
          <div className="mt-4 rounded-lg bg-muted p-4">
            <pre className="overflow-x-auto text-sm">
              {viewingCompany ? JSON.stringify(viewingCompany, null, 2) : 'No data available.'}
            </pre>
          </div>
        </SheetContent>
      </Sheet>

      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className="sm:max-w-[800px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Company</DialogTitle>
          </DialogHeader>
          <div className="grid gap-6 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">ID</Label>
              <Input value={editingCompany?.id || ''} disabled className="col-span-3" />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label className="text-right">Name</Label>
              <Input value={editingCompany?.name || ''} disabled className="col-span-3" />
            </div>

            <div className="grid grid-cols-4 items-start gap-4">
              <Label className="text-right pt-2">Thumbnail Image</Label>
              <div className="col-span-3 space-y-2">
                <div className="relative group w-full max-w-md aspect-video bg-muted rounded-lg overflow-hidden border">
                  {editForm.thumb_image ? (
                    <img
                      src={editForm.thumb_image}
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
                  placeholder="Thumbnail Image URL"
                />
              </div>
            </div>

            <div className="grid grid-cols-4 items-start gap-4">
              <Label className="text-right pt-2">Images</Label>
              <div className="col-span-3 space-y-4">
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                  {editForm.images.map((url, index) => (
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
                              const newImages = [...editForm.images];
                              [newImages[index - 1], newImages[index]] = [newImages[index], newImages[index - 1]];
                              setEditForm({ ...editForm, images: newImages });
                            }}
                          >
                            <ArrowUp className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="secondary"
                            size="icon"
                            className="h-8 w-8"
                            disabled={index === editForm.images.length - 1}
                            onClick={() => {
                              const newImages = [...editForm.images];
                              [newImages[index + 1], newImages[index]] = [newImages[index], newImages[index + 1]];
                              setEditForm({ ...editForm, images: newImages });
                            }}
                          >
                            <ArrowDown className="h-4 w-4" />
                          </Button>
                        </div>
                        <Button
                          variant="destructive"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => removeImage(index)}
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

            <div className="grid grid-cols-4 items-start gap-4">
              <Label className="text-right pt-2">Social URLs</Label>
              <div className="col-span-3 space-y-2">
                {editForm.social_url.map((social, index) => (
                  <div key={index} className="flex items-center gap-2">
                    <Input
                      placeholder="Name"
                      value={social.name}
                      onChange={(e) => handleSocialUrlChange(index, 'name', e.target.value)}
                      className="w-1/3"
                    />
                    <Input
                      placeholder="URL"
                      value={social.url}
                      onChange={(e) => handleSocialUrlChange(index, 'url', e.target.value)}
                      className="flex-1"
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => removeSocialUrl(index)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
                <Button
                  type="button"
                  variant="outline"
                  onClick={addSocialUrl}
                  className="w-full"
                >
                  <Plus className="mr-2 h-4 w-4" /> Add Social URL
                </Button>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleUpdateCompany}>Save changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <ImageSelectionModal
        open={isImageSelectOpen}
        onOpenChange={setIsImageSelectOpen}
        onSelect={handleImageSelect}
      />
    </div>
  );
}
