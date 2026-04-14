import { useParams, Link } from 'react-router-dom';
import { useCallback } from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"
import { useApi } from '../hooks/useApi';
import { fetchRocket } from '../services/rocketsService';

const RocketDetail = () => {
  const { id } = useParams<{ id: string }>();
  
  const fetchRocketCallback = useCallback(() => fetchRocket(id!), [id]);
  const { data: rocket, loading, error } = useApi(fetchRocketCallback);

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading rocket details...</p>
        </div>
      </div>
    );
  }

  if (error || !rocket) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8 max-w-md text-center">
          <h1 className="text-2xl font-bold mb-4 text-white">Rocket Not Found</h1>
          <p className="text-gray-400 mb-6">{error?.message || 'The requested rocket does not exist.'}</p>
          <Button asChild className="bg-blue-600 hover:bg-blue-700">
            <Link to="/rockets">Back to Rockets</Link>
          </Button>
        </div>
      </div>
    );
  }

  // Helper to format currency
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumSignificantDigits: 3 }).format(value);
  };

  // Helper to format large numbers
  const formatNumber = (value: number) => {
    return new Intl.NumberFormat('en-US').format(value);
  };

  return (
    <div className="min-h-screen bg-background text-[#f0f0f0]">
      {/* Hero Section - Two Column Layout */}
      <section className="container mx-auto px-4 py-16 md:py-24">
        <div className="max-w-7xl mx-auto">
          <div className="grid md:grid-cols-2 gap-12 items-start">
            {/* Left: Rocket Image */}
            <div>
              <img 
                src={rocket.main_image} 
                alt={rocket.name} 
                className="w-full h-auto rounded-lg shadow-2xl"
              />
            </div>

            {/* Right: Rocket Info */}
            <div className="space-y-6">
              <div>
                <Link to="/rockets" className="text-sm text-gray-500 hover:text-gray-300 mb-2 inline-block transition-colors">
                  ← Back to List
                </Link>
                <h1 className="text-5xl md:text-6xl font-black text-white mb-2">
                  {rocket.name}
                </h1>
                {rocket.company && (
                  <p className="text-2xl text-gray-400 mb-4">{rocket.company.name}</p>
                )}
                <div className="flex gap-2">
                  {rocket.active && (
                    <Badge className="bg-green-600 text-white px-4 py-1 text-sm">✓ Active</Badge>
                  )}
                  {rocket.reusable && (
                    <Badge className="bg-blue-600 text-white px-4 py-1 text-sm">↻ Reusable</Badge>
                  )}
                </div>
              </div>

              <p className="text-lg text-gray-300 leading-relaxed">
                {rocket.description}
              </p>

              {/* Quick Stats Grid */}
              <div className="grid grid-cols-2 gap-4 pt-4">
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-4">
                  <p className="text-sm text-gray-400 mb-1">Height</p>
                  <p className="text-3xl font-bold text-white">{rocket.length}m</p>
                </div>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-4">
                  <p className="text-sm text-gray-400 mb-1">Diameter</p>
                  <p className="text-3xl font-bold text-white">{rocket.diameter}m</p>
                </div>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-4">
                  <p className="text-sm text-gray-400 mb-1">Launch Mass</p>
                  <p className="text-3xl font-bold text-white">{(rocket.launch_mass / 1000).toFixed(0)}t</p>
                </div>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-4">
                  <p className="text-sm text-gray-400 mb-1">Liftoff Thrust</p>
                  <p className="text-3xl font-bold text-white">{formatNumber(rocket.liftoff_thrust)} kN</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Technical Details Sections */}
      <section className="container mx-auto px-4 py-12">
        <div className="max-w-7xl mx-auto space-y-12">
        
          {/* Overview Section */}
        <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8">
          <h3 className="text-2xl font-bold text-white mb-6">Technical Specifications</h3>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div className="space-y-1">
                <p className="text-sm text-gray-400">Launch Cost</p>
                <p className="text-lg text-white font-semibold">{rocket.launch_cost ? formatCurrency(rocket.launch_cost) : 'N/A'}</p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-400">Success Rate</p>
                <p className="text-lg text-white font-semibold">
                  {rocket.total_launches > 0 
                    ? `${((rocket.success_launches / rocket.total_launches) * 100).toFixed(1)}%` 
                    : 'N/A'}
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-400">LEO Capacity</p>
                <p className="text-lg text-white font-semibold">{rocket.leo_capacity ? `${formatNumber(rocket.leo_capacity)} kg` : 'N/A'}</p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-400">GTO Capacity</p>
                <p className="text-lg text-white font-semibold">{rocket.gto_capacity ? `${formatNumber(rocket.gto_capacity)} kg` : 'N/A'}</p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-400">GEO Capacity</p>
                <p className="text-lg text-white font-semibold">{rocket.geo_capacity ? `${formatNumber(rocket.geo_capacity)} kg` : 'N/A'}</p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-gray-400">SSO Capacity</p>
                <p className="text-lg text-white font-semibold">{rocket.sso_capacity ? `${formatNumber(rocket.sso_capacity)} kg` : 'N/A'}</p>
              </div>
          </div>
        </div>

        {/* Launch History Section */}
        {rocket.launches && rocket.launches.length > 0 && (
          <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8">
            <h3 className="text-2xl font-bold text-white mb-6">Launch Records</h3>
            {Object.entries(
              rocket.launches.reduce((acc, launch) => {
                const year = new Date(launch.launch_time).getFullYear();
                if (!acc[year]) acc[year] = [];
                acc[year].push(launch);
                return acc;
              }, {} as Record<string, typeof rocket.launches>)
            )
            .sort(([yearA], [yearB]) => Number(yearB) - Number(yearA))
            .map(([year, launches]) => (
              <div key={year} className="mb-8 last:mb-0">
                <h4 className="text-xl font-bold text-gray-400 mb-4 border-b border-gray-800 pb-2">{year}</h4>
                <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {launches.map((launch) => (
                    <Link 
                      key={launch.id} 
                      to={`/launches/${launch.id}`}
                      className="flex bg-[#0a0a0a] border border-[#2a2a2a] rounded-lg overflow-hidden hover:border-blue-500 transition-colors group h-24"
                    >
                      <div className="w-32 relative shrink-0">
                        <img 
                          src={launch.thumb_image || 'https://images.unsplash.com/photo-1517976487492-5750f3195933?q=80&w=2670&auto=format&fit=crop'} 
                          alt={launch.name}
                          className="w-full h-full object-cover"
                        />
                      </div>
                      <div className="flex-1 p-3 flex flex-col justify-center min-w-0">
                        <div className="flex justify-between items-start mb-1">
                          <h4 className="font-bold text-white text-sm group-hover:text-blue-400 transition-colors truncate pr-2">{launch.name}</h4>
                          <Badge className={`text-[10px] px-1.5 py-0 h-5 shrink-0 ${launch.status === 'success' ? "bg-green-600" : launch.status === 'failure' ? "bg-red-600" : "bg-blue-600"}`}>
                            {launch.status === 'success' ? 'Success' : launch.status === 'failure' ? 'Failure' : launch.status_label}
                          </Badge>
                        </div>
                        <p className="text-xs text-gray-400 mb-1">{new Date(launch.launch_time).toLocaleDateString()}</p>
                        <p className="text-xs text-gray-500 truncate">{launch.launch_base.name}</p>
                      </div>
                    </Link>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Image Gallery Section */}
        {rocket.image_list && rocket.image_list.length > 0 && (
          <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8">
            <h3 className="text-2xl font-bold text-white mb-6">Gallery</h3>
            <div className="flex justify-center">
              <Carousel
                opts={{
                  align: "start",
                }}
                className="w-full max-w-6xl"
              >
                <CarouselContent>
                  {rocket.image_list.map((img, index) => (
                    <CarouselItem key={index} className="md:basis-1/2 lg:basis-1/3">
                      <div className="p-1">
                        <div className="aspect-video relative rounded-lg overflow-hidden bg-[#0a0a0a]">
                          <img 
                            src={img} 
                            alt={`${rocket.name} gallery image ${index + 1}`} 
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

        </div>
      </section>
    </div>
  );
};

export default RocketDetail;
