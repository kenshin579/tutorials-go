import { useEffect, useState } from 'react';
import apiClient from '../api/client';
import PermissionGate from '../components/PermissionGate';

interface Permission {
  id: number;
  resource: string;
  action: string;
  description: string;
}

interface Role {
  id: number;
  name: string;
  description: string;
  permissions: Permission[];
}

export default function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [editRole, setEditRole] = useState<Role | null>(null);
  const [selectedPermIds, setSelectedPermIds] = useState<Set<number>>(new Set());
  const [formName, setFormName] = useState('');
  const [formDesc, setFormDesc] = useState('');
  const [isCreate, setIsCreate] = useState(false);
  const [saving, setSaving] = useState(false);

  const fetchData = async () => {
    try {
      const [rolesRes, permsRes] = await Promise.all([
        apiClient.get('/roles'),
        apiClient.get('/permissions'),
      ]);
      setRoles(rolesRes.data ?? []);
      setAllPermissions(permsRes.data ?? []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const openCreate = () => {
    setIsCreate(true);
    setEditRole(null);
    setFormName('');
    setFormDesc('');
    setSelectedPermIds(new Set());
  };

  const openEdit = (role: Role) => {
    setIsCreate(false);
    setEditRole(role);
    setFormName(role.name);
    setFormDesc(role.description);
    setSelectedPermIds(new Set((role.permissions ?? []).map((p) => p.id)));
  };

  const togglePerm = (permId: number) => {
    setSelectedPermIds((prev) => {
      const next = new Set(prev);
      if (next.has(permId)) {
        next.delete(permId);
      } else {
        next.add(permId);
      }
      return next;
    });
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      if (isCreate) {
        const res = await apiClient.post('/roles', { name: formName, description: formDesc });
        const newRole = res.data;
        // Assign permissions
        for (const permId of selectedPermIds) {
          await apiClient.post(`/roles/${newRole.id}/permissions`, { permission_id: permId });
        }
      } else if (editRole) {
        await apiClient.put(`/roles/${editRole.id}`, { name: formName, description: formDesc });
        const currentPermIds = new Set((editRole.permissions ?? []).map((p) => p.id));
        // Remove unchecked permissions
        for (const permId of currentPermIds) {
          if (!selectedPermIds.has(permId)) {
            await apiClient.delete(`/roles/${editRole.id}/permissions/${permId}`);
          }
        }
        // Add newly checked permissions
        for (const permId of selectedPermIds) {
          if (!currentPermIds.has(permId)) {
            await apiClient.post(`/roles/${editRole.id}/permissions`, { permission_id: permId });
          }
        }
      }
      closeModal();
      fetchData();
    } catch {
      // ignore
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (roleId: number) => {
    if (!confirm('Are you sure you want to delete this role?')) return;
    try {
      await apiClient.delete(`/roles/${roleId}`);
      fetchData();
    } catch {
      // ignore
    }
  };

  const closeModal = () => {
    setEditRole(null);
    setIsCreate(false);
  };

  // Group permissions by resource
  const groupedPerms = allPermissions.reduce<Record<string, Permission[]>>((acc, perm) => {
    if (!acc[perm.resource]) acc[perm.resource] = [];
    acc[perm.resource].push(perm);
    return acc;
  }, {});

  if (loading) {
    return <div className="text-gray-500">Loading...</div>;
  }

  const isModalOpen = isCreate || editRole !== null;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Roles</h2>
        <PermissionGate permission="roles:create">
          <button
            onClick={openCreate}
            className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            + Add Role
          </button>
        </PermissionGate>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Role</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Description</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Permissions</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {roles.map((role) => (
              <tr key={role.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm text-gray-800 font-medium">{role.name}</td>
                <td className="px-6 py-4 text-sm text-gray-600">{role.description}</td>
                <td className="px-6 py-4 text-sm text-gray-800">{(role.permissions ?? []).length}</td>
                <td className="px-6 py-4 space-x-2">
                  <PermissionGate permission="roles:update">
                    <button
                      onClick={() => openEdit(role)}
                      className="px-3 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition-colors"
                    >
                      Edit
                    </button>
                  </PermissionGate>
                  <PermissionGate permission="roles:delete">
                    <button
                      onClick={() => handleDelete(role.id)}
                      className="px-3 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition-colors"
                    >
                      Delete
                    </button>
                  </PermissionGate>
                </td>
              </tr>
            ))}
            {roles.length === 0 && (
              <tr>
                <td colSpan={4} className="px-6 py-8 text-center text-sm text-gray-500">
                  No roles found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Create/Edit Role Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-lg p-6 max-h-[80vh] overflow-y-auto">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">
              {isCreate ? 'Add Role' : `Edit Role: ${editRole?.name}`}
            </h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                <input
                  type="text"
                  value={formName}
                  onChange={(e) => setFormName(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <input
                  type="text"
                  value={formDesc}
                  onChange={(e) => setFormDesc(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Permission Assignment</label>
                <div className="bg-gray-50 rounded-lg p-3 space-y-3 max-h-60 overflow-y-auto">
                  {Object.entries(groupedPerms).map(([resource, perms]) => (
                    <div key={resource}>
                      <p className="text-xs font-semibold text-gray-500 uppercase mb-1">{resource}</p>
                      <div className="space-y-1 pl-2">
                        {perms.map((perm) => (
                          <label key={perm.id} className="flex items-center gap-2 text-sm cursor-pointer">
                            <input
                              type="checkbox"
                              checked={selectedPermIds.has(perm.id)}
                              onChange={() => togglePerm(perm.id)}
                              className="rounded border-gray-300"
                            />
                            <span className="text-gray-700">{perm.resource}:{perm.action}</span>
                            <span className="text-gray-400 text-xs">- {perm.description}</span>
                          </label>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={closeModal}
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
