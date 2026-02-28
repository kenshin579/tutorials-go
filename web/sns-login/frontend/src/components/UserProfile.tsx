import type { User } from '../types/auth';

interface Props {
  user: User;
  onLogout: () => void;
}

export function UserProfile({ user, onLogout }: Props) {
  return (
    <div style={{ textAlign: 'center', marginTop: '100px' }}>
      <h1>환영합니다!</h1>
      <div style={{ margin: '20px 0' }}>
        {user.avatar_url && (
          <img
            src={user.avatar_url}
            alt="프로필"
            width="80"
            height="80"
            style={{ borderRadius: '50%' }}
          />
        )}
        <h2>{user.name}</h2>
        <p>{user.email}</p>
        <p style={{ color: '#666' }}>Provider: {user.provider}</p>
      </div>
      <button
        onClick={onLogout}
        style={{
          padding: '10px 20px',
          fontSize: '14px',
          border: '1px solid #dadce0',
          borderRadius: '4px',
          backgroundColor: '#fff',
          cursor: 'pointer',
        }}
      >
        로그아웃
      </button>
    </div>
  );
}
