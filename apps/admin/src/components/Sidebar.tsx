import { Link, useLocation } from 'react-router-dom';
import { cn } from '@/lib/utils';
import {
  Rocket,
  Calendar,
  Newspaper,
  MapPin,
  Building2,
  Home,
} from 'lucide-react';

const menuItems = [
  { path: '/', label: 'Dashboard', icon: Home },
  { path: '/rockets', label: 'Rockets', icon: Rocket },
  { path: '/launches', label: 'Launches', icon: Calendar },
  { path: '/news', label: 'News', icon: Newspaper },
  { path: '/bases', label: 'Launch Bases', icon: MapPin },
  { path: '/companies', label: 'Companies', icon: Building2 },
];

export function Sidebar() {
  const location = useLocation();

  return (
    <div className="flex h-screen w-64 flex-col border-r bg-card">
      <div className="flex h-16 items-center border-b px-6">
        <h1 className="text-xl font-bold">LaunchDate Admin</h1>
      </div>
      <nav className="flex-1 space-y-1 p-4">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname === item.path;
          
          return (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
              )}
            >
              <Icon className="h-5 w-5" />
              {item.label}
            </Link>
          );
        })}
      </nav>
    </div>
  );
}
