import { useEffect, useState } from 'react';
import './ComingSoon.css';

const ComingSoon = () => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  return (
    <div className="coming-soon-container">
      <div className={`coming-soon-content ${mounted ? 'fade-in' : ''}`}>
        <div className="rocket-animation">
          <div className="rocket">
            <div className="rocket-body">
              <div className="rocket-window"></div>
              <div className="rocket-fin rocket-fin-left"></div>
              <div className="rocket-fin rocket-fin-right"></div>
            </div>
            <div className="rocket-flame">
              <div className="flame"></div>
              <div className="flame flame-2"></div>
              <div className="flame flame-3"></div>
            </div>
          </div>
        </div>

        <h1 className="main-title">Launch Date</h1>
        <h2 className="subtitle">Coming Soon</h2>
        
        <p className="description">
          We're preparing for liftoff! Our mission control center is currently under development.
        </p>
        
        
      </div>
    </div>
  );
};

export default ComingSoon;
