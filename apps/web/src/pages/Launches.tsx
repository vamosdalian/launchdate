import { Link } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Timeline, TimelineItem, type TimelineColor, type TimelineStatus } from '@/components/ui/timeline';
import { motion } from 'framer-motion';
import { CalendarClock, CalendarSync, CalendarX2, Check, EyeOff, MapPin, Rocket, Search, Sunrise, X, X as StatusX } from 'lucide-react';
import { useEffect, useState } from 'react';
import { fetchRocketLaunches } from '../services/launchesService';
import PageHero from '../components/PageHero';
import type { PublicLaunchSummary } from '../types';

const PAGE_SIZE = 20;

const launchTimelineDateFormat: Intl.DateTimeFormatOptions = {
  year: 'numeric',
  month: 'long',
  day: 'numeric',
};

const launchCardTimeFormat: Intl.DateTimeFormatOptions = {
  year: 'numeric',
  month: 'long',
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
  timeZone: 'UTC',
  timeZoneName: 'short',
};

const Launches = () => {
  const [launches, setLaunches] = useState<PublicLaunchSummary[]>([]);
  const [page, setPage] = useState(0);
  const [searchInput, setSearchInput] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [refreshToken, setRefreshToken] = useState(0);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadLaunches = async () => {
      try {
        setLoading(true);
        const response = await fetchRocketLaunches({ page, search: appliedSearch });
        if (!isMounted) {
          return;
        }

        setError(null);
        setTotalCount(response.count);
        setLaunches((previousLaunches) => {
          if (page === 0) {
            return response.launches;
          }
          return [...previousLaunches, ...response.launches];
        });
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

    loadLaunches();

    return () => {
      isMounted = false;
    };
  }, [page, appliedSearch, refreshToken]);

  useEffect(() => {
    setHasMore(launches.length < totalCount);
  }, [launches.length, totalCount]);

  const handleSearchSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLaunches([]);
    setTotalCount(0);
    setHasMore(true);
    setPage(0);
    setAppliedSearch(searchInput.trim());
    setRefreshToken((currentValue) => currentValue + 1);
  };

  const handleClearSearch = () => {
    setSearchInput('');
    setLaunches([]);
    setTotalCount(0);
    setHasMore(true);
    setPage(0);
    setAppliedSearch('');
    setRefreshToken((currentValue) => currentValue + 1);
  };

  const handleLoadMore = () => {
    if (loading || !hasMore) {
      return;
    }
    setPage((currentPage) => currentPage + 1);
  };

  const formatLaunchDateTime = (launchTime: string) => {
    return new Intl.DateTimeFormat('en-US', launchCardTimeFormat).format(new Date(launchTime));
  };

  const getTimelineStatus = (status: PublicLaunchSummary['status']): TimelineStatus => {
    switch (status) {
      case 'success':
        return 'completed';
      case 'in_flight':
        return 'in-progress';
      case 'failure':
      case 'cancelled':
        return 'error';
      case 'scheduled':
      case 'delayed':
      case 'unknown':
      default:
        return 'pending';
    }
  };

  const getTimelineColor = (): TimelineColor => {
    return 'primary';
  };

  const getStatusBadge = (status: PublicLaunchSummary['status'], statusLabel: string) => {
    switch (status) {
      case 'scheduled':
        return (
          <Badge className="gap-1.5">
            <CalendarClock className="h-3.5 w-3.5" />
            Scheduled
          </Badge>
        );
      case 'delayed':
        return (
          <Badge variant="secondary" className="gap-1.5 bg-amber-500/15 text-amber-200 hover:bg-amber-500/20">
            <CalendarSync className="h-3.5 w-3.5" />
            Delayed
          </Badge>
        );
      case 'cancelled':
        return (
          <Badge variant="outline" className="gap-1.5 border-zinc-500/50 bg-zinc-500/10 text-zinc-200">
            <CalendarX2 className="h-3.5 w-3.5" />
            Cancelled
          </Badge>
        );
      case 'in_flight':
        return (
          <Badge variant="secondary" className="gap-1.5 bg-sky-500/15 text-sky-200 hover:bg-sky-500/20">
            <Sunrise className="h-3.5 w-3.5" />
            In Flight
          </Badge>
        );
      case 'unknown':
        return (
          <Badge variant="outline" className="gap-1.5 border-zinc-500/50 bg-zinc-500/10 text-zinc-200">
            <EyeOff className="h-3.5 w-3.5" />
            Unknown
          </Badge>
        );
      case 'success':
        return (
          <Badge className="gap-1.5 bg-green-600">
            <Check className="h-3.5 w-3.5" />
            Successful
          </Badge>
        );
      case 'failure':
        return (
          <Badge className="gap-1.5 bg-red-600 text-white hover:bg-red-500">
            <StatusX className="h-3.5 w-3.5" />
            Failed
          </Badge>
        );
      default:
        return (
          <Badge variant="outline" className="gap-1.5 border-zinc-500/50 bg-zinc-500/10 text-zinc-200">
            <EyeOff className="h-3.5 w-3.5" />
            {statusLabel}
          </Badge>
        );
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
      <PageHero
        pageKey="launches"
        title="Launch Schedule"
        description="Stay updated with upcoming and past rocket launches"
      />

      <section className="sticky top-20 z-40 border-b border-white/10 bg-[#0F2854]/85 py-6 backdrop-blur-md">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-5xl">
            <form className="flex flex-col gap-3 md:flex-row" onSubmit={handleSearchSubmit}>
              <div className="relative flex-1">
                <Search className="pointer-events-none absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  value={searchInput}
                  onChange={(event) => setSearchInput(event.target.value)}
                  placeholder="Search launches, rockets, companies, or launch sites"
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
              {appliedSearch ? `Showing ${totalCount} results for "${appliedSearch}".` : `Showing ${totalCount} launches.`}
            </p>
          </div>
        </div>
      </section>

      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="mx-auto max-w-4xl">
            {launches.length > 0 ? (
              <Timeline size="lg">
                {launches.map((launch, index) => {
                  const timelineColor = getTimelineColor();
                  const timelineStatus = getTimelineStatus(launch.status);

                  return (
                    <TimelineItem
                      key={launch.id}
                      date={launch.launch_time}
                      dateFormat={launchTimelineDateFormat}
                      icon={<Rocket className="h-5 w-5" />}
                      iconColor={timelineColor}
                      connectorColor="muted"
                      status={timelineStatus}
                      showConnector
                      className="last:pb-0"
                      animationDelay={Math.min(index, 7) * 0.06}
                    >
                      <motion.div
                        initial={{ opacity: 0, y: 24 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{
                          duration: 0.45,
                          delay: Math.min(index, 7) * 0.06,
                          ease: 'easeOut',
                        }}
                      >
                        <Link to={`/launches/${launch.id}`} className="group block">
                          <Card className="overflow-hidden border-[var(--color-card-border)] bg-[var(--color-card-bg)] text-[var(--color-card-fg)] shadow-[0_18px_60px_rgba(28,77,141,0.22)] backdrop-blur-sm transition duration-300 hover:-translate-y-1 hover:border-[var(--color-card-hover-border)] hover:bg-[var(--color-card-hover-bg)]">
                            <CardContent className="p-5 sm:p-6">
                              <div className="flex items-start gap-4 sm:gap-5">
                                {launch.thumb_image ? (
                                  <div className="h-24 w-32 shrink-0 overflow-hidden rounded-xl border border-white/10 bg-[#0d1118] sm:h-28 sm:w-40">
                                    <img
                                      src={launch.thumb_image}
                                      alt={launch.name}
                                      className="h-full w-full object-cover transition duration-500 group-hover:scale-[1.04]"
                                      onError={(e) => {
                                        e.currentTarget.parentElement?.remove();
                                      }}
                                    />
                                  </div>
                                ) : null}

                                <div className="min-w-0 flex-1">
                                  <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                    <div className="min-w-0 space-y-1.5">
                                      <h3 className="truncate text-lg font-semibold leading-tight text-foreground transition-colors group-hover:text-[var(--color-highlight-strong)] sm:text-xl">
                                        {launch.name}
                                      </h3>
                                      <p className="text-sm text-muted-foreground">
                                        {launch.company.name || 'Unknown launch provider'}
                                      </p>
                                      <div className="text-sm text-muted-foreground">
                                        {formatLaunchDateTime(launch.launch_time)}
                                      </div>
                                      <div className="flex min-w-0 items-center gap-2 text-sm text-muted-foreground">
                                        <MapPin className="mt-0.5 h-4 w-4 shrink-0 text-[var(--color-highlight-strong)]" />
                                        <span className="block truncate">{launch.launch_base.name}</span>
                                      </div>
                                    </div>
                                    <div className="shrink-0">{getStatusBadge(launch.status, launch.status_label)}</div>
                                  </div>
                                </div>
                              </div>
                            </CardContent>
                          </Card>
                        </Link>
                      </motion.div>
                    </TimelineItem>
                  );
                })}
              </Timeline>
            ) : null}
            {!loading && launches.length === 0 && (
              <div className="rounded-lg border border-dashed border-[#2a2a2a] bg-[#111111] p-8 text-center text-gray-400">
                No launches matched your search.
              </div>
            )}
            {loading && page > 0 && (
              <div className="py-4 text-center">
                <div className="mx-auto h-8 w-8 animate-spin rounded-full border-b-2 border-blue-500"></div>
              </div>
            )}
            {!loading && hasMore && launches.length >= PAGE_SIZE && (
              <div className="pt-6 text-center">
                <Button
                  type="button"
                  variant="outline"
                  className="border-[#3a3a3a] bg-[#141414] text-white hover:bg-[#1d1d1d]"
                  onClick={handleLoadMore}
                >
                  Load More
                </Button>
              </div>
            )}
          </div>
        </div>
      </section>
    </div>
  );
};

export default Launches;
