import { Link } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { fetchLaunchBases } from '../services/launchBasesService';

const LaunchBases = () => {
  const fetchLaunchBasesCallback = useCallback(() => fetchLaunchBases(), []);
  const { data: launchBases, loading, error } = useApi(fetchLaunchBasesCallback);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading launch bases...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
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
    <div className="min-h-screen bg-[#0a0a0a]">
      {/* Page Hero */}
      <section className="py-16 md:py-24 text-center bg-[#111]">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">Launch Sites</h1>
          <p className="max-w-3xl mx-auto text-lg md:text-xl text-gray-400">
            Explore the facilities where rockets are launched into space
          </p>
        </div>
      </section>

      {/* Launch Bases Grid */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {launchBases && launchBases.map((base) => (
              <Link key={base.id} to={`/bases/${base.id}`}>
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg overflow-hidden hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300 cursor-pointer h-full flex flex-col">
                  <div className="aspect-video bg-[#0a0a0a]">
                    <img 
                      src={base.imageUrl} 
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
        </div>
      </section>
    </div>
  );
};

export default LaunchBases;
