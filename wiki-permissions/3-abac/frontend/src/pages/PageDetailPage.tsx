import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { apiClient } from '../api/client';

interface Department {
  id: number;
  name: string;
}

interface Page {
  id: number;
  title: string;
  content: string;
  owner_id: number;
  confidentiality: 'public' | 'internal' | 'confidential';
  department: Department | null;
}

interface Decision {
  allowed: boolean;
  reason: string;
  policy: string;
}

interface PageWithDecision {
  page: Page;
  can_read: Decision;
  can_edit: Decision;
}

// DecisionCard는 ABAC 결정의 reason과 policy 식별자를 사용자에게 표시한다.
// "왜 허용/거부됐는지"를 그대로 보여주는 것이 ABAC의 미덕이다 (1·2편 단순 yes/no와 대비).
function DecisionCard({ label, decision }: { label: string; decision: Decision }) {
  const bg = decision.allowed ? 'bg-green-50 border-green-300' : 'bg-red-50 border-red-300';
  const icon = decision.allowed ? '✓' : '✗';
  return (
    <div className={`rounded border-l-4 p-3 text-sm ${bg}`}>
      <div className="flex items-center gap-2">
        <span className="font-semibold">{icon} {label}</span>
        <span className="rounded bg-white px-2 py-0.5 text-xs text-slate-600">
          policy: {decision.policy}
        </span>
      </div>
      <p className="mt-1 text-slate-700">{decision.reason}</p>
    </div>
  );
}

export default function PageDetailPage() {
  const { id } = useParams();
  const [data, setData] = useState<PageWithDecision | null>(null);
  const [editing, setEditing] = useState(false);
  const [draftTitle, setDraftTitle] = useState('');
  const [draftContent, setDraftContent] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    apiClient
      .get<PageWithDecision>(`/api/pages/${id}`)
      .then((r) => {
        setData(r.data);
        setDraftTitle(r.data.page.title);
        setDraftContent(r.data.page.content);
      })
      .catch((e) => setError(e.response?.data?.message ?? 'failed'));
  }, [id]);

  async function save() {
    try {
      const r = await apiClient.put<Page>(`/api/pages/${id}`, {
        title: draftTitle,
        content: draftContent,
      });
      setData((prev) => prev && { ...prev, page: r.data });
      setEditing(false);
    } catch (e: unknown) {
      const msg =
        (e as { response?: { data?: { message?: string } } }).response?.data?.message ?? 'failed';
      setError(msg);
    }
  }

  if (error) return <p className="text-red-600">{error}</p>;
  if (!data) return <p>로딩...</p>;

  const { page, can_read, can_edit } = data;

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
          {!editing && can_edit.allowed && (
            <button
              onClick={() => setEditing(true)}
              className="rounded bg-blue-600 px-3 py-1 text-sm text-white"
            >
              편집
            </button>
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

      <div className="text-xs text-slate-500">
        분류: {page.confidentiality}
        {page.department && ` · 부서: ${page.department.name}`} · owner: user #{page.owner_id}
      </div>

      <div className="space-y-2">
        <DecisionCard label="읽기 권한" decision={can_read} />
        <DecisionCard label="편집 권한" decision={can_edit} />
      </div>

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
