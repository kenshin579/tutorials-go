import { Outlet } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';
import Sidebar from './Sidebar';

const roleBadgeColors: Record<string, string> = {
  admin: 'bg-red-100 text-red-700',
  manager: 'bg-blue-100 text-blue-700',
  user: 'bg-green-100 text-green-700',
};

export default function Layout() {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6 shadow-sm">
        <h1 className="text-xl font-bold text-gray-800">RBAC Admin</h1>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-gray-700">{user?.name}</span>
            {user?.roles.map((role) => (
              <span
                key={role}
                className={`px-2 py-0.5 rounded-full text-xs font-semibold ${roleBadgeColors[role] ?? 'bg-gray-100 text-gray-700'}`}
              >
                {role}
              </span>
            ))}
          </div>
          <button
            onClick={logout}
            className="px-3 py-1.5 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg transition-colors"
          >
            Logout
          </button>
        </div>
      </header>

      <div className="flex">
        <Sidebar />
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
