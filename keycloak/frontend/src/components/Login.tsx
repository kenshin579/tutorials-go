import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import keycloak from '../services/keycloak';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [isRedirecting, setIsRedirecting] = useState(false);

  const handleLogin = () => {
    keycloak.login({
      redirectUri: window.location.origin + '/profile'
    });
  };

  useEffect(() => {
    // Check authentication status periodically
    const checkAuth = () => {
      if (keycloak.authenticated === true) {
        console.log('User is authenticated, redirecting to profile...');
        setIsRedirecting(true);
        setTimeout(() => {
          navigate('/profile', { replace: true });
        }, 100);
      }
    };

    // Check immediately
    checkAuth();

    // Set up interval to check authentication status
    const interval = setInterval(checkAuth, 500);

    return () => clearInterval(interval);
  }, [navigate]);

  if (keycloak.authenticated === true || isRedirecting) {
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
