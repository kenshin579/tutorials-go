import { createContext, useContext, useState, type ReactNode } from 'react';
import { apiClient } from '../api/client';

export interface User {
  id: number;
  email: string;
  name: string;
  // 1편(ACL)과 차이: RBAC에서는 사전 게이팅(PermissionGate)에 사용하기 위해
  // login 응답에서 받은 permissions / roles를 클라이언트 상태로 보관한다.
  permissions: string[];
  roles: { id: number; name: string }[];
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
}

interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

// AuthProvider는 로그인 상태/토큰/권한 정보를 localStorage에 보존하며
// login/logout 동작을 자식 트리에 노출한다.
export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(() => {
    const saved = localStorage.getItem('user');
    const token = localStorage.getItem('access_token');
    if (saved && token) {
      return { user: JSON.parse(saved), isAuthenticated: true };
    }
    return { user: null, isAuthenticated: false };
  });

  async function login(email: string, password: string) {
    const res = await apiClient.post('/auth/login', { email, password });
    const { token, user, permissions, roles } = res.data;
    const merged: User = { ...user, permissions, roles };
    localStorage.setItem('access_token', token);
    localStorage.setItem('user', JSON.stringify(merged));
    setState({ user: merged, isAuthenticated: true });
  }

  function logout() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
    setState({ user: null, isAuthenticated: false });
  }

  return (
    <AuthContext.Provider value={{ ...state, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// useAuth는 AuthProvider 자식 어디서든 인증 상태/액션을 가져오는 hook이다.
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
