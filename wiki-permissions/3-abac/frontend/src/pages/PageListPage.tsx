import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { apiClient } from '../api/client';

interface Department {
  id: number;
  name: string;
}

interface Page {
  id: number;
  title: string;
  owner_id: number;
  confidentiality: 'public' | 'internal' | 'confidential';
  department: Department | null;
}

// confidentiality에 따라 색상 뱃지를 반환한다.
function confidentialityBadgeClass(c: Page['confidentiality']): string {
  switch (c) {
    case 'public':
      return 'bg-green-100 text-green-800';
    case 'internal':
      return 'bg-orange-100 text-orange-800';
    case 'confidential':
      return 'bg-red-100 text-red-800';
  }
}

// PageListPage는 ABAC read 정책을 통과한 페이지(서버에서 이미 필터링됨)만 표시한다.
// 각 페이지의 confidentiality와 department를 뱃지로 표시해 ABAC 속성을 직관적으로 보여준다.
export default function PageListPage() {
  const [pages, setPages] = useState<Page[]>([]);
  const [error, setError] = useState('');

  useEffect(() => {
    apiClient
      .get<Page[]>('/api/pages')
      .then((r) => setPages(r.data))
      .catch((e) => setError(e.response?.data?.message ?? 'failed'));
  }, []);

  return (
    <div>
      <h2 className="mb-4 text-lg font-semibold">접근 가능한 페이지</h2>
      {error && <p className="mb-2 text-red-600">{error}</p>}
      <ul className="divide-y rounded bg-white shadow">
        {pages.map((p) => (
          <li key={p.id} className="flex items-center justify-between p-4 hover:bg-slate-50">
            <Link to={`/pages/${p.id}`} className="text-blue-600">
              {p.title}
            </Link>
            <div className="flex items-center gap-2">
              <span className={`rounded px-2 py-0.5 text-xs ${confidentialityBadgeClass(p.confidentiality)}`}>
                {p.confidentiality}
              </span>
              {p.department && (
                <span className="rounded bg-slate-100 px-2 py-0.5 text-xs">
                  {p.department.name}
                </span>
              )}
            </div>
          </li>
        ))}
        {pages.length === 0 && !error && (
          <li className="p-4 text-slate-500">접근 가능한 페이지가 없습니다.</li>
        )}
      </ul>
      <p className="mt-3 text-xs text-slate-500">
        ABAC read 정책을 통과한 페이지만 표시됩니다 (서버에서 필터링).
      </p>
    </div>
  );
}
