import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { apiClient } from '../api/client';
import { useAuth } from '../auth/AuthContext';
import ShareModal from '../components/ShareModal';

interface Page {
  id: number;
  title: string;
  content: string;
  owner_id: number;
}

// PageDetailPage는 페이지 상세 조회/편집 화면이다.
// 편집 버튼은 모두에게 보이지만, 권한이 없으면 PUT 요청이 서버에서 403으로 거부된다.
// (UX 차원에서 사전 게이팅은 본 1편의 ACL 모델 한계로 글에서 별도 다룸.)
export default function PageDetailPage() {
  const { id } = useParams();
  const { user } = useAuth();
  const [page, setPage] = useState<Page | null>(null);
  const [editing, setEditing] = useState(false);
  const [draftTitle, setDraftTitle] = useState('');
  const [draftContent, setDraftContent] = useState('');
  const [shareOpen, setShareOpen] = useState(false);
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

  if (error) return <p className="text-red-600">{error}</p>;
  if (!page) return <p>로딩...</p>;

  const isOwner = user?.id === page.owner_id;

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
          {isOwner && (
            <button
              onClick={() => setShareOpen(true)}
              className="rounded border px-3 py-1 text-sm"
            >
              공유 관리
            </button>
          )}
          {!editing && (
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
      {editing ? (
        <textarea
          value={draftContent}
          onChange={(e) => setDraftContent(e.target.value)}
          className="h-48 w-full rounded border p-2 font-mono text-sm"
        />
      ) : (
        <pre className="whitespace-pre-wrap rounded bg-white p-4 shadow">{page.content}</pre>
      )}
      {isOwner && shareOpen && (
        <ShareModal pageId={page.id} onClose={() => setShareOpen(false)} />
      )}
    </div>
  );
}
