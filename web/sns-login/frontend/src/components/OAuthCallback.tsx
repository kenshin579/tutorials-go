import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { handleCallback } from '../services/authService';
import type { User } from '../types/auth';

interface Props {
  onLogin: (user: User, accessToken: string, refreshToken: string) => void;
}

export function OAuthCallback({ onLogin }: Props) {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state');

    if (!code || !state) {
      setError('인증 정보가 없습니다');
      return;
    }

    handleCallback(code, state)
      .then((res) => {
        onLogin(res.user, res.tokens.access_token, res.tokens.refresh_token);
        navigate('/');
      })
      .catch(() => {
        setError('로그인에 실패했습니다');
      });
  }, [searchParams, onLogin, navigate]);

  if (error) {
    return (
      <div style={{ textAlign: 'center', marginTop: '100px' }}>
        <h2>오류</h2>
        <p>{error}</p>
        <a href="/login">로그인 페이지로 돌아가기</a>
      </div>
    );
  }

  return (
    <div style={{ textAlign: 'center', marginTop: '100px' }}>
      <p>로그인 처리 중...</p>
    </div>
  );
}
