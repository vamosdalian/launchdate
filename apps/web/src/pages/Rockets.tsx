import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Search, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import ListPagination from '../components/ListPagination';
import PageHero from '../components/PageHero';
import { fetchRockets } from '../services/rocketsService';
import type { RocketListItem } from '../types';

const PAGE_SIZE = 20;

const Rockets = () => {
  const [rockets, setRockets] = useState<RocketListItem[]>([]);
  const [page, setPage] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [refreshToken, setRefreshToken] = useState(0);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadRockets = async () => {
      try {
        setLoading(true);
        const response = await fetchRockets({ page, search: appliedSearch });
        if (!isMounted) {
          return;
        }

        setRockets(response.rockets);
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

    loadRockets();

    return () => {
      isMounted = false;
    };
  }, [page, appliedSearch, refreshToken]);

  const handleSearchSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setPage(0);
    setAppliedSearch(searchTerm.trim());
    setRefreshToken((currentValue) => currentValue + 1);
  };

  const handleClearSearch = () => {
    setSearchTerm('');
    setPage(0);
    setAppliedSearch('');
    setRefreshToken((currentValue) => currentValue + 1);
  };

  if (loading && rockets.length === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading rockets...</p>
        </div>
      </div>
    );
  }

  if (error && rockets.length === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-lg mb-4">Error loading rockets: {error.message}</p>
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
        pageKey="rockets"
        title="Rockets"
        description="Explore the rockets that are shaping the future of space exploration"
      />

      {/* Filters and Search */}
      <section className="sticky top-20 z-40 bg-[#0F2854]/85 py-8 backdrop-blur-md">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <form className="flex flex-col gap-4 md:flex-row" onSubmit={handleSearchSubmit}>
              <div className="relative flex-grow">
                <Search
                  className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
                  aria-hidden="true"
                />
                <input
                  type="text"
                  placeholder="Search rockets, e.g. 'Falcon 9'..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full bg-[#2a2a2a] border border-[#4a4a4a] rounded-lg py-3 px-4 pl-10 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                />
              </div>
              <div className="flex gap-3">
                <Button type="submit" className="h-[50px] bg-blue-600 px-5 text-white hover:bg-blue-700">
                  Search
                </Button>
                {appliedSearch && (
                  <Button type="button" variant="outline" className="h-[50px] border-[#4a4a4a] bg-[#111111] px-5 text-white hover:bg-[#1b1b1b]" onClick={handleClearSearch}>
                    <X className="h-4 w-4" />
                    Clear
                  </Button>
                )}
              </div>
            </form>
            <p className="mt-3 text-sm text-gray-400">
              {appliedSearch ? `Showing ${totalCount} rockets for "${appliedSearch}".` : `Showing ${totalCount} rockets.`}
            </p>
          </div>
        </div>
      </section>

      <section className="py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {rockets.map((rocket) => (
                <Link
                  key={rocket.id}
                  to={`/rockets/${rocket.id}`}
                  className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg overflow-hidden hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 flex flex-col group"
                >
                  <div className="aspect-[2/3] bg-[#0a0a0a] relative overflow-hidden">
                    <img 
                      src={rocket.thumb_image} 
                      alt={rocket.name} 
                      className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" 
                    />
                  </div>
                  <div className="p-4 text-center">
                    <h3 className="text-lg font-bold group-hover:text-blue-400 transition-colors">{rocket.name}</h3>
                  </div>
                </Link>
              ))}
            </div>
            {!loading && rockets.length === 0 && (
              <div className="text-center py-12">
                <p className="text-gray-400 text-lg">No rockets found matching your criteria.</p>
              </div>
            )}
            <ListPagination currentPage={page} totalCount={totalCount} pageSize={PAGE_SIZE} onPageChange={setPage} />
            {loading && rockets.length > 0 && (
              <p className="mt-6 text-sm text-gray-500">Updating results…</p>
            )}
          </div>
        </div>
      </section>
    </div>
  );
};

export default Rockets;
