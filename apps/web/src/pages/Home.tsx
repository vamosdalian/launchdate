import { Link } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useEffect, useState, useCallback, useMemo } from 'react';
import { ArrowRight } from 'lucide-react';
import { usePageBackground } from '../contexts/PageBackgroundContext';
import { useApi } from '../hooks/useApi';
import { fetchRocketLaunches } from '../services/launchesService';
import { fetchRockets } from '../services/rocketsService';
import { fetchCompanies } from '../services/companiesService';
import './Home.css';

const Home = () => {
  const configuredHeroImage = usePageBackground('home');
  const fetchLaunchesCallback = useCallback(() => fetchRocketLaunches({ page: 0 }), []);
  const { data: launchesData } = useApi(fetchLaunchesCallback);
  const launches = launchesData?.launches;
  
  const fetchRocketsCallback = useCallback(() => fetchRockets({ page: 0 }), []);
  const { data: rocketsData } = useApi(fetchRocketsCallback);
  const rockets = rocketsData?.rockets;
  
  const fetchCompaniesCallback = useCallback(() => fetchCompanies({ page: 0, homepageOnly: true }), []);
  const { data: companiesData } = useApi(fetchCompaniesCallback);
  const companies = companiesData?.companies;

  // Get the next upcoming launch
  const nextLaunch = useMemo(() => {
    if (!launches) return null;
    const now = new Date();
    const upcomingLaunches = launches
      .filter((launch) => {
        const launchDate = new Date(launch.launch_time);
        return launch.status === 'scheduled' && launchDate >= now;
      })
      .sort((a, b) => {
        return new Date(a.launch_time).getTime() - new Date(b.launch_time).getTime();
      });
    
    return upcomingLaunches[0];
  }, [launches]);

  // Countdown state
  const [countdown, setCountdown] = useState({
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0,
  });
  const primaryFallbackImage = 'https://images.unsplash.com/photo-1614728263952-84ea256ec346?q=80&w=2574&auto=format&fit=crop';
  const secondaryFallbackImage = 'https://images.unsplash.com/photo-1516849841032-87cbac4d88f7?q=80&w=2574&auto=format&fit=crop';
  const launchHeroImage = nextLaunch?.background_image?.trim() || '';
  const defaultHeroImage = configuredHeroImage || primaryFallbackImage;
  const preferredHeroImage = launchHeroImage || defaultHeroImage;
  const [heroImageSrc, setHeroImageSrc] = useState(preferredHeroImage);

  useEffect(() => {
    setHeroImageSrc(preferredHeroImage);
  }, [preferredHeroImage]);

  // Update countdown every second
  useEffect(() => {
    if (!nextLaunch) return;

    const updateCountdown = () => {
      const launchDate = new Date(nextLaunch.launch_time).getTime();
      const now = new Date().getTime();
      const distance = launchDate - now;

      if (distance < 0) {
        setCountdown({ days: 0, hours: 0, minutes: 0, seconds: 0 });
        return;
      }

      setCountdown({
        days: Math.floor(distance / (1000 * 60 * 60 * 24)),
        hours: Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
        minutes: Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)),
        seconds: Math.floor((distance % (1000 * 60)) / 1000),
      });
    };

    updateCountdown();
    const interval = setInterval(updateCountdown, 1000);
    return () => clearInterval(interval);
  }, [nextLaunch]);

  // Get recent launch results (successful or failed)
  const recentResults = useMemo(() => {
    if (!launches) return [];
    return launches
      .filter((launch) => launch.status === 'success' || launch.status === 'failure')
      .sort((a, b) => {
        return new Date(b.launch_time).getTime() - new Date(a.launch_time).getTime();
      })
      .slice(0, 6);
  }, [launches]);

  // Get featured rockets from the compact public list.
  const popularRockets = useMemo(() => {
    if (!rockets) return [];
    return rockets.slice(0, 4);
  }, [rockets]);

  // Get popular companies
  const popularCompanies = useMemo(() => {
    if (!companies) return [];
    return companies.slice(0, 5);
  }, [companies]);
  
  // Get upcoming launches for the section
  const upcomingLaunches = useMemo(() => {
    if (!launches) return [];
    const now = new Date();
    return launches
      .filter((launch) => {
        const launchDate = new Date(launch.launch_time);
        return launch.status === 'scheduled' && launchDate >= now;
      })
      .sort((a, b) => {
        return new Date(a.launch_time).getTime() - new Date(b.launch_time).getTime();
      });
  }, [launches]);

  const launchCadence = useMemo(() => {
    const now = Date.now();
    const windows = [
      { label: 'Next 24H', horizonMs: 24 * 60 * 60 * 1000 },
      { label: 'Next 3D', horizonMs: 3 * 24 * 60 * 60 * 1000 },
      { label: 'Next 7D', horizonMs: 7 * 24 * 60 * 60 * 1000 },
    ];

    return windows.map((window) => ({
      label: window.label,
      value: upcomingLaunches.filter((launch) => {
        const launchTime = new Date(launch.launch_time).getTime();
        return launchTime >= now && launchTime <= now + window.horizonMs;
      }).length,
    }));
  }, [upcomingLaunches]);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC',
    });
  };

  const formatHeroDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
      weekday: 'short',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC',
    });
  };

  const countdownItems = [
    { label: 'Days', value: countdown.days },
    { label: 'Hours', value: countdown.hours },
    { label: 'Minutes', value: countdown.minutes },
    { label: 'Seconds', value: countdown.seconds },
  ];

  return (
    <div className="min-h-screen bg-background">
      <section className="hero-section relative flex min-h-screen items-end overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img 
            src={heroImageSrc}
            alt="Rocket Launch Background"
            className="w-full h-full object-cover"
            onError={() => {
              if (launchHeroImage && heroImageSrc === launchHeroImage) {
                setHeroImageSrc(defaultHeroImage);
                return;
              }

              if (
                configuredHeroImage &&
                configuredHeroImage !== primaryFallbackImage &&
                heroImageSrc === configuredHeroImage
              ) {
                setHeroImageSrc(primaryFallbackImage);
                return;
              }

              if (heroImageSrc !== secondaryFallbackImage) {
                setHeroImageSrc(secondaryFallbackImage);
              }
            }}
          />
          <div className="absolute inset-0 bg-black/25" />
          <div className="hero-aurora hero-aurora-cyan" />
          <div className="hero-aurora hero-aurora-amber" />
          <div className="hero-grid-overlay" />
          <div
            className="absolute inset-x-0 bottom-0 h-2/3"
            style={{
              background: 'linear-gradient(180deg, rgba(5, 11, 20, 0) 0%, rgba(5, 11, 20, 0.52) 45%, rgba(5, 11, 20, 0.96) 100%)',
            }}
          />
        </div>

        <div className="relative z-10 container mx-auto px-4 pb-20 pt-28 md:pb-24 md:pt-36">
          <div className="hero-layout grid items-end gap-10 lg:grid-cols-[minmax(0,1.15fr)_minmax(320px,420px)] lg:gap-12">
            <div className="max-w-3xl text-left">
              <div className="hero-kicker mb-5 inline-flex items-center gap-3 rounded-full border border-white/10 bg-white/8 px-4 py-2 text-[11px] font-semibold uppercase tracking-[0.28em] text-white/90 backdrop-blur-xl">
                <span className="h-2 w-2 rounded-full bg-[var(--color-highlight-warm)] shadow-[0_0_18px_rgba(255,179,107,0.7)]" />
                {nextLaunch ? 'Next Mission Window' : 'Live Launch Tracking'}
              </div>

              {nextLaunch ? (
                <>
                  <h1 className="max-w-4xl text-5xl font-black leading-[0.95] tracking-[-0.04em] text-white md:text-7xl lg:text-[5.5rem]">
                    {nextLaunch.name}
                  </h1>
                  <p className="mt-6 max-w-2xl text-lg leading-8 text-gray-300 md:text-xl">
                    Follow the next scheduled ascent with live timing, vehicle details, and launch-site context before the window opens.
                  </p>

                  <div className="hero-info-list mt-8 max-w-2xl space-y-3">
                    <div className="hero-info-row">
                      <span className="hero-info-value">{formatHeroDate(nextLaunch.launch_time)} UTC</span>
                    </div>
                    <div className="hero-info-row">
                      <span className="hero-info-value">{nextLaunch.company.name}</span>
                    </div>
                    <div className="hero-info-row">
                      <span className="hero-info-value">{nextLaunch.launch_base.location}, {nextLaunch.launch_base.country}</span>
                    </div>
                  </div>

                  <div className="mt-8 flex flex-wrap gap-4">
                    <Button asChild size="lg" className="hero-primary-button px-7 text-base font-semibold">
                      <Link className="inline-flex items-center" to={`/launches/${nextLaunch.id}`}>
                        View Launch Details
                      </Link>
                    </Button>
                    <Button asChild size="lg" variant="outline" className="hero-secondary-button px-7 text-base font-semibold">
                      <Link className="inline-flex items-center whitespace-nowrap" to="/launches">
                        Browse Schedule
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </Link>
                    </Button>
                  </div>
                </>
              ) : (
                <>
                  <h1 className="max-w-4xl text-5xl font-black leading-[0.95] tracking-[-0.04em] text-white md:text-7xl lg:text-[5.5rem]">
                    Launch windows, rockets, and operators in one mission-ready feed.
                  </h1>
                  <p className="mt-6 max-w-2xl text-lg leading-8 text-gray-300 md:text-xl">
                    The next scheduled mission has not landed in the feed yet, but the full launch archive and upcoming schedule are still available.
                  </p>
                  <div className="mt-8 flex flex-wrap gap-4">
                    <Button asChild size="lg" className="hero-primary-button px-7 text-base font-semibold">
                      <Link className="inline-flex items-center" to="/launches">Browse Schedule</Link>
                    </Button>
                    <Button asChild size="lg" variant="outline" className="hero-secondary-button px-7 text-base font-semibold">
                      <Link className="inline-flex items-center whitespace-nowrap" to="/rockets">
                        Explore Rockets
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </Link>
                    </Button>
                  </div>
                </>
              )}
            </div>

            <div className="hero-panel rounded-[1.75rem] p-5 md:p-6 lg:-translate-y-10 lg:self-start">
              <div className="flex items-center justify-between gap-4 border-b border-white/10 pb-5">
                <div>
                  <p className="text-xs font-semibold uppercase tracking-[0.22em] text-[var(--color-highlight-strong)]">Launch Countdown</p>
                  <h2 className="mt-2 text-2xl font-bold text-white">{nextLaunch ? 'T-minus to liftoff' : 'Mission queue updating'}</h2>
                </div>
                <Badge className="hero-status-badge">{nextLaunch ? 'Scheduled' : 'Standby'}</Badge>
              </div>

              {nextLaunch ? (
                <>
                  <div className="countdown-grid mt-6 grid grid-cols-2 gap-3">
                    {countdownItems.map((item) => (
                      <div key={item.label} className="countdown-box rounded-2xl px-4 py-4 md:px-5 md:py-5">
                        <div className="countdown-value text-4xl font-black tracking-[-0.05em] md:text-5xl">
                          {String(item.value).padStart(2, '0')}
                        </div>
                        <div className="mt-2 text-xs font-semibold uppercase tracking-[0.2em] text-gray-400">{item.label}</div>
                      </div>
                    ))}
                  </div>

                  <div className="hero-cadence mt-6">
                    <div className="hero-cadence-header">
                      <p className="hero-cadence-title">Launch cadence</p>
                      <span className="hero-cadence-caption">Current tracked feed</span>
                    </div>
                    <div className="hero-cadence-list mt-4 space-y-3">
                      {launchCadence.map((item) => (
                        <div key={item.label} className="hero-cadence-row">
                          <div className="hero-cadence-label">{item.label}</div>
                          <div className="hero-cadence-value">{item.value}</div>
                        </div>
                      ))}
                      <div className="hero-cadence-row">
                        <div className="hero-cadence-label">Tracked Upcoming</div>
                        <div className="hero-cadence-value">{upcomingLaunches.length}</div>
                      </div>
                    </div>
                  </div>
                </>
              ) : (
                <div className="mt-6 rounded-2xl border border-white/10 bg-white/5 px-5 py-6 text-left text-gray-300">
                  <p className="text-base leading-7">
                    We are waiting for the next confirmed launch window. Open the schedule to review upcoming missions and recently completed flights.
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* Upcoming Launches Section */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-12 text-center">Upcoming Launches</h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {upcomingLaunches.slice(0, 6).map((launch) => (
              <Link key={launch.id} to={`/launches/${launch.id}`} className="block h-full">
                <div className="flex h-full flex-col bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-6 hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer">
                  <div className="flex justify-between items-start mb-4">
                    <h3 className="launch-card-title text-xl font-bold">{launch.name}</h3>
                    <Badge className="bg-blue-600">Scheduled</Badge>
                  </div>
                  <div className="space-y-2 text-sm text-gray-400 mb-4">
                    <p className="launch-card-meta">{launch.company.name}</p>
                    <p className="launch-card-meta">{formatDate(launch.launch_time)} UTC</p>
                    <p className="launch-card-meta">{launch.launch_base.name}</p>
                  </div>
                </div>
              </Link>
            ))}
          </div>
          <div className="text-center mt-8">
            <Button asChild variant="outline" size="lg">
              <Link to="/launches">View All Launches</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Recent Launch Results */}
      <section className="bg-[#0F2854] py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-12 text-center">Recent Launch Results</h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {recentResults.map((launch) => (
              <Link key={launch.id} to={`/launches/${launch.id}`} className="block h-full">
                <div className="flex h-full flex-col bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-6 hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer">
                  <div className="flex justify-between items-start mb-4">
                    <div className="flex-1">
                      <h3 className="launch-card-title mb-2 text-xl font-bold">{launch.name}</h3>
                      <div className="space-y-2 text-sm text-gray-400 mb-2">
                        <p className="launch-card-meta">{launch.company.name}</p>
                        <p className="launch-card-meta">{formatDate(launch.launch_time)} UTC</p>
                        <p className="launch-card-meta">{launch.launch_base.name}</p>
                      </div>
                    </div>
                    <Badge 
                      className={launch.status === 'success' ? 'bg-green-600' : 'bg-red-600'}
                    >
                      {launch.status === 'success' ? 'Success' : 'Failed'}
                    </Badge>
                  </div>
                </div>
              </Link>
            ))}
          </div>
          <div className="text-center mt-8">
            <Button asChild variant="outline" size="lg">
              <Link to="/launches">View More Launches</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Popular Rockets */}
      <section className="bg-[#0F2854] py-20">
        <div className="container mx-auto px-4">
          <div className="flex justify-between items-center mb-12">
            <h2 className="text-3xl font-bold">Popular Rockets</h2>
            <Button asChild variant="outline">
              <Link to="/rockets">View All</Link>
            </Button>
          </div>
          <div className="grid md:grid-cols-4 gap-6">
            {popularRockets.map((rocket) => (
              <Link key={rocket.id} to={`/rockets/${rocket.id}`}>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg overflow-hidden hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer h-full">
                  <div className="aspect-[2/3] bg-[#0a0a0a]">
                    <img
                      src={rocket.thumb_image}
                      alt={rocket.name}
                      className="w-full h-full object-cover"
                    />
                  </div>
                  <div className="p-6">
                    <h3 className="text-lg font-bold mb-1">{rocket.name}</h3>
                    <div className="flex gap-2 text-xs">
                      <Badge variant="secondary" className="bg-[#2a2a2a]">Rocket</Badge>
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Space Agencies Section */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">Explore Major Space Agencies</h2>
          <div className="flex flex-wrap justify-center items-center gap-x-12 sm:gap-x-16 md:gap-x-20 gap-y-8">
            {popularCompanies.map((company) => (
              <Link
                key={company.id} 
                to={`/companies/${company.id}`}
                title={company.name}
                className="grayscale brightness-75 hover:grayscale-0 hover:brightness-100 transition-all duration-300"
              >
                <img 
                  src={company.image_url} 
                  alt={`${company.name} Logo`} 
                  className="h-12 md:h-16 max-w-[150px] object-contain"
                />
              </Link>
            ))}
          </div>
        </div>
      </section>


    </div>
  );
};

export default Home;
