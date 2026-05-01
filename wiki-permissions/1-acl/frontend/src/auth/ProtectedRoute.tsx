import { Navigate } from 'react-router-dom';
import type { ReactNode } from 'react';
import { useAuth } from './AuthContext';

// ProtectedRoute는 미인증 사용자를 /login으로 리다이렉트한다.
// 권한 단위(action) 게이팅은 별도 컴포넌트(예: PermissionGate) 책임이며 여기서는 인증만 검증한다.
export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
