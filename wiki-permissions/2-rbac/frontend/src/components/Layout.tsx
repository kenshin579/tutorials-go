import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';
import PermissionGate from './PermissionGate';

// Layout는 인증된 라우트의 공통 헤더(현재 사용자 + 로그아웃)와 네비게이션을 제공한다.
// 사용자 관리 메뉴는 users:manage 권한이 있는 사용자(=admin)에게만 노출된다.
export default function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function onLogout() {
    logout();
    navigate('/login');
  }

  const navLink = ({ isActive }: { isActive: boolean }) =>
    `rounded px-2 py-1 text-sm ${isActive ? 'bg-slate-900 text-white' : 'text-slate-700 hover:bg-slate-100'}`;

  return (
    <div className="min-h-screen bg-slate-50">
      <header className="flex items-center justify-between border-b bg-white px-6 py-3">
        <div className="flex items-center gap-4">
          <h1 className="font-bold">wiki-permissions / 2-rbac</h1>
          <nav className="flex items-center gap-2">
            <NavLink to="/pages" className={navLink}>
              페이지
            </NavLink>
            <PermissionGate permission="users:manage">
              <NavLink to="/users" className={navLink}>
                사용자 관리
              </NavLink>
            </PermissionGate>
          </nav>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span>
            {user?.email}
            {user && user.roles.length > 0 && (
              <span className="ml-2 rounded bg-slate-200 px-2 py-0.5 text-xs">
                {user.roles.map((r) => r.name).join(', ')}
              </span>
            )}
          </span>
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
