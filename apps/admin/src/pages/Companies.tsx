import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Companies() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate('/companies/prod');
  }, [navigate]);

  return null;
}
