import { Link, useLocation } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useTheme } from '@/hooks/useTheme';

const Navigation = () => {
  const location = useLocation();
  const { theme, toggleTheme } = useTheme();
  const navItems = [
    { label: 'Home', path: '/' },
    { label: 'Launches', path: '/launches' },
    { label: 'Rockets', path: '/rockets' },
    { label: 'Companies', path: '/companies' },
    { label: 'Locations', path: '/bases' },
  ];

  const isActive = (path: string) => location.pathname === path;

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-white/10 bg-[#0F2854]/80 backdrop-blur-xl">
      <div className="container mx-auto px-4">
        <div className="flex h-20 items-center justify-between">
          <Link to="/" className="flex items-center gap-3 text-sm font-semibold uppercase tracking-[0.24em] text-[var(--color-page-fg)]">
            <span className="relative flex h-12 w-12 items-center justify-center">
              <span className="absolute inset-1 rounded-full bg-[var(--color-highlight-soft)] blur-lg" />
              <img src="/rocket.png" alt="Rocket" className="relative h-8 w-8" />
            </span>
            <span className="flex items-center gap-3 tracking-[0.18em]">
              LaunchDate
              <span className="inline-flex items-center rounded-[4px] border border-[rgba(255,179,107,0.24)] bg-[rgba(255,179,107,0.1)] px-1.5 py-0.5 text-[8px] font-bold tracking-[0.22em] text-[var(--color-highlight-warm)] shadow-[0_0_14px_rgba(255,179,107,0.14)]">
                LIVE
              </span>
            </span>
          </Link>
          <div className="hidden items-center gap-1 md:flex">
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                className={`px-4 py-2 text-sm font-medium tracking-[0.12em] transition-colors ${
                  isActive(item.path)
                    ? 'text-white'
                    : 'text-gray-300 hover:text-white'
                }`}
              >
                {item.label}
              </Link>
            ))}
          </div>
          <div className="md:hidden">
            <Button 
              variant="ghost" 
              size="icon"
              className="rounded-2xl border border-white/10 bg-[var(--color-icon-button-bg)] text-[var(--color-icon-button-fg)] hover:bg-[var(--color-icon-button-hover-bg)] hover:text-[var(--color-icon-button-fg)]"
              onClick={toggleTheme}
              aria-label="Toggle theme"
            >
              {theme === 'light' ? '🌙' : '☀️'}
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navigation;

