import { Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

// Layout는 인증된 라우트의 공통 헤더와 콘텐츠 영역을 제공한다.
// ABAC에서는 사용자 속성(department/employment_type)을 헤더에 표시해
// 어떤 정책이 적용되는지 사용자가 직관적으로 볼 수 있게 한다.
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
        <h1 className="font-bold">wiki-permissions / 3-abac</h1>
        <div className="flex items-center gap-3 text-sm">
          <span>
            {user?.email}
            {user && (
              <span className="ml-2 inline-flex gap-1">
                <span className="rounded bg-slate-200 px-2 py-0.5 text-xs">
                  {user.department?.name ?? 'no dept'}
                </span>
                <span className="rounded bg-slate-200 px-2 py-0.5 text-xs">
                  {user.employment_type}
                </span>
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
