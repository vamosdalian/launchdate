import { useParams, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useCallback } from 'react';
import { useApi } from '../hooks/useApi';
import { fetchNewsArticle } from '../services/newsService';
import type { News } from '../types';

const NewsDetail = () => {
  const { id } = useParams<{ id: string }>();
  
  const fetchNewsCallback = useCallback(() => fetchNewsArticle(id!), [id]);
  const { data: article, loading, error } = useApi(fetchNewsCallback);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
          <p className="text-gray-400">Loading article...</p>
        </div>
      </div>
    );
  }

  if (error || !article) {
    return (
      <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center">
        <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg p-8 max-w-md text-center">
          <h1 className="text-2xl font-bold mb-4">Article Not Found</h1>
          <p className="text-gray-400 mb-6">{error?.message || 'The requested article does not exist.'}</p>
          <Button asChild className="bg-blue-600 hover:bg-blue-700">
            <Link to="/news">Back to News</Link>
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
    });
  };

  // Get related articles (exclude current one) - will be loaded via API in future
  const relatedArticles: News[] = [];

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      {/* Article Hero */}
      <section className="relative h-[50vh] md:h-[60vh] flex items-end overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img 
            src={article.imageUrl} 
            alt={article.title}
            className="w-full h-full object-cover"
            onError={(e) => {
              e.currentTarget.src = 'https://images.unsplash.com/photo-1630343719250-6c9d0407420c?q=80&w=1200&auto=format&fit=crop';
            }}
          />
          <div className="absolute inset-0 bg-black/60"></div>
          <div className="absolute inset-x-0 bottom-0 h-1/2" style={{
            background: 'linear-gradient(0deg, rgba(10, 10, 10, 1) 0%, rgba(10, 10, 10, 0.7) 30%, rgba(10, 10, 10, 0) 100%)'
          }}></div>
        </div>
        <div className="relative z-10 container mx-auto px-4 pb-12">
          <Badge className="mb-4 bg-blue-600">Space News</Badge>
          <h1 className="text-3xl md:text-5xl font-extrabold mb-4">{article.title}</h1>
          <div className="flex items-center gap-4 text-gray-300">
            <span>{formatDate(article.date)}</span>
          </div>
        </div>
      </section>

      {/* Article Content */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <Button asChild variant="ghost" className="mb-8 text-gray-300 hover:text-white">
              <Link to="/news">← Back to News</Link>
            </Button>

            {/* Article content rendered directly without card wrapper */}
            <article className="prose prose-lg prose-invert max-w-none">
              {article.content ? (
                <ReactMarkdown 
                  remarkPlugins={[remarkGfm]}
                  components={{
                    // Customize heading styles
                    h1: (props) => <h1 className="text-4xl font-extrabold mb-6 mt-8 text-white" {...props} />,
                    h2: (props) => <h2 className="text-3xl font-bold mb-5 mt-7 text-white" {...props} />,
                    h3: (props) => <h3 className="text-2xl font-semibold mb-4 mt-6 text-white" {...props} />,
                    // Customize paragraph styles
                    p: (props) => <p className="text-gray-300 mb-4 leading-relaxed text-lg" {...props} />,
                    // Customize list styles
                    ul: (props) => <ul className="text-gray-300 mb-4 ml-6 list-disc space-y-2" {...props} />,
                    ol: (props) => <ol className="text-gray-300 mb-4 ml-6 list-decimal space-y-2" {...props} />,
                    li: (props) => <li className="text-gray-300 leading-relaxed" {...props} />,
                    // Customize blockquote styles
                    blockquote: (props) => (
                      <blockquote className="border-l-4 border-blue-600 pl-6 py-2 my-6 italic text-gray-400 bg-[#1a1a1a] rounded-r-lg" {...props} />
                    ),
                    // Customize table styles
                    table: (props) => (
                      <div className="overflow-x-auto my-6">
                        <table className="min-w-full border border-[#2a2a2a] rounded-lg" {...props} />
                      </div>
                    ),
                    thead: (props) => <thead className="bg-[#1a1a1a]" {...props} />,
                    th: (props) => <th className="border border-[#2a2a2a] px-4 py-3 text-left font-semibold text-white" {...props} />,
                    td: (props) => <td className="border border-[#2a2a2a] px-4 py-3 text-gray-300" {...props} />,
                    // Customize code styles
                    code: (props) => {
                      const {children, className} = props;
                      const isInline = !className || !className.includes('language-');
                      return isInline ? (
                        <code className="bg-[#1a1a1a] border border-[#2a2a2a] px-2 py-1 rounded text-sm text-blue-400">{children}</code>
                      ) : (
                        <code className="block bg-[#1a1a1a] border border-[#2a2a2a] p-4 rounded-lg text-sm overflow-x-auto">{children}</code>
                      );
                    },
                    // Customize link styles - only open external links in new tab
                    a: ({href, ...props}) => {
                      const isExternal = href && (href.startsWith('http://') || href.startsWith('https://'));
                      return (
                        <a 
                          href={href}
                          className="text-blue-400 hover:text-blue-300 underline transition-colors" 
                          {...(isExternal ? { target: '_blank', rel: 'noopener noreferrer' } : {})}
                          {...props} 
                        />
                      );
                    },
                  }}
                >
                  {article.content}
                </ReactMarkdown>
              ) : (
                <div className="text-gray-300">
                  <p className="text-xl leading-relaxed">{article.summary}</p>
                </div>
              )}
            </article>

            {/* External Link at bottom */}
            {article.url && article.url !== '#' && (
              <div className="mt-12 pt-8 border-t border-[#2a2a2a]">
                <div className="flex items-center justify-between">
                  <p className="text-gray-400">Read the original article from the source:</p>
                  <Button asChild className="bg-blue-600 hover:bg-blue-700">
                    <a href={article.url} target="_blank" rel="noopener noreferrer">
                      Visit Source →
                    </a>
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      </section>

      {/* Related Articles */}
      {relatedArticles.length > 0 && (
        <section className="py-16 bg-[#111]">
          <div className="container mx-auto px-4">
            <h2 className="text-3xl font-bold mb-8">Related Articles</h2>
            <div className="grid md:grid-cols-3 gap-6">
              {relatedArticles.map((related) => (
                <Link key={related.id} to={`/news/${related.id}`}>
                  <div className="bg-[#1a1a1a] border border-[#2a2a2a] rounded-lg overflow-hidden hover:border-[#4a4a4a] hover:-translate-y-1 transition-all duration-300">
                    <div className="aspect-video bg-[#0a0a0a]">
                      <img
                        src={related.imageUrl}
                        alt={related.title}
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          e.currentTarget.src = 'https://images.unsplash.com/photo-1454789548928-9efd52dc4031?q=80&w=800&auto=format&fit=crop';
                        }}
                      />
                    </div>
                    <div className="p-6">
                      <p className="text-sm text-gray-400 mb-2">{formatDate(related.date)}</p>
                      <h3 className="text-xl font-bold mb-3 line-clamp-2">{related.title}</h3>
                      <p className="text-gray-400 line-clamp-2">{related.summary}</p>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </section>
      )}
    </div>
  );
};

export default NewsDetail;
