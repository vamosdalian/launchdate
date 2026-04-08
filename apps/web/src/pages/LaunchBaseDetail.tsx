import { useParams, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useCallback, useMemo } from 'react';
import { useApi } from '../hooks/useApi';
import { fetchLaunchBase } from '../services/launchBasesService';
import { fetchRocketLaunches } from '../services/launchesService';
import type { Launch } from '../types';

// Helper to get the best date field from launch data
const getLaunchDate = (launch: Launch): string => {
  return launch.t0 || launch.window_open || launch.win_open || launch.date || launch.created_at || '';
};

const LaunchBaseDetail = () => {
  const { id } = useParams<{ id: string }>();
  
  const fetchLaunchBaseCallback = useCallback(() => fetchLaunchBase(id!), [id]);
  const { data: base, loading: baseLoading, error: baseError } = useApi(fetchLaunchBaseCallback);
  
  const fetchLaunchesCallback = useCallback(() => fetchRocketLaunches(), []);
  const { data: launchesData } = useApi(fetchLaunchesCallback);
  const launches = launchesData?.launches;

  // Get launches from this base
  const baseLaunches = useMemo(() => {
    if (!launches || !base) return [];
    return launches.filter(l => l.location.includes(base.name)).slice(0, 6);
  }, [launches, base]);

  if (baseLoading) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading launch base details...</p>
        </div>
      </div>
    );
  }

  if (baseError || !base) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8 max-w-md text-center">
          <h1 className="text-2xl font-bold mb-4">Launch Site Not Found</h1>
          <p className="text-gray-400 mb-6">{baseError?.message || 'The requested launch site does not exist.'}</p>
          <Button asChild className="bg-blue-600 hover:bg-blue-700">
            <Link to="/bases">Back to Launch Sites</Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      {/* Hero Section with Image */}
      <section className="relative h-[50vh] md:h-[60vh] flex items-end overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img 
            src={base.imageUrl} 
            alt={base.name}
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-black/60"></div>
          <div className="absolute inset-x-0 bottom-0 h-1/2" style={{
            background: 'linear-gradient(0deg, rgba(10, 10, 10, 1) 0%, rgba(10, 10, 10, 0.7) 30%, rgba(10, 10, 10, 0) 100%)'
          }}></div>
        </div>
        <div className="relative z-10 container mx-auto px-4 pb-12">
          <Badge className="mb-4 bg-blue-600">📍 Launch Site</Badge>
          <h1 className="text-3xl md:text-5xl font-extrabold mb-4">{base.name}</h1>
          <p className="text-xl md:text-2xl text-gray-300 mb-2">{base.location}</p>
          <p className="text-lg md:text-xl text-gray-400">{base.country}</p>
        </div>
      </section>

      {/* Main Content */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <Button asChild variant="ghost" className="mb-8 text-gray-300 hover:text-white">
              <Link to="/bases">← Back to Launch Sites</Link>
            </Button>

            <div className="grid lg:grid-cols-3 gap-8">
              {/* Main Info */}
              <div className="lg:col-span-2 space-y-6">
                {/* Description */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8">
                  <h2 className="text-3xl font-bold mb-6">About This Launch Site</h2>
                  <p className="text-lg text-gray-300 leading-relaxed mb-6">
                    {base.description}
                  </p>
                  
                  {/* Additional Details */}
                  <div className="grid md:grid-cols-2 gap-6 mt-8 pt-8 border-t border-[#2a2a2a]">
                    <div>
                      <h3 className="text-xl font-bold mb-4">Location Details</h3>
                      <ul className="space-y-3 text-gray-300">
                        <li>
                          <strong className="text-white">Country:</strong> {base.country}
                        </li>
                        <li>
                          <strong className="text-white">Region:</strong> {base.location}
                        </li>
                        <li>
                          <strong className="text-white">Coordinates:</strong><br />
                          {base.latitude.toFixed(4)}°N, {Math.abs(base.longitude).toFixed(4)}°W
                        </li>
                      </ul>
                    </div>
                    <div>
                      <h3 className="text-xl font-bold mb-4">Facilities</h3>
                      <p className="text-gray-300">
                        This launch site features state-of-the-art facilities including launch pads, 
                        integration buildings, and support infrastructure for various rocket systems.
                      </p>
                    </div>
                  </div>
                </div>

                {/* Map Placeholder */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8">
                  <h3 className="text-2xl font-bold mb-4">Location Map</h3>
                  <div className="aspect-video bg-[#0a0a0a] rounded-lg flex items-center justify-center text-gray-400">
                    <div className="text-center">
                      <p className="text-4xl mb-2">🗺️</p>
                      <p>Map View</p>
                      <p className="text-sm mt-2">
                        {base.latitude.toFixed(4)}°N, {Math.abs(base.longitude).toFixed(4)}°W
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Sidebar - Statistics */}
              <div className="space-y-6">
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-6">
                  <h3 className="text-xl font-bold mb-6">Statistics</h3>
                  <div className="space-y-4">
                    <div className="border-b border-[#2a2a2a] pb-4">
                      <p className="text-sm text-gray-400 mb-1">Total Launches</p>
                      <p className="text-3xl font-bold text-blue-500">{baseLaunches.length}</p>
                    </div>
                    <div className="border-b border-[#2a2a2a] pb-4">
                      <p className="text-sm text-gray-400 mb-1">Success Rate</p>
                      <p className="text-3xl font-bold text-green-500">
                        {baseLaunches.length > 0 
                          ? Math.round((baseLaunches.filter(l => l.status === 'successful').length / baseLaunches.length) * 100)
                          : 0}%
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-400 mb-1">Upcoming Launches</p>
                      <p className="text-3xl font-bold text-blue-500">
                        {baseLaunches.filter(l => l.status === 'scheduled').length}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Coordinates Card */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-6">
                  <h3 className="text-xl font-bold mb-4">Coordinates</h3>
                  <div className="space-y-2">
                    <Badge variant="secondary" className="w-full justify-center py-2 bg-[#2a2a2a]">
                      Lat: {base.latitude.toFixed(4)}°
                    </Badge>
                    <Badge variant="secondary" className="w-full justify-center py-2 bg-[#2a2a2a]">
                      Long: {base.longitude.toFixed(4)}°
                    </Badge>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Launches from this Base */}
      {baseLaunches.length > 0 && (
        <section className="py-20 bg-[#111]">
          <div className="container mx-auto px-4">
            <div className="flex justify-between items-center mb-8">
              <h2 className="text-3xl font-bold">Launches from {base.name}</h2>
              <Button asChild variant="outline">
                <Link to="/launches">View All</Link>
              </Button>
            </div>
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {baseLaunches.map((launch) => {
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

                const getStatusBadge = (status: string) => {
                  switch (status) {
                    case 'scheduled':
                      return <Badge>Scheduled</Badge>;
                    case 'successful':
                      return <Badge className="bg-green-600">✅ Success</Badge>;
                    case 'failed':
                      return <Badge className="bg-red-600">❌ Failed</Badge>;
                    default:
                      return <Badge variant="outline">{status}</Badge>;
                  }
                };

                return (
                  <Link key={launch.id} to={`/launches/${launch.id}`}>
                    <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-6 hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer h-full">
                      <div className="flex justify-between items-start mb-4">
                        <h3 className="text-xl font-bold">{launch.name}</h3>
                        {getStatusBadge(launch.status)}
                      </div>
                      <div className="space-y-2 text-sm text-gray-400 mb-4">
                        <p>🚀 {launch.vehicle?.name || launch.rocket}</p>
                        <p>📅 {launch.date_str || formatDate(getLaunchDate(launch))} UTC</p>
                      </div>
                      <p className="text-gray-300 line-clamp-2">{launch.quicktext || launch.launch_description || launch.mission_description || launch.description}</p>
                    </div>
                  </Link>
                );
              })}
            </div>
          </div>
        </section>
      )}
    </div>
  );
};

export default LaunchBaseDetail;
