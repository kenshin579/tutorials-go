import { useEffect, useState } from 'react';
import { useAuth } from '../auth/AuthContext';
import apiClient from '../api/client';
import { usePermission } from '../auth/usePermission';

const roleBadgeColors: Record<string, string> = {
  admin: 'bg-red-100 text-red-700',
  manager: 'bg-blue-100 text-blue-700',
  user: 'bg-green-100 text-green-700',
};

export default function DashboardPage() {
  const { user } = useAuth();
  const { hasPermission } = usePermission();
  const [stats, setStats] = useState({ users: 0, roles: 0, permissions: 0 });

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const promises: Promise<unknown>[] = [];

        if (hasPermission('users:read')) {
          promises.push(apiClient.get('/users').then((r) => ({ users: r.data?.length ?? 0 })));
        } else {
          promises.push(Promise.resolve({ users: 0 }));
        }

        if (hasPermission('roles:read')) {
          promises.push(apiClient.get('/roles').then((r) => ({ roles: r.data?.length ?? 0 })));
          promises.push(
            apiClient.get('/permissions').then((r) => ({ permissions: r.data?.length ?? 0 }))
          );
        } else {
          promises.push(Promise.resolve({ roles: 0 }));
          promises.push(Promise.resolve({ permissions: 0 }));
        }

        const results = await Promise.all(promises);
        const merged = Object.assign({}, ...results as Record<string, number>[]);
        setStats(merged);
      } catch {
        // ignore errors for stats
      }
    };

    fetchStats();
  }, [hasPermission]);

  const statCards = [
    { label: 'Users', value: stats.users, color: 'bg-blue-500' },
    { label: 'Roles', value: stats.roles, color: 'bg-green-500' },
    { label: 'Permissions', value: stats.permissions, color: 'bg-purple-500' },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h2>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        {statCards.map((card) => (
          <div key={card.label} className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <p className="text-sm font-medium text-gray-500">{card.label}</p>
            <p className="text-3xl font-bold text-gray-800 mt-1">{card.value}</p>
          </div>
        ))}
      </div>

      {/* My Info */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">My Info</h3>
        <div className="space-y-3">
          <div className="flex">
            <span className="text-sm font-medium text-gray-500 w-24">Name:</span>
            <span className="text-sm text-gray-800">{user?.name}</span>
          </div>
          <div className="flex">
            <span className="text-sm font-medium text-gray-500 w-24">Email:</span>
            <span className="text-sm text-gray-800">{user?.email}</span>
          </div>
          <div className="flex items-center">
            <span className="text-sm font-medium text-gray-500 w-24">Roles:</span>
            <div className="flex gap-2">
              {user?.roles.map((role) => (
                <span
                  key={role}
                  className={`px-2.5 py-0.5 rounded-full text-xs font-semibold ${roleBadgeColors[role] ?? 'bg-gray-100 text-gray-700'}`}
                >
                  {role}
                </span>
              ))}
            </div>
          </div>
          <div className="flex">
            <span className="text-sm font-medium text-gray-500 w-24">Permissions:</span>
            <div className="flex flex-wrap gap-1.5">
              {user?.permissions.map((perm) => (
                <span
                  key={perm}
                  className="px-2 py-0.5 bg-gray-100 text-gray-600 rounded text-xs"
                >
                  {perm}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
