import { Navigate } from 'react-router-dom';
import { useAuth } from './AuthContext';
import { usePermission } from './usePermission';
import type { ReactNode } from 'react';

interface ProtectedRouteProps {
  children: ReactNode;
  permission?: string;
}

export default function ProtectedRoute({ children, permission }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuth();
  const { hasPermission } = usePermission();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (permission && !hasPermission(permission)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}
