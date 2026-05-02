import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './auth/AuthContext';
import ProtectedRoute from './auth/ProtectedRoute';
import Layout from './components/Layout';
import LoginPage from './pages/LoginPage';
import PageListPage from './pages/PageListPage';
import PageDetailPage from './pages/PageDetailPage';
import UsersPage from './pages/UsersPage';

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            element={
              <ProtectedRoute>
                <Layout />
              </ProtectedRoute>
            }
          >
            <Route path="/pages" element={<PageListPage />} />
            <Route path="/pages/:id" element={<PageDetailPage />} />
            <Route path="/users" element={<UsersPage />} />
            <Route index element={<Navigate to="/pages" replace />} />
          </Route>
          <Route path="*" element={<Navigate to="/pages" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
