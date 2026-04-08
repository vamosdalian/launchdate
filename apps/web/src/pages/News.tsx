import { Link } from 'react-router-dom';
import { useCallback, useMemo } from 'react';
import { useApi } from '../hooks/useApi';
import { fetchNews } from '../services/newsService';

const News = () => {
  const fetchNewsCallback = useCallback(() => fetchNews(), []);
  const { data: news, loading, error } = useApi(fetchNewsCallback);
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  // Get featured article (first one)
  const { featuredArticle, regularNews } = useMemo(() => {
    if (!news || news.length === 0) {
      return { featuredArticle: null, regularNews: [] };
    }
    return {
      featuredArticle: news[0],
      regularNews: news.slice(1)
    };
  }, [news]);

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading news...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-lg mb-4">Error loading news: {error.message}</p>
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
          <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">Space News</h1>
          <p className="max-w-3xl mx-auto text-lg md:text-xl text-gray-400">
            Stay informed about the latest developments in space exploration, from industry leaders to innovative newcomers.
          </p>
        </div>
      </section>

      {/* Main Content */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="flex flex-col lg:flex-row gap-12">
            
            {/* News Feed (Left Column) */}
            <div className="w-full lg:w-2/3">
              {/* Featured Article */}
              {featuredArticle && (
                <Link to={`/news/${featuredArticle.id}`}>
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl overflow-hidden mb-12 hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300">
                    <img 
                      src={featuredArticle.imageUrl} 
                      alt={featuredArticle.title} 
                      className="w-full h-auto object-cover"
                      onError={(e) => {
                        e.currentTarget.src = 'https://images.unsplash.com/photo-1630343719250-6c9d0407420c?q=80&w=1200&auto=format&fit=crop';
                      }}
                    />
                    <div className="p-8">
                      <div className="flex items-center gap-4 mb-3">
                        <span className="text-sm font-semibold text-blue-400">Featured Story</span>
                        <span className="text-sm text-gray-400">{formatDate(featuredArticle.date)}</span>
                      </div>
                      <h2 className="text-3xl font-bold mb-4 hover:text-blue-400 transition-colors">
                        {featuredArticle.title}
                      </h2>
                      <p className="text-gray-400 mb-6 text-lg">{featuredArticle.summary}</p>
                      <span className="font-semibold text-blue-400 hover:text-blue-300 transition-colors">
                        Read More →
                      </span>
                    </div>
                  </div>
                </Link>
              )}

              {/* News Grid */}
              <div className="grid md:grid-cols-2 gap-6">
                {regularNews.map((article) => (
                  <Link key={article.id} to={`/news/${article.id}`}>
                    <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl overflow-hidden flex flex-col hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300">
                      <img 
                        src={article.imageUrl} 
                        alt={article.title} 
                        className="w-full h-48 object-cover"
                        onError={(e) => {
                          e.currentTarget.src = 'https://images.unsplash.com/photo-1454789548928-9efd52dc4031?q=80&w=800&auto=format&fit=crop';
                        }}
                      />
                      <div className="p-6 flex flex-col flex-grow">
                        <div className="flex items-center gap-4 mb-2">
                          <span className="text-sm font-semibold text-blue-400">Space News</span>
                          <span className="text-sm text-gray-400">{formatDate(article.date)}</span>
                        </div>
                        <h3 className="text-xl font-bold mb-3 flex-grow hover:text-blue-400 transition-colors">
                          {article.title}
                        </h3>
                        <p className="text-gray-400 text-sm mb-4 line-clamp-3">{article.summary}</p>
                        <span className="font-semibold text-blue-400 hover:text-blue-300 transition-colors">
                          Read More →
                        </span>
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            </div>

            {/* Sidebar (Right Column) */}
            <div className="w-full lg:w-1/3">
              <div className="sticky top-28">
                {/* Popular Tags */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-6 mb-8">
                  <h3 className="text-xl font-bold mb-4">Popular Tags</h3>
                  <div className="flex flex-wrap gap-2">
                    {['Starship', 'Artemis', 'Webb Telescope', 'Commercial Space', 'Mars Exploration', 'Space Station'].map((tag) => (
                      <a 
                        key={tag}
                        href="#" 
                        className="bg-[#2a2a2a] border border-[#4a4a4a] px-3 py-1 rounded-full text-sm hover:bg-[#3a3a3a] hover:border-[#6a6a6a] transition-all"
                      >
                        {tag}
                      </a>
                    ))}
                  </div>
                </div>

                {/* Weekly Hotspot */}
                <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-xl overflow-hidden">
                  <h3 className="text-xl font-bold p-6">This Week's Spotlight</h3>
                  <img 
                    src="https://images.unsplash.com/photo-1605103041133-7a6c62a83f3b?q=80&w=800&auto=format&fit=crop" 
                    alt="Saturn V" 
                    className="w-full h-48 object-cover"
                    onError={(e) => {
                      e.currentTarget.src = 'https://images.unsplash.com/photo-1516849841032-87cbac4d88f7?q=80&w=800&auto=format&fit=crop';
                    }}
                  />
                  <div className="p-6">
                    <h4 className="text-2xl font-bold mb-2">Saturn V</h4>
                    <p className="text-gray-400 text-sm mb-4">Revisiting the legendary rocket that took humanity to the Moon, exploring the engineering marvel and enduring legacy.</p>
                    <a href="#" className="font-semibold text-blue-400 hover:text-blue-300 transition-colors">Explore Archive →</a>
                  </div>
                </div>
              </div>
            </div>

          </div>
        </div>
      </section>
    </div>
  );
};

export default News;
