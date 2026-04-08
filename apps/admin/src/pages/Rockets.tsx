import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Rockets() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate("/rockets/prod", { replace: true });
  }, [navigate]);

  return null;
}
