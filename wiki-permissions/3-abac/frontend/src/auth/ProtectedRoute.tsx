import { Navigate } from 'react-router-dom';
import type { ReactNode } from 'react';
import { useAuth } from './AuthContext';

// ProtectedRoute는 미인증 사용자를 /login으로 리다이렉트한다.
// ABAC 권한 평가는 서버가 매 요청마다 수행하며, 클라이언트는 결과(decision)를 받아 표시한다.
export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
