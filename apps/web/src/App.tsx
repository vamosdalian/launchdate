import { BrowserRouter as Router } from 'react-router-dom';
import { Routes, Route, useLocation } from 'react-router-dom';
import { useEffect } from 'react';
import Navigation from './components/Navigation';
import Footer from './components/Footer';
// import ComingSoon from './pages/ComingSoon';
import Home from './pages/Home';
import Rockets from './pages/Rockets';
import RocketDetail from './pages/RocketDetail';
import Launches from './pages/Launches';
import LaunchDetail from './pages/LaunchDetail';
import News from './pages/News';
import NewsDetail from './pages/NewsDetail';
import LaunchBases from './pages/LaunchBases';
import LaunchBaseDetail from './pages/LaunchBaseDetail';
import Companies from './pages/Companies';
import CompanyDetail from './pages/CompanyDetail';
import './App.css';

function AppContent() {
  const location = useLocation();
  // Don't show global footer on LaunchDetail pages (they have their own)
  const showFooter = !location.pathname.match(/^\/launches\/\d+$/);

  // Scroll to top on route change
  useEffect(() => {
    // Only scroll if not already at the top to avoid unnecessary operations
    if (window.scrollY > 0) {
      window.scrollTo(0, 0);
    }
  }, [location.pathname]);

  return (
    <div className="app">
      <Navigation />
      <main className="main-content">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/rockets" element={<Rockets />} />
          <Route path="/rockets/:id" element={<RocketDetail />} />
          <Route path="/launches" element={<Launches />} />
          <Route path="/launches/:id" element={<LaunchDetail />} />
          <Route path="/news" element={<News />} />
          <Route path="/news/:id" element={<NewsDetail />} />
          <Route path="/bases" element={<LaunchBases />} />
          <Route path="/bases/:id" element={<LaunchBaseDetail />} />
          <Route path="/companies" element={<Companies />} />
          <Route path="/companies/:id" element={<CompanyDetail />} />
        </Routes>
      </main>
      {showFooter && <Footer />}
    </div>
  );
}

function App() {
  // Show Coming Soon page without navigation and footer
  return <AppContent />;
}

function AppWrapper() {
  return (
    <Router>
      <App />
    </Router>
  );
}

export default AppWrapper;
