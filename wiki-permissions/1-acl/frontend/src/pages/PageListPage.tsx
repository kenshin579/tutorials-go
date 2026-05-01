import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { apiClient } from '../api/client';

interface Page {
  id: number;
  title: string;
  owner_id: number;
}

// PageListPage는 본인이 access 가능한 페이지(owner이거나 ACL 받은) 목록을 보여준다.
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
      <h2 className="mb-4 text-lg font-semibold">내가 접근 가능한 페이지</h2>
      {error && <p className="text-red-600">{error}</p>}
      <ul className="divide-y rounded bg-white shadow">
        {pages.map((p) => (
          <li key={p.id} className="p-4 hover:bg-slate-50">
            <Link to={`/pages/${p.id}`} className="text-blue-600">
              {p.title}
            </Link>
          </li>
        ))}
        {pages.length === 0 && !error && (
          <li className="p-4 text-slate-500">접근 가능한 페이지가 없습니다.</li>
        )}
      </ul>
    </div>
  );
}
