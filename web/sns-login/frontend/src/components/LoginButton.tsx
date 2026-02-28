import { getGoogleAuthURL } from '../services/authService';

export function LoginButton() {
  const handleGoogleLogin = async () => {
    const url = await getGoogleAuthURL();
    window.location.href = url;
  };

  return (
    <div style={{ textAlign: 'center', marginTop: '100px' }}>
      <h1>SNS 로그인 데모</h1>
      <p>Google 계정으로 로그인하세요</p>
      <button
        onClick={handleGoogleLogin}
        style={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: '8px',
          padding: '12px 24px',
          fontSize: '16px',
          border: '1px solid #dadce0',
          borderRadius: '4px',
          backgroundColor: '#fff',
          cursor: 'pointer',
        }}
      >
        <img
          src="https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg"
          alt="Google"
          width="20"
          height="20"
        />
        Google로 로그인
      </button>
    </div>
  );
}
