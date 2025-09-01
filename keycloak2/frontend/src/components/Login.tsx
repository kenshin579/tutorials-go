import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import keycloak from '../services/keycloak';

const Login: React.FC = () => {
  const navigate = useNavigate();

  const handleLogin = () => {
    keycloak.login({
      redirectUri: 'http://localhost:3000/'
    });
  };

  useEffect(() => {
    if (keycloak.authenticated === true) {
      navigate('/profile');
    }
  }, [navigate]);

  if (keycloak.authenticated === true) {
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
      <button 
        onClick={handleLogin}
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
        Login with Keycloak
      </button>
    </div>
  );
};

export default Login;
