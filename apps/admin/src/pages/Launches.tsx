import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Launches() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate('/launches/prod', { replace: true });
  }, [navigate]);

  return null;
}
