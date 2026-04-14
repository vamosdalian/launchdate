import { Link } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Search, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import ListPagination from '../components/ListPagination';
import PageHero from '../components/PageHero';
import { fetchLaunchBases } from '../services/launchBasesService';
import type { LaunchBase } from '../types';

const PAGE_SIZE = 20;

const LaunchBases = () => {
  const [launchBases, setLaunchBases] = useState<LaunchBase[]>([]);
  const [page, setPage] = useState(0);
  const [searchInput, setSearchInput] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [refreshToken, setRefreshToken] = useState(0);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadLaunchBases = async () => {
      try {
        setLoading(true);
        const response = await fetchLaunchBases({ page, search: appliedSearch });
        if (!isMounted) {
          return;
        }

        setLaunchBases(response.launch_bases);
        setTotalCount(response.count);
        setError(null);
      } catch (err) {
        if (isMounted) {
          setError(err as Error);
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    loadLaunchBases();

    return () => {
      isMounted = false;
    };
  }, [page, appliedSearch, refreshToken]);

  const handleSearchSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setPage(0);
    setAppliedSearch(searchInput.trim());
    setRefreshToken((currentValue) => currentValue + 1);
  };

  const handleClearSearch = () => {
    setSearchInput('');
    setPage(0);
    setAppliedSearch('');
    setRefreshToken((currentValue) => currentValue + 1);
  };

  if (loading && launchBases.length === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading launch bases...</p>
        </div>
      </div>
    );
  }

  if (error && launchBases.length === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-lg mb-4">Error loading launch bases: {error.message}</p>
          <button 
            onClick={() => window.location.reload()} 
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <PageHero
        pageKey="launch-bases"
        title="Launch Sites"
        description="Explore the facilities where rockets are launched into space"
      />

      <section className="border-b border-white/10 bg-[#0F2854]/80 py-8 backdrop-blur-md">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-6xl">
            <form className="flex flex-col gap-4 md:flex-row" onSubmit={handleSearchSubmit}>
              <div className="relative flex-1">
                <Search className="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  value={searchInput}
                  onChange={(event) => setSearchInput(event.target.value)}
                  placeholder="Search launch sites, countries, or celestial bodies"
                  className="w-full rounded-lg border border-[#3a3a3a] bg-[#141414] py-3 pl-10 pr-4 text-white outline-none transition focus:border-blue-500"
                />
              </div>
              <div className="flex gap-3">
                <Button type="submit" className="h-[50px] bg-blue-600 px-5 text-white hover:bg-blue-700">
                  Search
                </Button>
                {appliedSearch && (
                  <Button type="button" variant="outline" className="h-[50px] border-[#3a3a3a] bg-[#141414] px-5 text-white hover:bg-[#1d1d1d]" onClick={handleClearSearch}>
                    <X className="h-4 w-4" />
                    Clear
                  </Button>
                )}
              </div>
            </form>
            <p className="mt-3 text-sm text-gray-400">
              {appliedSearch ? `Showing ${totalCount} launch sites for "${appliedSearch}".` : `Showing ${totalCount} launch sites.`}
            </p>
          </div>
        </div>
      </section>

      {/* Launch Bases Grid */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {launchBases.map((base) => (
              <Link key={base.id} to={`/bases/${base.id}`}>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg overflow-hidden hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer h-full flex flex-col">
                  <div className="aspect-video bg-[#0a0a0a]">
                    <img 
                      src={base.image_url} 
                      alt={base.name} 
                      className="w-full h-full object-cover"
                    />
                  </div>
                  <div className="p-6 flex flex-col flex-grow">
                    <h3 className="text-2xl font-bold mb-2">{base.name}</h3>
                    <p className="text-sm text-gray-400 mb-4 flex items-center gap-2">
                      <span>📍</span>
                      <span>{base.location}, {base.country}</span>
                    </p>
                    <p className="text-gray-300 mb-4 flex-grow">{base.description}</p>
                    
                    <div className="border-t border-[#2a2a2a] pt-4 mt-auto">
                      <div className="flex gap-2 mb-4">
                        <Badge variant="secondary" className="bg-[#2a2a2a]">
                          Lat: {base.latitude.toFixed(4)}°
                        </Badge>
                        <Badge variant="secondary" className="bg-[#2a2a2a]">
                          Long: {base.longitude.toFixed(4)}°
                        </Badge>
                      </div>
                      <Badge variant="secondary" className="w-full justify-center py-2 bg-[#2a2a2a] hover:bg-[#3a3a3a]">
                        View Details →
                      </Badge>
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
          {!loading && launchBases.length === 0 && (
            <div className="rounded-lg border border-dashed border-[#2a2a2a] bg-[#111111] p-8 text-center text-gray-400">
              No launch sites matched your search.
            </div>
          )}
          <ListPagination currentPage={page} totalCount={totalCount} pageSize={PAGE_SIZE} onPageChange={setPage} />
          {loading && launchBases.length > 0 && (
            <p className="mt-6 text-sm text-gray-500">Updating results…</p>
          )}
        </div>
      </section>
    </div>
  );
};

export default LaunchBases;
