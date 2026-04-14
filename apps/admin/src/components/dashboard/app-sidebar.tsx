import type { ComponentProps } from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  Building2,
  Calendar,
  Home,
  Image as ImageIcon,
  MapPin,
  Newspaper,
  Rocket,
} from 'lucide-react';
import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from '@/components/ui/hover-card';
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarSeparator,
} from '../ui/sidebar';

const navigation = [
  { path: '/', label: 'Dashboard', icon: Home },
  {
    label: 'Rockets',
    icon: Rocket,
    sub: [
      { path: '/rockets/prod', label: 'Prod' },
      { path: '/rockets/ll2-launchers', label: 'LL2 Launchers' },
      { path: '/rockets/ll2-launcher-families', label: 'LL2 Launcher Families' },
    ],
  },
  {
    label: 'Launches',
    icon: Calendar,
    sub: [
      { path: '/launches/prod', label: 'Prod' },
      { path: '/launches/ll2', label: 'LL2' },
    ],
  },
  { path: '/news', label: 'News', icon: Newspaper },
  {
    label: 'Launch Bases',
    icon: MapPin,
    sub: [
      { path: '/launch-bases/prod', label: 'Prod' },
      { path: '/launch-bases/ll2-locations', label: 'LL2 Locations' },
      { path: '/launch-bases/ll2-pads', label: 'LL2 Pads' },
    ],
  },
  {
    label: 'Companies',
    icon: Building2,
    sub: [
      { path: '/companies/prod', label: 'Prod' },
      { path: '/companies/ll2', label: 'LL2' },
    ],
  },
  { path: '/images', label: 'Images', icon: ImageIcon },
  { path: '/page-backgrounds', label: 'Page Backgrounds', icon: ImageIcon },
];

export function AppSidebar(props: ComponentProps<typeof Sidebar>) {
  const location = useLocation();

  return (
    <Sidebar variant="sidebar" collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild className="h-12 text-left text-base font-semibold">
              <Link to="/">
                <span className="truncate">LaunchDate Admin</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarSeparator />
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navigation.map((item) => {
                const Icon = item.icon;

                if ('sub' in item && item.sub) {
                  const isSubActive = item.sub.some((subItem) =>
                    location.pathname.startsWith(subItem.path)
                  );
                  return (
                    <SidebarMenuItem key={item.label}>
                      <HoverCard openDelay={200} closeDelay={100}>
                        <HoverCardTrigger asChild>
                          <SidebarMenuButton
                            asChild
                            isActive={isSubActive}
                            className="w-full justify-start"
                          >
                            <Link to={item.sub[0].path}>
                              <Icon className="h-4 w-4" />
                              <span>{item.label}</span>
                            </Link>
                          </SidebarMenuButton>
                        </HoverCardTrigger>
                        <HoverCardContent
                          side="right"
                          align="start"
                          className="p-1"
                        >
                          <div className="flex flex-col space-y-1">
                            {item.sub.map((subItem) => (
                              <Link
                                key={subItem.path}
                                to={subItem.path}
                                className="rounded-md px-3 py-2 text-sm hover:bg-gray-100"
                              >
                                {subItem.label}
                              </Link>
                            ))}
                          </div>
                        </HoverCardContent>
                      </HoverCard>
                    </SidebarMenuItem>
                  );
                }

                if ('path' in item) {
                  const isActive = location.pathname === item.path;
                  return (
                    <SidebarMenuItem key={item.path}>
                      <SidebarMenuButton asChild isActive={isActive}>
                        <Link to={item.path}>
                          <Icon className="h-4 w-4" />
                          <span>{item.label}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  );
                }
                return null;
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  );
}
