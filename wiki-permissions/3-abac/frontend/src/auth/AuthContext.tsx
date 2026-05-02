import { createContext, useContext, useState, type ReactNode } from 'react';
import { apiClient } from '../api/client';

export interface Department {
  id: number;
  name: string;
}

export interface User {
  id: number;
  email: string;
  name: string;
  // 1·2편과 차이: ABAC에서는 사용자에게 속성(department/employment_type)이 직접 붙는다.
  // 클라이언트는 이 속성을 단순히 표시 용도로만 사용 — 권한 평가는 항상 서버가 한다.
  department: Department | null;
  employment_type: 'fulltime' | 'contract';
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
    const { token, user } = res.data;
    localStorage.setItem('access_token', token);
    localStorage.setItem('user', JSON.stringify(user));
    setState({ user, isAuthenticated: true });
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

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
