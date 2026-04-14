import { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Rockets from './pages/Rockets';
import Launches from './pages/Launches';
import News from './pages/News';
import LaunchBases from './pages/LaunchBases';
import Companies from './pages/Companies';
import LoginPage from './pages/Login';
import MobileWarningPage from './pages/Mobile';
import { SidebarInset, SidebarProvider } from './components/ui/sidebar';
import { AppSidebar } from './components/dashboard/app-sidebar';
import { SiteHeader } from './components/dashboard/site-header';
import { ScrollArea } from './components/ui/scroll-area';
import { useIsMobile } from './hooks/use-mobile';
import { clearAuthState, getAccessToken } from './lib/authStore';
import { getCurrentUser, refreshAccessToken } from './services/authService';
import ProdRockets from './pages/rockets/ProdRockets';
import LL2Launchers from './pages/rockets/LL2Launchers';
import LL2LauncherFamilies from './pages/rockets/LL2LauncherFamilies';
import ProdLaunches from './pages/launches/ProdLaunches';
import LL2Launches from './pages/launches/LL2Launches';
import ProdLaunchBases from './pages/launch-bases/ProdLaunchBases';
import LL2Locations from './pages/launch-bases/LL2Locations';
import LL2LaunchPads from './pages/launch-bases/LL2LaunchPads';
import ProdCompanies from './pages/companies/ProdCompanies';
import LL2Companies from './pages/companies/LL2Companies';
import Images from './pages/Images';
import PageBackgrounds from './pages/page-backgrounds/PageBackgrounds';

const ProtectedRoute: React.FC = () => {
  const [status, setStatus] = useState<'pending' | 'authorized' | 'unauthorized'>('pending');

  useEffect(() => {
    let active = true;

    const ensureSession = async () => {
      let token = getAccessToken();
      if (!token) {
        token = await refreshAccessToken();
      }

      if (!token) {
        if (!active) return;
        clearAuthState();
        setStatus('unauthorized');
        return;
      }

  const user = await getCurrentUser(true);
      if (!active) return;

      if (user) {
        setStatus('authorized');
      } else {
        clearAuthState();
        setStatus('unauthorized');
      }
    };

    ensureSession();

    return () => {
      active = false;
    };
  }, []);

  if (status === 'pending') return null;
  if (status === 'unauthorized') return <Navigate to="/login" replace />;
  return <Outlet />;
};

const MobileRedirector: React.FC = () => {
  const isMobile = useIsMobile();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    if (isMobile && location.pathname !== '/mobile') {
      navigate('/mobile', { replace: true });
    }
    if (!isMobile && location.pathname === '/mobile') {
      navigate('/', { replace: true });
    }
  }, [isMobile, location.pathname, navigate]);

  return null;
};

const AppLayout: React.FC = () => (
  <SidebarProvider>
    <div className="flex min-h-screen w-full bg-muted/40">
      <AppSidebar />
      <SidebarInset>
        <SiteHeader />
        <ScrollArea className="flex-1 min-h-0 h-full">
          <div className="px-6 py-6">
            <Outlet />
          </div>
        </ScrollArea>
      </SidebarInset>
    </div>
  </SidebarProvider>
);

function App() {
  return (
    <Router>
      <MobileRedirector />
      <Routes>
        <Route path="/mobile" element={<MobileWarningPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/rockets" element={<Rockets />} />
            <Route path="/rockets/prod" element={<ProdRockets />} />
            <Route path="/rockets/ll2-launchers" element={<LL2Launchers />} />
            <Route path="/rockets/ll2-launcher-families" element={<LL2LauncherFamilies />} />
            <Route path="/launches" element={<Launches />} />
            <Route path="/launches/prod" element={<ProdLaunches />} />
            <Route path="/launches/ll2" element={<LL2Launches />} />
            <Route path="/news" element={<News />} />
            <Route path="/launch-bases" element={<LaunchBases />} />
            <Route path="/launch-bases/prod" element={<ProdLaunchBases />} />
            <Route path="/launch-bases/ll2-locations" element={<LL2Locations />} />
            <Route path="/launch-bases/ll2-pads" element={<LL2LaunchPads />} />
            <Route path="/companies" element={<Companies />} />
            <Route path="/companies/prod" element={<ProdCompanies />} />
            <Route path="/companies/ll2" element={<LL2Companies />} />
            <Route path="/images" element={<Images />} />
            <Route path="/page-backgrounds" element={<PageBackgrounds />} />
          </Route>
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
