import { Link } from 'react-router-dom';
import { useState, useMemo, useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { fetchRockets } from '../services/rocketsService';

const Rockets = () => {
  const [searchTerm, setSearchTerm] = useState('');

  const fetchRocketsCallback = useCallback(() => fetchRockets(), []);
  const { data: rockets, loading, error } = useApi(fetchRocketsCallback);

  const filteredRockets = useMemo(() => {
    if (!rockets) return [];
    return rockets.filter((rocket) => {
      return rocket.name.toLowerCase().includes(searchTerm.toLowerCase());
    });
  }, [rockets, searchTerm]);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading rockets...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
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
    <div className="min-h-screen bg-[#0a0a0a]">
      {/* Page Hero */}
      <section className="py-16 md:py-24 text-center bg-[#111]">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">Rockets</h1>
            <p className="max-w-3xl mx-auto text-lg md:text-xl text-gray-400">
              Explore the rockets that are shaping the future of space exploration
            </p>
          </div>
        </div>
      </section>

      {/* Filters and Search */}
      <section className="py-8 sticky top-20 bg-black/80 backdrop-blur-md z-40">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <div className="flex flex-col md:flex-row gap-4">
              <div className="relative flex-grow">
                <input
                  type="text"
                  placeholder="Search rockets, e.g. 'Falcon 9'..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full bg-[#2a2a2a] border border-[#4a4a4a] rounded-lg py-3 px-4 pl-10 focus:outline-none focus:ring-2 focus:ring-blue-500 text-white"
                />
                <svg
                  className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <circle cx="11" cy="11" r="8"></circle>
                  <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                </svg>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-7xl mx-auto">
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
              {filteredRockets.map((rocket) => (
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
            {filteredRockets.length === 0 && (
              <div className="text-center py-12">
                <p className="text-gray-400 text-lg">No rockets found matching your criteria.</p>
              </div>
            )}
          </div>
        </div>
      </section>
    </div>
  );
};

export default Rockets;
