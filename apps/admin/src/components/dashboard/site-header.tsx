import { Search } from 'lucide-react';

import UserNav from '@/components/UserNav';
import { Input } from '@/components/ui/input';
import { Separator } from '../ui/separator';
import { SidebarTrigger } from '../ui/sidebar';

export function SiteHeader() {
  return (
    <header className="flex h-16 shrink-0 items-center border-b bg-background px-6">
      <div className="flex w-full items-center gap-4">
        <SidebarTrigger className="hidden lg:flex" />
        <div className="hidden lg:block">
          <Separator orientation="vertical" className="h-6" />
        </div>
        <div className="grid flex-1 gap-1">
          <p className="text-lg font-medium">Control Center</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input placeholder="Search" className="w-64 pl-10" />
          </div>
          <UserNav />
        </div>
      </div>
    </header>
  );
}
