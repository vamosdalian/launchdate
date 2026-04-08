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
import { useApi } from '../hooks/useApi';
import { fetchRocketLaunch } from '../services/launchesService';

const LaunchDetail = () => {
  const { id } = useParams<{ id: string }>();
  
  const fetchLaunchCallback = useCallback(() => fetchRocketLaunch(id!), [id]);
  const { data: launch, loading: launchLoading, error: launchError } = useApi(fetchLaunchCallback);

  // Countdown state
  const [countdown, setCountdown] = useState({
    days: 0,
    hours: 0,
    minutes: 0,
    seconds: 0,
  });

  // Update countdown every second
  useEffect(() => {
    if (!launch || launch.status !== 1) return; // 1 is scheduled

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

  if (launchLoading) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading launch details...</p>
        </div>
      </div>
    );
  }

  if (launchError || !launch) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
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

  const getStatusBadge = (status: number) => {
    switch (status) {
      case 3: // Successful
        return <Badge className="bg-green-600">✅ Success</Badge>;
      case 4: // Failed
        return <Badge className="bg-red-600">❌ Failed</Badge>;
      case 2: // Cancelled (Assuming 2 is cancelled based on typical status codes, adjust if needed)
        return <Badge className="bg-gray-600">🚫 Cancelled</Badge>;
      default:
        return <Badge className="bg-blue-600">🕒 Scheduled</Badge>;
    }
  };

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-[#f0f0f0]">
      {/* Hero Section */}
      <section className="relative h-[60vh] md:h-[70vh] flex items-end justify-center text-center overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img 
            src={launch.background_image || "https://images.unsplash.com/photo-1614728263952-84ea256ec346?q=80&w=2574&auto=format&fit=crop"}
            onError={(e) => {
              e.currentTarget.src = 'https://images.unsplash.com/photo-1516849841032-87cbac4d88f7?q=80&w=2574&auto=format&fit=crop';
            }}
            alt="Rocket Launch Background" 
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-black/40"></div>
          <div className="absolute inset-x-0 bottom-0 h-1/2" style={{
            background: 'linear-gradient(0deg, rgba(10, 10, 10, 1) 0%, rgba(10, 10, 10, 0.7) 30%, rgba(10, 10, 10, 0) 100%)'
          }}></div>
        </div>
        <div className="relative z-10 p-4 pb-16 md:pb-24 container mx-auto">
          <div className="mb-8">
            {getStatusBadge(launch.status)}
          </div>

          {/* Countdown Timer - Only show for scheduled launches */}
          {launch.status === 1 && (
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
          <p className="text-xl md:text-2xl text-gray-300 mb-2">{launch.rocket_info?.name}</p>
          <p className="text-lg md:text-xl text-gray-400 mb-2">📍 {launch.location_info?.name}</p>
          {launch.status !== 1 && (
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
                {launch.rocket_info && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Rocket Information</h3>
                    <div className="flex flex-col sm:flex-row gap-6">
                      <div className="w-32 h-48 bg-[#0a0a0a] rounded-lg overflow-hidden flex-shrink-0">
                        <img src={launch.rocket_info.thumb_image} alt={launch.rocket_info.name} className="w-full h-full object-cover" />
                      </div>
                      <div className="flex-1">
                        <h4 className="text-2xl font-bold mb-2">{launch.rocket_info.name}</h4>
                        <Button asChild className="bg-blue-600 hover:bg-blue-700">
                          <Link to={`/rockets/${launch.rocket_info.id}`}>View Rocket Details →</Link>
                        </Button>
                      </div>
                    </div>
                  </div>
                )}

                {/* Mission Info */}
                {launch.mission_info && launch.mission_info.length > 0 && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Mission Details</h3>
                    {launch.mission_info.map((mission) => (
                      <div key={mission.id} className="mb-4 last:mb-0">
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
                {launch.location_info && typeof launch.location_info.lat === 'number' && typeof launch.location_info.lon === 'number' && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                    <h3 className="text-xl font-bold mb-3">Location</h3>
                    <div className="aspect-video w-full rounded-lg overflow-hidden">
                      <iframe
                        width="100%"
                        height="100%"
                        frameBorder="0"
                        scrolling="no"
                        marginHeight={0}
                        marginWidth={0}
                        src={`https://www.openstreetmap.org/export/embed.html?bbox=${launch.location_info.lon-0.5}%2C${launch.location_info.lat-0.5}%2C${launch.location_info.lon+0.5}%2C${launch.location_info.lat+0.5}&amp;layer=mapnik&amp;marker=${launch.location_info.lat}%2C${launch.location_info.lon}`}
                      ></iframe>
                    </div>
                    <p className="mt-2 text-gray-400 text-sm">{launch.location_info.name}</p>
                  </div>
                )}
              </div>

              {/* Right Column */}
              <div className="space-y-6">
                {/* Agency Logo */}
                {launch.agency_info && (
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg flex flex-col items-center text-center">
                    <h3 className="text-xl font-bold mb-4">Launch Agency</h3>
                    {launch.agency_info.thumb_image ? (
                      <img src={launch.agency_info.thumb_image} alt={launch.agency_info.name} className="max-w-[200px] max-h-[200px] object-contain mb-4" />
                    ) : (
                      <div className="w-32 h-32 bg-gray-800 rounded-full flex items-center justify-center mb-4">
                        <span className="text-4xl">🚀</span>
                      </div>
                    )}
                    <h4 className="text-lg font-semibold">{launch.agency_info.name}</h4>
                  </div>
                )}

                {/* Timeline */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] p-6 rounded-lg">
                  <h3 className="text-xl font-bold mb-6">Launch Timeline</h3>
                  <div className="relative pl-8">
                    {launch.timeline_event && launch.timeline_event.map((event, index) => (
                      <div key={index} className="relative pb-8 last:pb-0">
                        <div 
                          className="absolute -left-[30px] top-[5px] w-5 h-5 rounded-full bg-[#2a2a2a] border-4 border-[#007bff]"
                        ></div>
                        {index !== launch.timeline_event.length - 1 && (
                          <div 
                            className="absolute -left-[21px] top-[5px] bottom-[-5px] w-0.5 bg-[#2a2a2a]"
                          ></div>
                        )}
                        <p className="font-bold text-lg">{event.relative_time} ({event.abbrev})</p>
                        <p className="text-gray-400">{event.description}</p>
                      </div>
                    ))}
                    {(!launch.timeline_event || launch.timeline_event.length === 0) && (
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
