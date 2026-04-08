import { Link } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { useCallback, useEffect, useState } from 'react';
import { fetchRocketLaunches } from '../services/launchesService';
import type { PublicCompactLaunch } from '../types';

const Launches = () => {
  const [launches, setLaunches] = useState<PublicCompactLaunch[]>([]);
  const [page, setPage] = useState(0);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const loadLaunches = useCallback(async (pageNum: number) => {
    try {
      setLoading(true);
      const response = await fetchRocketLaunches(pageNum);
      
      if (pageNum === 0) {
        setLaunches(response.launches);
      } else {
        setLaunches(prev => [...prev, ...response.launches]);
      }
      
      // If we got fewer items than expected (e.g. 20), or count suggests we are done
      // For now, just check if we got any launches. 
      // Ideally we check against response.count or a fixed page size.
      if (response.launches.length === 0) {
        setHasMore(false);
      }
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadLaunches(page);
  }, [page, loadLaunches]);

  // Infinite scroll handler
  useEffect(() => {
    const handleScroll = () => {
      if (
        window.innerHeight + document.documentElement.scrollTop < document.documentElement.offsetHeight - 100 ||
        loading ||
        !hasMore
      ) {
        return;
      }
      setPage(prev => prev + 1);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [loading, hasMore]);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC',
      timeZoneName: 'short',
    });
  };

  const getStatusBadge = (status: number) => {
    switch (status) {
      case 1: // Scheduled
        return <Badge>🕒 Scheduled</Badge>;
      case 3: // Successful
        return <Badge className="bg-green-600">✅ Successful</Badge>;
      case 4: // Failed
        return <Badge variant="destructive">❌ Failed</Badge>;
      default:
        return <Badge variant="outline">Status: {status}</Badge>;
    }
  };

  if (loading && page === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading launches...</p>
        </div>
      </div>
    );
  }

  if (error && page === 0) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-lg mb-4">Error loading launches: {error.message}</p>
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
      {/* Page Hero */}
      <section className="py-16 md:py-24 text-center bg-[#111]">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">Launch Schedule</h1>
          <p className="max-w-3xl mx-auto text-lg md:text-xl text-gray-400">
            Stay updated with upcoming and past rocket launches
          </p>
        </div>
      </section>

      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto">
            <div className="relative ml-3">
              {/* Timeline line */}
              <div className="absolute left-0 top-4 bottom-0 border-l-2 border-border" />

              {launches.map((launch) => (
                <div key={launch.id} className="relative pl-8 pb-12 last:pb-0">
                  {/* Timeline dot */}
                  <div className="absolute h-3 w-3 -translate-x-1/2 left-px top-3 rounded-full border-2 border-primary bg-background ring-8 ring-background" />

                  {/* Content */}
                  <Link to={`/launches/${launch.id}`}>
                    <div className="space-y-3 bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-6 hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer">
                      <div className="flex flex-col sm:flex-row gap-6">
                        {/* Thumbnail Image */}
                        {launch.thumb_image && (
                          <div className="w-full sm:w-32 h-32 bg-[#0a0a0a] rounded-lg overflow-hidden flex-shrink-0">
                            <img 
                              src={launch.thumb_image} 
                              alt={launch.name} 
                              className="w-full h-full object-cover"
                              onError={(e) => {
                                e.currentTarget.style.display = 'none';
                              }}
                            />
                          </div>
                        )}
                        
                        <div className="flex-1 space-y-3">
                          <div className="flex items-start justify-between gap-4">
                            <div className="flex-1">
                              <h3 className="text-xl font-semibold mb-1">{launch.name}</h3>
                              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <span>🗓️</span>
                                <span>{formatDate(launch.launch_time)}</span>
                              </div>
                            </div>
                            {getStatusBadge(launch.status)}
                          </div>
                          
                          <div className="flex flex-wrap gap-2">
                            {launch.rocket_name && (
                              <Badge variant="secondary" className="rounded-full">
                                🚀 {launch.rocket_name}
                              </Badge>
                            )}
                            {launch.location && (
                              <Badge variant="secondary" className="rounded-full">
                                📍 {launch.location}
                              </Badge>
                            )}
                            {launch.agency_name && (
                              <Badge variant="secondary" className="rounded-full">
                                🏢 {launch.agency_name}
                              </Badge>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </Link>
                </div>
              ))}
              {loading && page > 0 && (
                 <div className="text-center py-4">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
                 </div>
              )}
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Launches;
