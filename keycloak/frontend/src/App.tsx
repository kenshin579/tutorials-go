import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import authService from './services/authService';
import Login from './components/Login';
import UserProfile from './components/UserProfile';
import ProtectedRoute from './components/ProtectedRoute';
import Callback from './components/Callback';

function App() {
  const [authInitialized, setAuthInitialized] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    // Initialize authentication state
    const initializeAuth = () => {
      try {
        const authenticated = authService.isAuthenticated();
        console.log('Auth initialized. Authenticated:', authenticated);
        setIsAuthenticated(authenticated);
        setAuthInitialized(true);
      } catch (error) {
        console.error('Auth initialization failed:', error);
        setAuthInitialized(true);
        setIsAuthenticated(false);
      }
    };

    initializeAuth();
  }, []);

  // Listen for authentication changes
  useEffect(() => {
    const interval = setInterval(() => {
      const currentAuthState = authService.isAuthenticated();
      if (currentAuthState !== isAuthenticated) {
        console.log('Authentication status changed:', currentAuthState);
        setIsAuthenticated(currentAuthState);
      }
    }, 1000); // Check every second

    return () => clearInterval(interval);
  }, [isAuthenticated]);

  if (!authInitialized) {
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
          <Route path="/callback" element={<Callback />} />
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