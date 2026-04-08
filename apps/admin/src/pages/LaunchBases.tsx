import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function LaunchBases() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate('/launch-bases/prod', { replace: true });
  }, [navigate]);

  return null;
}
