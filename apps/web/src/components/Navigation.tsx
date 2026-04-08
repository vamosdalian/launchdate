import { Link, useLocation } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useTheme } from '@/hooks/useTheme';

const Navigation = () => {
  const location = useLocation();
  const { theme, toggleTheme } = useTheme();

  const isActive = (path: string) => location.pathname === path;

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-black/50 backdrop-blur-md">
      <div className="container mx-auto px-4">
        <div className="flex h-20 items-center justify-between">
          <Link to="/" className="flex items-center gap-2 text-2xl font-bold">
            <img src="/rocket.png" alt="Rocket" className="w-8 h-8" />
            LaunchDate
          </Link>
          <div className="hidden md:flex items-center gap-8">
            <Link 
              to="/launches" 
              className={`font-medium transition-colors ${isActive('/launches') ? 'text-white' : 'text-gray-300 hover:text-white'}`}
            >
              Launches
            </Link>
            <Link 
              to="/rockets" 
              className={`font-medium transition-colors ${isActive('/rockets') ? 'text-white' : 'text-gray-300 hover:text-white'}`}
            >
              Rockets
            </Link>
            <Link 
              to="/news" 
              className={`font-medium transition-colors ${isActive('/news') ? 'text-white' : 'text-gray-300 hover:text-white'}`}
            >
              News
            </Link>
            <Link 
              to="/companies" 
              className={`font-medium transition-colors ${isActive('/companies') ? 'text-white' : 'text-gray-300 hover:text-white'}`}
            >
              Companies
            </Link>
            <Link 
              to="/bases" 
              className={`font-medium transition-colors ${isActive('/bases') ? 'text-white' : 'text-gray-300 hover:text-white'}`}
            >
              Locations
            </Link>
          </div>
          <div className="md:hidden">
            <Button 
              variant="ghost" 
              size="icon"
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

