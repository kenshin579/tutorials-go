import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import keycloak from './services/keycloak';
import Login from './components/Login';
import UserProfile from './components/UserProfile';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  const [keycloakInitialized, setKeycloakInitialized] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    // Check if already initialized
    if (keycloak.authenticated !== undefined) {
      setKeycloakInitialized(true);
      setIsAuthenticated(keycloak.authenticated === true);
      return;
    }

    keycloak.init({ 
      onLoad: 'check-sso',
      checkLoginIframe: false,
      silentCheckSsoRedirectUri: window.location.origin + '/silent-check-sso.html',
      enableLogging: true
    })
      .then((authenticated: boolean) => {
        console.log('Keycloak initialized. Authenticated:', authenticated);
        console.log('Current URL:', window.location.href);
        setKeycloakInitialized(true);
        setIsAuthenticated(authenticated);
      })
      .catch((error: any) => {
        console.error('Keycloak initialization failed:', error);
        console.log('Continuing without authentication for demo purposes');
        setKeycloakInitialized(true); // Continue anyway for demo
        setIsAuthenticated(false);
      });
  }, []);

  // Listen for authentication changes
  useEffect(() => {
    const interval = setInterval(() => {
      if (keycloak.authenticated !== isAuthenticated) {
        console.log('Authentication status changed:', keycloak.authenticated);
        setIsAuthenticated(keycloak.authenticated === true);
      }
    }, 500);

    return () => clearInterval(interval);
  }, [isAuthenticated]);

  if (!keycloakInitialized) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontFamily: 'Arial, sans-serif'
      }}>
        <div>
          <h2>Loading...</h2>
          <p>Initializing authentication...</p>
        </div>
      </div>
    );
  }

  // Remove error handling for demo purposes

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route 
            path="/profile" 
            element={
              <ProtectedRoute>
                <UserProfile />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/" 
            element={
              isAuthenticated ? 
                <Navigate to="/profile" replace /> : 
                <Navigate to="/login" replace />
            } 
          />
        </Routes>
      </div>
    </Router>
  );
}

export default App;