import { useEffect, useState } from 'react';
import { apiClient } from '../api/client';

interface RoleSummary {
  id: number;
  name: string;
  description?: string;
}

interface UserRow {
  id: number;
  email: string;
  name: string;
  roles?: RoleSummary[];
}

// UsersPage는 admin이 사용자 목록을 보고 role을 부여/회수하는 화면이다.
// users:manage 권한이 없는 사용자가 직접 URL로 접근하면 서버가 403 응답 → 화면이 비어 보인다.
// (메뉴는 PermissionGate로 가려져 있고, 이 페이지 자체는 단순 fetch로 처리.)
export default function UsersPage() {
  const [users, setUsers] = useState<UserRow[]>([]);
  const [roles, setRoles] = useState<RoleSummary[]>([]);
  const [error, setError] = useState('');

  function load() {
    Promise.all([
      apiClient.get<UserRow[]>('/api/users'),
      apiClient.get<RoleSummary[]>('/api/roles'),
    ])
      .then(([uRes, rRes]) => {
        setUsers(uRes.data);
        setRoles(rRes.data);
      })
      .catch((e) => setError(e.response?.data?.message ?? 'failed'));
  }

  useEffect(load, []);

  async function assign(userId: number, roleId: number) {
    if (!roleId) return;
    try {
      await apiClient.post(`/api/users/${userId}/roles`, { role_id: roleId });
      load();
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  async function revoke(userId: number, roleId: number) {
    try {
      await apiClient.delete(`/api/users/${userId}/roles/${roleId}`);
      load();
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  return (
    <div>
      <h2 className="mb-4 text-lg font-semibold">사용자 관리 (admin 전용)</h2>
      {error && <p className="mb-2 text-red-600">{error}</p>}

      <table className="w-full rounded bg-white text-sm shadow">
        <thead className="border-b text-left">
          <tr>
            <th className="p-3">Email</th>
            <th className="p-3">이름</th>
            <th className="p-3">현재 Role</th>
            <th className="p-3">관리</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id} className="border-b">
              <td className="p-3">{u.email}</td>
              <td className="p-3">{u.name}</td>
              <td className="p-3">
                <div className="flex flex-wrap gap-1">
                  {(u.roles ?? []).map((r) => (
                    <span
                      key={r.id}
                      className="inline-flex items-center gap-1 rounded bg-slate-100 px-2 py-0.5 text-xs"
                    >
                      {r.name}
                      <button
                        onClick={() => revoke(u.id, r.id)}
                        className="text-red-600"
                        title="회수"
                      >
                        ×
                      </button>
                    </span>
                  ))}
                </div>
              </td>
              <td className="p-3">
                <select
                  defaultValue=""
                  onChange={(e) => {
                    const v = Number(e.target.value);
                    if (v) assign(u.id, v);
                    e.currentTarget.value = '';
                  }}
                  className="rounded border p-1 text-xs"
                >
                  <option value="">+ role 추가</option>
                  {roles.map((r) => (
                    <option key={r.id} value={r.id}>
                      {r.name}
                    </option>
                  ))}
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
