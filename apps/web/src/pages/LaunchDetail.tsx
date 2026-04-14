import { useParams, Link } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"
import { useEffect, useState, useCallback } from 'react';
import { CircleMarker, MapContainer, TileLayer, Tooltip } from 'react-leaflet';
import { useApi } from '../hooks/useApi';
import { fetchRocketLaunch } from '../services/launchesService';

const LaunchDetail = () => {
  const { id } = useParams<{ id: string }>();
  
  const fetchLaunchCallback = useCallback(() => fetchRocketLaunch(id!), [id]);
  const { data: launch, loading: launchLoading, error: launchError } = useApi(fetchLaunchCallback);
  const [companyImageFailed, setCompanyImageFailed] = useState(false);
  const [rocketImageFailed, setRocketImageFailed] = useState(false);

  // Countdown state
  const [countdown, setCountdown] = useState({
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0,
  });

  // Update countdown every second
  useEffect(() => {
    if (!launch || launch.status !== 'scheduled') return;

    const updateCountdown = () => {
      const launchDate = new Date(launch.launch_time).getTime();
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
  }, [launch]);

  useEffect(() => {
    setCompanyImageFailed(false);
  }, [launch?.company?.id, launch?.company?.image_url]);

  useEffect(() => {
    setRocketImageFailed(false);
  }, [launch?.rocket?.id, launch?.rocket?.image_url, launch?.rocket?.thumb_image]);

  if (launchLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading launch details...</p>
        </div>
      </div>
    );
  }

  if (launchError || !launch) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8 max-w-md text-center">
          <h1 className="text-2xl font-bold mb-4">Launch Not Found</h1>
          <p className="text-gray-400 mb-6">{launchError?.message || 'The requested launch does not exist.'}</p>
          <Button asChild className="bg-blue-600 hover:bg-blue-700">
            <Link to="/launches">Back to Launches</Link>
          </Button>
        </div>
      </div>
    );
  }

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

  const getStatusBadge = (status: string, statusLabel: string) => {
    switch (status) {
      case 'success':
        return <Badge className="bg-green-600">✅ Success</Badge>;
      case 'failure':
        return <Badge className="bg-red-600">❌ Failed</Badge>;
      case 'cancelled':
        return <Badge className="bg-gray-600">🚫 Cancelled</Badge>;
      case 'delayed':
        return <Badge className="bg-amber-600">⏸ Delayed</Badge>;
      default:
        return <Badge className="bg-blue-600">{statusLabel}</Badge>;
    }
  };

  const launchBaseCenter: [number, number] | null =
    launch.launch_base &&
    typeof launch.launch_base.latitude === 'number' &&
    typeof launch.launch_base.longitude === 'number'
      ? [launch.launch_base.latitude, launch.launch_base.longitude]
      : null;

  const companyHref = launch.company?.id ? `/companies/${launch.company.id}` : null;
  const companyInitials = launch.company?.name
    ? launch.company.name
        .split(/\s+/)
        .filter(Boolean)
        .slice(0, 2)
        .map((segment) => segment.charAt(0).toUpperCase())
        .join('')
    : 'LD';

  return (
    <div className="min-h-screen bg-background text-[#f0f0f0]">
      {/* Hero Section */}
      <section className="relative -mt-20 h-screen flex items-end justify-center text-center overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img 
            src={launch.background_image || "https://images.unsplash.com/photo-1614728263952-84ea256ec346?q=80&w=2574&auto=format&fit=crop"}
            onError={(e) => {
              e.currentTarget.src = 'https://images.unsplash.com/photo-1516849841032-87cbac4d88f7?q=80&w=2574&auto=format&fit=crop';
            }}
            alt="Rocket Launch Background" 
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-black/25"></div>
          <div className="absolute inset-x-0 bottom-0 h-1/2" style={{
            background: 'linear-gradient(0deg, rgba(10, 10, 10, 0.9) 0%, rgba(10, 10, 10, 0.45) 30%, rgba(10, 10, 10, 0) 100%)'
          }}></div>
        </div>
        <div className="relative z-10 p-4 pb-24 md:pb-32 container mx-auto">
          <div className="mb-8">
            {getStatusBadge(launch.status, launch.status_label)}
          </div>

          {/* Countdown Timer - Only show for scheduled launches */}
          {launch.status === 'scheduled' && (
            <div className="mb-8">
              <h2 className="text-lg md:text-xl font-medium text-gray-300 mb-4">Launch Countdown</h2>
              <div className="flex justify-center gap-4 md:gap-6">
                <div className="rounded-lg px-4 py-3 md:px-6 md:py-4 min-w-[70px] md:min-w-[100px]" style={{
                  background: 'rgba(255, 255, 255, 0.05)',
                  backdropFilter: 'blur(10px)',
                  WebkitBackdropFilter: 'blur(10px)',
                  border: '1px solid rgba(255, 255, 255, 0.1)'
                }}>
                  <div className="text-3xl md:text-5xl font-bold">{String(countdown.days).padStart(2, '0')}</div>
                  <div className="text-xs md:text-sm text-gray-400 mt-1">Days</div>
                </div>
                <div className="rounded-lg px-4 py-3 md:px-6 md:py-4 min-w-[70px] md:min-w-[100px]" style={{
                  background: 'rgba(255, 255, 255, 0.05)',
                  backdropFilter: 'blur(10px)',
                  WebkitBackdropFilter: 'blur(10px)',
                  border: '1px solid rgba(255, 255, 255, 0.1)'
                }}>
                  <div className="text-3xl md:text-5xl font-bold">{String(countdown.hours).padStart(2, '0')}</div>
                  <div className="text-xs md:text-sm text-gray-400 mt-1">Hours</div>
                </div>
                <div className="rounded-lg px-4 py-3 md:px-6 md:py-4 min-w-[70px] md:min-w-[100px]" style={{
                  background: 'rgba(255, 255, 255, 0.05)',
                  backdropFilter: 'blur(10px)',
                  WebkitBackdropFilter: 'blur(10px)',
                  border: '1px solid rgba(255, 255, 255, 0.1)'
                }}>
                  <div className="text-3xl md:text-5xl font-bold">{String(countdown.minutes).padStart(2, '0')}</div>
                  <div className="text-xs md:text-sm text-gray-400 mt-1">Minutes</div>
                </div>
                <div className="rounded-lg px-4 py-3 md:px-6 md:py-4 min-w-[70px] md:min-w-[100px]" style={{
                  background: 'rgba(255, 255, 255, 0.05)',
                  backdropFilter: 'blur(10px)',
                  WebkitBackdropFilter: 'blur(10px)',
                  border: '1px solid rgba(255, 255, 255, 0.1)'
                }}>
                  <div className="text-3xl md:text-5xl font-bold">{String(countdown.seconds).padStart(2, '0')}</div>
                  <div className="text-xs md:text-sm text-gray-400 mt-1">Seconds</div>
                </div>
              </div>
            </div>
          )}

          <h1 className="text-4xl md:text-6xl font-extrabold mb-4">{launch.name}</h1>
          <p className="text-xl md:text-2xl text-gray-300 mb-2">{launch.rocket.name}</p>
          <p className="text-lg md:text-xl text-gray-400 mb-2">📍 {launch.launch_base.name}</p>
          {launch.status !== 'scheduled' && (
            <p className="text-lg md:text-xl text-gray-400 mb-8">🕒 {formatDate(launch.launch_time)}</p>
          )}
        </div>
      </section>

      {/* Main Content */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <Button asChild variant="ghost" className="mb-4 text-gray-300 hover:text-white">
              <Link to="/launches">← Back to Launches</Link>
            </Button>
            <div className="grid lg:grid-cols-3 gap-8">
              {/* Left Column */}
              <div className="lg:col-span-2 space-y-6">

                {/* Rocket Information */}
                {launch.rocket && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Rocket Information</h3>
                    <div className="flex flex-col sm:flex-row gap-6">
                      <div className="w-32 h-48 bg-[#0a0a0a] rounded-lg overflow-hidden flex-shrink-0">
                        {!rocketImageFailed && (launch.rocket.image_url || launch.rocket.thumb_image) ? (
                          <img
                            src={launch.rocket.image_url || launch.rocket.thumb_image}
                            alt={launch.rocket.name}
                            className="w-full h-full object-cover"
                            onError={() => setRocketImageFailed(true)}
                          />
                        ) : (
                          <div className="flex h-full w-full items-center justify-center bg-[radial-gradient(circle_at_top,_rgba(56,189,248,0.18),_rgba(8,15,28,0.92)_72%)] px-3 text-center">
                            <span className="text-sm font-semibold uppercase tracking-[0.28em] text-slate-300">
                              {launch.rocket.name}
                            </span>
                          </div>
                        )}
                      </div>
                      <div className="flex-1">
                        <h4 className="text-2xl font-bold mb-2">{launch.rocket.name}</h4>
                        {launch.rocket.id && (
                          <Button asChild className="bg-blue-600 hover:bg-blue-700">
                            <Link to={`/rockets/${launch.rocket.id}`}>View Rocket Details →</Link>
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                )}

                {/* Mission Info */}
                {launch.missions && launch.missions.length > 0 && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Mission Details</h3>
                    {launch.missions.map((mission) => (
                      <div key={mission.name} className="mb-4 last:mb-0">
                        <h4 className="font-semibold text-lg mb-2">{mission.name}</h4>
                        {mission.description && <p className="text-gray-400">{mission.description}</p>}
                      </div>
                    ))}
                  </div>
                )}

                {/* Image List */}
                {launch.image_list && launch.image_list.length > 0 && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Images</h3>
                    <div className="flex justify-center px-12">
                      <Carousel
                        opts={{
                          align: "start",
                        }}
                        className="w-full"
                      >
                        <CarouselContent>
                          {launch.image_list.map((img, idx) => (
                            <CarouselItem key={idx} className="md:basis-1/2 lg:basis-1/2">
                              <div className="p-1">
                                <div className="aspect-video relative rounded-lg overflow-hidden bg-[#0a0a0a]">
                                  <img 
                                    src={img} 
                                    alt={`Mission image ${idx + 1}`} 
                                    className="w-full h-full object-cover" 
                                  />
                                </div>
                              </div>
                            </CarouselItem>
                          ))}
                        </CarouselContent>
                        <CarouselPrevious />
                        <CarouselNext />
                      </Carousel>
                    </div>
                  </div>
                )}

                {/* Location Map */}
                {launchBaseCenter && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Location</h3>
                    <div className="aspect-video w-full overflow-hidden rounded-lg border border-white/10">
                      <MapContainer
                        center={launchBaseCenter}
                        zoom={7}
                        scrollWheelZoom={true}
                        className="h-full w-full"
                      >
                        <TileLayer
                          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        />
                        <CircleMarker
                          center={launchBaseCenter}
                          radius={10}
                          pathOptions={{
                            color: '#eef4fb',
                            weight: 2,
                            fillColor: '#78d7ff',
                            fillOpacity: 0.95,
                          }}
                        >
                          <Tooltip direction="top" offset={[0, -10]} opacity={1} permanent>
                            {launch.launch_base.name}
                          </Tooltip>
                        </CircleMarker>
                      </MapContainer>
                    </div>
                    <p className="mt-2 text-gray-400 text-sm">{launch.launch_base.name}</p>
                  </div>
                )}
              </div>

              {/* Right Column */}
              <div className="space-y-6">
                {/* Agency Logo */}
                {launch.company && (
                  companyHref ? (
                    <Link to={companyHref} className="group block">
                      <div className="flex flex-col items-center rounded-lg border border-[#2a2a2a] bg-[#1a1a1a] p-6 text-center transition-all duration-300 hover:-translate-y-1 hover:border-[var(--color-card-hover-border)]">
                        <h3 className="mb-4 text-xl font-bold">Launch Agency</h3>
                        {launch.company.image_url && !companyImageFailed ? (
                          <div className="mb-4 flex h-32 w-32 items-center justify-center overflow-hidden rounded-3xl border border-white/10 bg-white/5 p-4">
                            <img
                              src={launch.company.image_url}
                              alt={launch.company.name}
                              className="max-h-full max-w-full object-contain"
                              onError={() => setCompanyImageFailed(true)}
                            />
                          </div>
                        ) : (
                          <div className="mb-4 flex h-32 w-32 items-center justify-center rounded-3xl border border-white/10 bg-[linear-gradient(180deg,rgba(120,215,255,0.16),rgba(255,255,255,0.03))] text-3xl font-black tracking-[0.08em] text-[var(--color-highlight-strong)] shadow-[0_18px_50px_rgba(0,0,0,0.18)]">
                            {companyInitials}
                          </div>
                        )}
                        <h4 className="text-lg font-semibold transition-colors group-hover:text-[var(--color-highlight-strong)]">{launch.company.name}</h4>
                      </div>
                    </Link>
                  ) : (
                    <div className="flex flex-col items-center rounded-lg border border-[#2a2a2a] bg-[#1a1a1a] p-6 text-center">
                      <h3 className="mb-4 text-xl font-bold">Launch Agency</h3>
                      {launch.company.image_url && !companyImageFailed ? (
                        <div className="mb-4 flex h-32 w-32 items-center justify-center overflow-hidden rounded-3xl border border-white/10 bg-white/5 p-4">
                          <img
                            src={launch.company.image_url}
                            alt={launch.company.name}
                            className="max-h-full max-w-full object-contain"
                            onError={() => setCompanyImageFailed(true)}
                          />
                        </div>
                      ) : (
                        <div className="mb-4 flex h-32 w-32 items-center justify-center rounded-3xl border border-white/10 bg-[linear-gradient(180deg,rgba(120,215,255,0.16),rgba(255,255,255,0.03))] text-3xl font-black tracking-[0.08em] text-[var(--color-highlight-strong)] shadow-[0_18px_50px_rgba(0,0,0,0.18)]">
                          {companyInitials}
                        </div>
                      )}
                      <h4 className="text-lg font-semibold">{launch.company.name}</h4>
                    </div>
                  )
                )}

                {/* Timeline */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                  <h3 className="text-xl font-bold mb-6">Launch Timeline</h3>
                  <div className="relative pl-8">
                    {launch.timeline && launch.timeline.map((event, index) => (
                      <div key={index} className="relative pb-8 last:pb-0">
                        <div 
                          className="absolute -left-[30px] top-[5px] w-5 h-5 rounded-full bg-[#2a2a2a] border-4 border-[#007bff]"
                        ></div>
                        {index !== launch.timeline.length - 1 && (
                          <div 
                            className="absolute -left-[21px] top-[5px] bottom-[-5px] w-0.5 bg-[#2a2a2a]"
                          ></div>
                        )}
                        <p className="font-bold text-lg">{event.relative_time} ({event.abbrev})</p>
                        <p className="text-gray-400">{event.description}</p>
                      </div>
                    ))}
                    {(!launch.timeline || launch.timeline.length === 0) && (
                      <p className="text-gray-400">No timeline information available.</p>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-[#111] border-t border-gray-800">
        <div className="container mx-auto px-4 py-12">
          <div className="grid md:grid-cols-3 gap-8">
            <div>
              <h3 className="text-xl font-bold mb-2">LaunchDate</h3>
              <p className="text-gray-400">Your comprehensive source for rocket launches, space news, and aerospace information.</p>
            </div>
            <div>
              <h4 className="font-semibold text-lg mb-3">Quick Links</h4>
              <ul className="space-y-2 text-gray-400">
                <li><Link to="/launches" className="hover:text-white">Launch Dates</Link></li>
                <li><Link to="/rockets" className="hover:text-white">Rockets</Link></li>
                <li><Link to="/news" className="hover:text-white">News</Link></li>
                <li><Link to="/companies" className="hover:text-white">Companies</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-lg mb-3">Stay Connected</h4>
              <p className="text-gray-400 mb-4">Subscribe for the latest launch updates and space news.</p>
              <form className="flex">
                <input 
                  type="email" 
                  placeholder="Your email" 
                  className="w-full rounded-l-md bg-gray-800 border-gray-700 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button 
                  type="submit" 
                  className="bg-blue-600 hover:bg-blue-700 text-white font-semibold px-4 rounded-r-md"
                >
                  Subscribe
                </button>
              </form>
            </div>
          </div>
          <div className="mt-12 border-t border-gray-800 pt-8 text-center text-gray-500 text-sm">
            <p>&copy; 2025 LaunchDate. All Rights Reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default LaunchDetail;
