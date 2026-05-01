import { Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

// Layout는 인증된 라우트의 공통 헤더(현재 사용자 + 로그아웃)와 콘텐츠 영역을 제공한다.
export default function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function onLogout() {
    logout();
    navigate('/login');
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <header className="flex items-center justify-between border-b bg-white px-6 py-3">
        <h1 className="font-bold">wiki-permissions / 1-acl</h1>
        <div className="flex items-center gap-3 text-sm">
          <span>{user?.email}</span>
          <button onClick={onLogout} className="rounded border px-2 py-1">
            로그아웃
          </button>
        </div>
      </header>
      <main className="mx-auto max-w-3xl p-6">
        <Outlet />
      </main>
    </div>
  );
}
