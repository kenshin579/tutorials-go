import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../services/authService';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  const handleKeycloakLogin = () => {
    try {
      authService.initiateLogin();
    } catch (err) {
      console.error('Keycloak login error:', err);
      setError('Keycloak 로그인 중 오류가 발생했습니다.');
    }
  };

  useEffect(() => {
    // Check if already authenticated
    if (authService.isAuthenticated()) {
      console.log('User is already authenticated, redirecting to profile...');
      navigate('/profile', { replace: true });
    }
  }, [navigate]);

  if (authService.isAuthenticated()) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontFamily: 'Arial, sans-serif'
      }}>
        <div>Redirecting to profile...</div>
      </div>
    );
  }

  return (
    <div style={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center', 
      height: '100vh',
      fontFamily: 'Arial, sans-serif'
    }}>
      <h1>Keycloak Tutorial</h1>
      <p>Please login to continue</p>
      
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column', 
        alignItems: 'center',
        gap: '20px'
      }}>
        {/* Keycloak Authorization Code Flow 로그인 */}
        <div style={{ 
          padding: '30px',
          border: '2px solid #007bff',
          borderRadius: '12px',
          backgroundColor: '#f8f9fa',
          textAlign: 'center',
          maxWidth: '400px'
        }}>
          <h3 style={{ margin: '0 0 15px 0', color: '#007bff', fontSize: '20px' }}>
            Keycloak 로그인
          </h3>
          <p style={{ margin: '0 0 20px 0', fontSize: '14px', color: '#666', lineHeight: '1.5' }}>
            Authorization Code Flow를 사용하여<br/>
            Keycloak 호스팅 페이지에서 로그인합니다
          </p>
          <button 
            onClick={handleKeycloakLogin}
            style={{
              padding: '15px 30px',
              fontSize: '16px',
              backgroundColor: '#007bff',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              cursor: 'pointer',
              fontWeight: 'bold',
              transition: 'background-color 0.2s'
            }}
            onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#0056b3'}
            onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#007bff'}
          >
            🔐 Keycloak으로 로그인
          </button>
        </div>

        {error && (
          <div style={{ 
            color: '#dc3545', 
            fontSize: '14px',
            textAlign: 'center',
            padding: '12px',
            backgroundColor: '#f8d7da',
            border: '1px solid #f5c6cb',
            borderRadius: '8px',
            maxWidth: '400px'
          }}>
            {error}
          </div>
        )}
      </div>
    </div>
  );
};

export default Login;
