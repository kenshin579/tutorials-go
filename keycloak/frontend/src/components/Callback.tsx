import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import authService from '../services/authService';

const Callback: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<'processing' | 'success' | 'error'>('processing');
  const [error, setError] = useState<string | null>(null);
  const [processed, setProcessed] = useState(false);

  useEffect(() => {
    // 이미 처리되었거나 처리 중이면 중복 실행 방지
    if (processed) return;

    const handleCallback = async () => {
      try {
        setProcessed(true); // 처리 시작 표시
        
        const code = searchParams.get('code');
        const errorParam = searchParams.get('error');

        console.log('Callback processing:', { code: !!code, error: errorParam });

        if (errorParam) {
          setError(`인증 오류: ${errorParam}`);
          setStatus('error');
          return;
        }

        if (!code) {
          setError('인증 코드가 없습니다.');
          setStatus('error');
          return;
        }

        // 이미 인증된 상태라면 바로 프로필로 이동
        if (authService.isAuthenticated()) {
          console.log('Already authenticated, redirecting to profile');
          setStatus('success');
          navigate('/profile', { replace: true });
          return;
        }

        console.log('Attempting token exchange...');
        const success = await authService.handleCallback(code);
        
        if (success) {
          console.log('Token exchange successful');
          setStatus('success');
          // URL에서 코드 파라미터 제거
          window.history.replaceState({}, document.title, '/callback');
          setTimeout(() => {
            navigate('/profile', { replace: true });
          }, 1000);
        } else {
          console.error('Token exchange failed');
          setError('토큰 교환에 실패했습니다.');
          setStatus('error');
        }
      } catch (err) {
        console.error('Callback handling error:', err);
        setError(`콜백 처리 중 오류가 발생했습니다: ${err instanceof Error ? err.message : '알 수 없는 오류'}`);
        setStatus('error');
      }
    };

    handleCallback();
  }, [searchParams, navigate, processed]);

  if (status === 'processing') {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontFamily: 'Arial, sans-serif'
      }}>
        <div style={{ fontSize: '18px', marginBottom: '20px' }}>로그인 처리 중...</div>
        <div style={{ 
          width: '50px', 
          height: '50px', 
          border: '3px solid #f3f3f3',
          borderTop: '3px solid #007bff',
          borderRadius: '50%',
          animation: 'spin 1s linear infinite'
        }}></div>
        <style>{`
          @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
          }
        `}</style>
      </div>
    );
  }

  if (status === 'success') {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontFamily: 'Arial, sans-serif'
      }}>
        <div style={{ fontSize: '18px', color: '#28a745', marginBottom: '10px' }}>
          ✓ 로그인 성공!
        </div>
        <div style={{ fontSize: '14px', color: '#666' }}>
          프로필 페이지로 이동 중...
        </div>
      </div>
    );
  }

  return (
    <div style={{ 
      display: 'flex', 
      flexDirection: 'column',
      justifyContent: 'center', 
      alignItems: 'center', 
      height: '100vh',
      fontFamily: 'Arial, sans-serif'
    }}>
      <div style={{ fontSize: '18px', color: '#dc3545', marginBottom: '10px' }}>
        ✗ 로그인 실패
      </div>
      <div style={{ fontSize: '14px', color: '#666', marginBottom: '20px' }}>
        {error}
      </div>
      <button 
        onClick={() => navigate('/login')}
        style={{
          padding: '10px 20px',
          fontSize: '16px',
          backgroundColor: '#007bff',
          color: 'white',
          border: 'none',
          borderRadius: '5px',
          cursor: 'pointer'
        }}
      >
        다시 로그인
      </button>
    </div>
  );
};

export default Callback;
