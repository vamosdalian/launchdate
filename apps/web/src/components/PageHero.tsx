import { useEffect, useState } from 'react';

import { cn } from '@/lib/utils';

import { usePageBackground } from '../contexts/PageBackgroundContext';
import type { PageBackgroundKey } from '../types';

interface PageHeroProps {
  pageKey: PageBackgroundKey;
  title: string;
  description: string;
  className?: string;
}

export default function PageHero({ pageKey, title, description, className }: PageHeroProps) {
  const backgroundImage = usePageBackground(pageKey);
  const [showImage, setShowImage] = useState(Boolean(backgroundImage));

  useEffect(() => {
    setShowImage(Boolean(backgroundImage));
  }, [backgroundImage]);

  return (
    <section
      className={cn(
        'relative -mt-20 flex min-h-[60vh] items-end justify-center overflow-hidden bg-[#111] text-center md:min-h-[68vh]',
        className,
      )}
    >
      {showImage ? (
        <div className="absolute inset-0">
          <img
            src={backgroundImage}
            alt={title}
            className="h-full w-full object-cover"
            onError={() => setShowImage(false)}
          />
          <div className="absolute inset-0 bg-black/18" />
          <div className="absolute inset-0" style={{ background: 'radial-gradient(circle at top, rgba(120, 215, 255, 0.16) 0%, transparent 38%)' }} />
          <div className="absolute inset-x-0 bottom-0 h-1/2" style={{ background: 'linear-gradient(180deg, rgba(5, 11, 20, 0) 0%, rgba(5, 11, 20, 0.78) 100%)' }} />
        </div>
      ) : (
        <div
          className="absolute inset-0"
          style={{
            background: 'radial-gradient(circle at top, rgba(120, 215, 255, 0.22) 0%, transparent 30%), radial-gradient(circle at 78% 20%, rgba(255, 179, 107, 0.14) 0%, transparent 18%), linear-gradient(180deg, rgba(10, 17, 29, 0.84) 0%, rgba(5, 11, 20, 0.98) 100%)',
          }}
        />
      )}

      <div className="relative container mx-auto px-4 pb-32 md:pb-40">
        <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-4">{title}</h1>
        <p className="max-w-3xl mx-auto text-lg md:text-xl text-gray-200">{description}</p>
      </div>
    </section>
  );
}