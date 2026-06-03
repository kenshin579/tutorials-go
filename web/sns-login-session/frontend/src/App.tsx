import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import { LoginButton } from './components/LoginButton';
import { UserProfile } from './components/UserProfile';
import { OAuthCallback } from './components/OAuthCallback';
import { ProtectedRoute } from './components/ProtectedRoute';

function App() {
  const { user, loading, isAuthenticated, login, logout } = useAuth();

  if (loading) {
    return (
      <div style={{ textAlign: 'center', marginTop: '100px' }}>
        <p>로딩 중...</p>
      </div>
    );
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            isAuthenticated ? <Navigate to="/" replace /> : <LoginButton />
          }
        />
        <Route
          path="/auth/callback"
          element={<OAuthCallback onLogin={login} />}
        />
        <Route
          path="/"
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <UserProfile user={user!} onLogout={logout} />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
