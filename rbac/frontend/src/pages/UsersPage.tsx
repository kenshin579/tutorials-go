import { useEffect, useState } from 'react';
import apiClient from '../api/client';
import PermissionGate from '../components/PermissionGate';

interface Role {
  id: number;
  name: string;
  description: string;
}

interface UserItem {
  id: number;
  name: string;
  email: string;
  roles: Role[];
}

const roleBadgeColors: Record<string, string> = {
  admin: 'bg-red-100 text-red-700',
  manager: 'bg-blue-100 text-blue-700',
  user: 'bg-green-100 text-green-700',
};

export default function UsersPage() {
  const [users, setUsers] = useState<UserItem[]>([]);
  const [allRoles, setAllRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [editUser, setEditUser] = useState<UserItem | null>(null);
  const [selectedRoleIds, setSelectedRoleIds] = useState<Set<number>>(new Set());
  const [saving, setSaving] = useState(false);

  const fetchData = async () => {
    try {
      const [usersRes, rolesRes] = await Promise.all([
        apiClient.get('/users'),
        apiClient.get('/roles'),
      ]);
      setUsers(usersRes.data ?? []);
      setAllRoles(rolesRes.data ?? []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const openEdit = (user: UserItem) => {
    setEditUser(user);
    setSelectedRoleIds(new Set(user.roles.map((r) => r.id)));
  };

  const toggleRole = (roleId: number) => {
    setSelectedRoleIds((prev) => {
      const next = new Set(prev);
      if (next.has(roleId)) {
        next.delete(roleId);
      } else {
        next.add(roleId);
      }
      return next;
    });
  };

  const handleSave = async () => {
    if (!editUser) return;
    setSaving(true);
    try {
      const currentRoleIds = new Set(editUser.roles.map((r) => r.id));
      // Remove roles that were unchecked
      for (const roleId of currentRoleIds) {
        if (!selectedRoleIds.has(roleId)) {
          await apiClient.delete(`/users/${editUser.id}/roles/${roleId}`);
        }
      }
      // Add roles that were checked
      for (const roleId of selectedRoleIds) {
        if (!currentRoleIds.has(roleId)) {
          await apiClient.post(`/users/${editUser.id}/roles`, { role_id: roleId });
        }
      }
      setEditUser(null);
      fetchData();
    } catch {
      // ignore
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (userId: number) => {
    if (!confirm('Are you sure you want to delete this user?')) return;
    try {
      await apiClient.delete(`/users/${userId}`);
      fetchData();
    } catch {
      // ignore
    }
  };

  if (loading) {
    return <div className="text-gray-500">Loading...</div>;
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Users</h2>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Name</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Email</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Roles</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {users.map((u) => (
              <tr key={u.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm text-gray-800 font-medium">{u.name}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{u.email}</td>
                <td className="px-6 py-4">
                  <div className="flex gap-1.5">
                    {u.roles.map((r) => (
                      <span
                        key={r.id}
                        className={`px-2 py-0.5 rounded-full text-xs font-semibold ${roleBadgeColors[r.name] ?? 'bg-gray-100 text-gray-700'}`}
                      >
                        {r.name}
                      </span>
                    ))}
                  </div>
                </td>
                <td className="px-6 py-4 space-x-2">
                  <PermissionGate permission="users:update">
                    <button
                      onClick={() => openEdit(u)}
                      className="px-3 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition-colors"
                    >
                      Edit
                    </button>
                  </PermissionGate>
                  <PermissionGate permission="users:delete">
                    <button
                      onClick={() => handleDelete(u.id)}
                      className="px-3 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition-colors"
                    >
                      Delete
                    </button>
                  </PermissionGate>
                </td>
              </tr>
            ))}
            {users.length === 0 && (
              <tr>
                <td colSpan={4} className="px-6 py-8 text-center text-sm text-gray-500">
                  No users found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Edit User Modal */}
      {editUser && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">
              Edit User: {editUser.name}
            </h3>
            <div className="space-y-4">
              <div>
                <span className="block text-sm font-medium text-gray-500">Name</span>
                <span className="text-sm text-gray-800">{editUser.name}</span>
              </div>
              <div>
                <span className="block text-sm font-medium text-gray-500">Email</span>
                <span className="text-sm text-gray-800">{editUser.email}</span>
              </div>
              <div>
                <span className="block text-sm font-medium text-gray-700 mb-2">Role Assignment</span>
                <div className="space-y-2 bg-gray-50 rounded-lg p-3">
                  {allRoles.map((role) => (
                    <label key={role.id} className="flex items-center gap-2 text-sm cursor-pointer">
                      <input
                        type="checkbox"
                        checked={selectedRoleIds.has(role.id)}
                        onChange={() => toggleRole(role.id)}
                        className="rounded border-gray-300"
                      />
                      <span className={`px-2 py-0.5 rounded-full text-xs font-semibold ${roleBadgeColors[role.name] ?? 'bg-gray-100 text-gray-700'}`}>
                        {role.name}
                      </span>
                      <span className="text-gray-500">{role.description}</span>
                    </label>
                  ))}
                </div>
              </div>
            </div>
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={() => setEditUser(null)}
                className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                disabled={saving}
                className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                {saving ? 'Saving...' : 'Save'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
