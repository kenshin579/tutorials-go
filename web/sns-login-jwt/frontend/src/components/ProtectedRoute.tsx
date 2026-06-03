import { Navigate } from 'react-router-dom';
import type { ReactNode } from 'react';

interface Props {
  isAuthenticated: boolean;
  children: ReactNode;
}

export function ProtectedRoute({ isAuthenticated, children }: Props) {
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}
