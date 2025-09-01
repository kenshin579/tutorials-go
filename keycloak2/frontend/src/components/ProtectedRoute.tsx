import React from 'react';
import { Navigate } from 'react-router-dom';
import keycloak from '../services/keycloak';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  if (keycloak.authenticated !== true) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
