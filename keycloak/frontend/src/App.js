import React, { useEffect, useState } from 'react'
import { Routes, Route, Link, Navigate } from 'react-router-dom'
import { initKeycloak, getKeycloak } from './services/keycloak'
import Login from './components/Login'
import Logout from './components/Logout'
import UserInfo from './components/UserInfo'

function ProtectedRoute({ children }) {
  const kc = getKeycloak()
  if (!kc || !kc.authenticated) {
    return <Navigate to="/" replace />
  }
  return children
}

export default function App() {
  const [ready, setReady] = useState(false)
  const [authenticated, setAuthenticated] = useState(false)

  useEffect(() => {
    ;(async () => {
      const ok = await initKeycloak()
      setAuthenticated(!!getKeycloak()?.authenticated)
      setReady(true)
    })()
  }, [])

  if (!ready) return <div className="container">초기화 중...</div>

  return (
    <div className="container">
      <header className="header">
        <h2>Keycloak React 샘플</h2>
        <nav>
          <Link to="/">홈</Link> | <Link to="/protected">보호된 페이지</Link>
        </nav>
        <div className="auth">
          {authenticated ? <Logout /> : <Login />}
        </div>
      </header>

      <main>
        <Routes>
          <Route
            path="/"
            element={
              <div className="card">
                <h3>환영합니다</h3>
                <p>로그인 후 보호된 사용자 정보를 확인할 수 있습니다.</p>
              </div>
            }
          />
          <Route
            path="/protected"
            element={
              <ProtectedRoute>
                <UserInfo />
              </ProtectedRoute>
            }
          />
        </Routes>
      </main>
    </div>
  )
}
