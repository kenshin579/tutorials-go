import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { apiClient } from '../api/client';
import PermissionGate from '../components/PermissionGate';

interface Page {
  id: number;
  title: string;
  content: string;
  owner_id: number;
}

// PageDetailPage는 페이지 상세 조회/편집/삭제 화면이다.
// 1편(ACL)과 차이: 편집/삭제 버튼은 PermissionGate로 사전 게이팅되어, 권한 없는 사용자에겐 아예 노출되지 않는다.
//
// 한계 — RBAC만으로는 "내가 만든 페이지만 편집"을 표현 못 한다.
// editor role이 있는 두 사용자가 서로의 페이지를 편집할 수 있다 (3편 ABAC에서 owner 속성으로 해결).
export default function PageDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [page, setPage] = useState<Page | null>(null);
  const [editing, setEditing] = useState(false);
  const [draftTitle, setDraftTitle] = useState('');
  const [draftContent, setDraftContent] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    apiClient
      .get<Page>(`/api/pages/${id}`)
      .then((r) => {
        setPage(r.data);
        setDraftTitle(r.data.title);
        setDraftContent(r.data.content);
      })
      .catch((e) => setError(e.response?.data?.message ?? 'failed'));
  }, [id]);

  async function save() {
    try {
      const r = await apiClient.put<Page>(`/api/pages/${id}`, {
        title: draftTitle,
        content: draftContent,
      });
      setPage(r.data);
      setEditing(false);
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  async function remove() {
    if (!confirm('정말 삭제할까요?')) return;
    try {
      await apiClient.delete(`/api/pages/${id}`);
      navigate('/pages');
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  if (error) return <p className="text-red-600">{error}</p>;
  if (!page) return <p>로딩...</p>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        {editing ? (
          <input
            value={draftTitle}
            onChange={(e) => setDraftTitle(e.target.value)}
            className="flex-1 rounded border p-2 text-xl font-semibold"
          />
        ) : (
          <h2 className="text-xl font-semibold">{page.title}</h2>
        )}
        <div className="flex gap-2">
          {!editing && (
            <PermissionGate permission="pages:edit">
              <button
                onClick={() => setEditing(true)}
                className="rounded bg-blue-600 px-3 py-1 text-sm text-white"
              >
                편집
              </button>
            </PermissionGate>
          )}
          {!editing && (
            <PermissionGate permission="pages:delete">
              <button
                onClick={remove}
                className="rounded bg-red-600 px-3 py-1 text-sm text-white"
              >
                삭제
              </button>
            </PermissionGate>
          )}
          {editing && (
            <>
              <button
                onClick={save}
                className="rounded bg-green-600 px-3 py-1 text-sm text-white"
              >
                저장
              </button>
              <button
                onClick={() => setEditing(false)}
                className="rounded border px-3 py-1 text-sm"
              >
                취소
              </button>
            </>
          )}
        </div>
      </div>
      <p className="text-xs text-slate-400">owner: user #{page.owner_id}</p>
      {editing ? (
        <textarea
          value={draftContent}
          onChange={(e) => setDraftContent(e.target.value)}
          className="h-48 w-full rounded border p-2 font-mono text-sm"
        />
      ) : (
        <pre className="whitespace-pre-wrap rounded bg-white p-4 shadow">{page.content}</pre>
      )}
    </div>
  );
}
