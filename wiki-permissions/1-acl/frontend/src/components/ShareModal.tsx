import { useEffect, useState } from 'react';
import { apiClient } from '../api/client';

interface ACLEntry {
  id: number;
  page_id: number;
  user_id: number;
  action: 'read' | 'edit';
}

interface Props {
  pageId: number;
  onClose: () => void;
}

// ShareModal은 페이지 owner가 다른 사용자에게 read/edit을 부여하거나 회수하는 UI다.
// 모든 ACL 변경은 서버 API를 통해 검증되며, 본 컴포넌트는 owner가 아닐 때 진입 자체가 차단된다(PageDetailPage에서 가드).
export default function ShareModal({ pageId, onClose }: Props) {
  const [entries, setEntries] = useState<ACLEntry[]>([]);
  const [userId, setUserId] = useState('');
  const [action, setAction] = useState<'read' | 'edit'>('read');
  const [error, setError] = useState('');

  function load() {
    apiClient
      .get<ACLEntry[]>(`/api/pages/${pageId}/acl`)
      .then((r) => setEntries(r.data))
      .catch((e) => setError(e.response?.data?.message ?? 'failed'));
  }

  useEffect(load, [pageId]);

  async function grant() {
    try {
      await apiClient.post(`/api/pages/${pageId}/acl`, { user_id: Number(userId), action });
      setUserId('');
      load();
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  async function revoke(uid: number) {
    try {
      await apiClient.delete(`/api/pages/${pageId}/acl/${uid}`);
      load();
    } catch (e: unknown) {
      const msg = (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  return (
    <div className="fixed inset-0 z-10 flex items-center justify-center bg-black/30">
      <div className="w-full max-w-md rounded-lg bg-white p-6 shadow">
        <div className="mb-3 flex items-center justify-between">
          <h3 className="font-semibold">공유 관리 — page #{pageId}</h3>
          <button onClick={onClose} className="text-slate-500">
            ✕
          </button>
        </div>
        {error && <p className="mb-2 text-sm text-red-600">{error}</p>}

        <ul className="mb-4 max-h-48 divide-y overflow-auto rounded border">
          {entries.map((e) => (
            <li key={e.id} className="flex items-center justify-between p-2 text-sm">
              <span>
                user #{e.user_id} — {e.action}
              </span>
              <button onClick={() => revoke(e.user_id)} className="text-red-600">
                회수
              </button>
            </li>
          ))}
          {entries.length === 0 && (
            <li className="p-2 text-sm text-slate-500">공유된 사용자 없음</li>
          )}
        </ul>

        <div className="flex gap-2">
          <input
            placeholder="user id"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            className="w-24 rounded border p-1"
          />
          <select
            value={action}
            onChange={(e) => setAction(e.target.value as 'read' | 'edit')}
            className="rounded border p-1"
          >
            <option value="read">read</option>
            <option value="edit">edit</option>
          </select>
          <button onClick={grant} className="rounded bg-blue-600 px-3 text-white">
            부여
          </button>
        </div>
      </div>
    </div>
  );
}
