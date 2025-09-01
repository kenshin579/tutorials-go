import React, { useEffect, useState } from 'react'
import { getUserInfo, validateToken } from '../services/api'

export default function UserInfo() {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [user, setUser] = useState(null)
  const [tokenValid, setTokenValid] = useState(null)

  useEffect(() => {
    let mounted = true
    ;(async () => {
      try {
        const [u, v] = await Promise.all([getUserInfo(), validateToken()])
        if (!mounted) return
        setUser(u)
        setTokenValid(v?.valid === true)
      } catch (e) {
        setError(e?.response?.data || e.message)
      } finally {
        setLoading(false)
      }
    })()
    return () => {
      mounted = false
    }
  }, [])

  if (loading) return <div className="card">불러오는 중...</div>
  if (error) return <div className="card error">에러: {typeof error === 'string' ? error : JSON.stringify(error)}</div>

  return (
    <div className="card">
      <h3>사용자 정보</h3>
      <pre style={{ whiteSpace: 'pre-wrap' }}>{JSON.stringify(user, null, 2)}</pre>
      <div>토큰 유효성: {tokenValid ? '유효' : '무효'}</div>
    </div>
  )
}
