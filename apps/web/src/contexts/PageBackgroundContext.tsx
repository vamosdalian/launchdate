import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';

import { fetchPageBackgrounds } from '../services/pageBackgroundsService';
import type { PageBackground, PageBackgroundKey } from '../types';

type PageBackgroundMap = Partial<Record<PageBackgroundKey, PageBackground>>;

interface PageBackgroundContextValue {
  backgrounds: PageBackgroundMap;
  loading: boolean;
}

const PageBackgroundContext = createContext<PageBackgroundContextValue>({
  backgrounds: {},
  loading: true,
});

export function PageBackgroundProvider({ children }: { children: ReactNode }) {
  const [backgrounds, setBackgrounds] = useState<PageBackgroundMap>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;

    const loadBackgrounds = async () => {
      try {
        setLoading(true);
        const data = await fetchPageBackgrounds();
        if (!active) {
          return;
        }

        const nextBackgrounds = data.reduce<PageBackgroundMap>((accumulator, background) => {
          accumulator[background.page_key] = background;
          return accumulator;
        }, {});

        setBackgrounds(nextBackgrounds);
      } catch (error) {
        console.error('Failed to load page backgrounds', error);
        if (active) {
          setBackgrounds({});
        }
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    };

    loadBackgrounds();

    return () => {
      active = false;
    };
  }, []);

  return (
    <PageBackgroundContext.Provider value={{ backgrounds, loading }}>
      {children}
    </PageBackgroundContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function usePageBackground(pageKey: PageBackgroundKey): string {
  const { backgrounds } = useContext(PageBackgroundContext);
  return backgrounds[pageKey]?.background_image ?? '';
}

// eslint-disable-next-line react-refresh/only-export-components
export function usePageBackgroundState() {
  return useContext(PageBackgroundContext);
}