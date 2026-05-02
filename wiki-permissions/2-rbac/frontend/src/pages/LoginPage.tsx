import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState('alice@example.com');
  const [password, setPassword] = useState('password');
  const [error, setError] = useState('');

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    try {
      await login(email, password);
      navigate('/pages');
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } }).response?.data?.message ??
        '로그인 실패';
      setError(msg);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50">
      <form onSubmit={onSubmit} className="w-full max-w-sm rounded-lg bg-white p-6 shadow">
        <h1 className="mb-4 text-xl font-bold">wiki-permissions / 2-rbac</h1>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="mb-2 w-full rounded border p-2"
          placeholder="email"
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="mb-2 w-full rounded border p-2"
          placeholder="password"
        />
        {error && <p className="mb-2 text-sm text-red-600">{error}</p>}
        <button className="w-full rounded bg-blue-600 p-2 text-white">로그인</button>
        <p className="mt-3 text-xs text-slate-500">
          시드 계정: alice (admin) / bob (editor) / carol·dave (viewer) — 비밀번호 모두 password
        </p>
      </form>
    </div>
  );
}
