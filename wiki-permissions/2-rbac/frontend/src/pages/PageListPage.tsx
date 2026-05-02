import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { apiClient } from '../api/client';
import PermissionGate from '../components/PermissionGate';

interface Page {
  id: number;
  title: string;
  owner_id: number;
}

// PageListPage는 사용자가 pages:read 권한을 가졌을 때 모든 페이지를 보여준다.
// (1편 ACL은 본인 owner / ACL 매칭 페이지만 노출했지만, RBAC은 role 기반이라 모두 노출.)
// pages:create 권한이 있으면 "새 페이지" 버튼이 표시된다.
export default function PageListPage() {
  const [pages, setPages] = useState<Page[]>([]);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    apiClient
      .get<Page[]>('/api/pages')
      .then((r) => setPages(r.data))
      .catch((e) => setError(e.response?.data?.message ?? 'failed'));
  }, []);

  async function createPage() {
    const title = prompt('새 페이지 제목');
    if (!title) return;
    try {
      const r = await apiClient.post<Page>('/api/pages', { title, content: '' });
      navigate(`/pages/${r.data.id}`);
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  return (
    <div>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">페이지 목록</h2>
        <PermissionGate permission="pages:create">
          <button onClick={createPage} className="rounded bg-blue-600 px-3 py-1 text-sm text-white">
            새 페이지
          </button>
        </PermissionGate>
      </div>
      {error && <p className="text-red-600">{error}</p>}
      <ul className="divide-y rounded bg-white shadow">
        {pages.map((p) => (
          <li key={p.id} className="p-4 hover:bg-slate-50">
            <Link to={`/pages/${p.id}`} className="text-blue-600">
              {p.title}
            </Link>
            <span className="ml-2 text-xs text-slate-400">owner #{p.owner_id}</span>
          </li>
        ))}
        {pages.length === 0 && !error && (
          <li className="p-4 text-slate-500">페이지가 없습니다.</li>
        )}
      </ul>
    </div>
  );
}
