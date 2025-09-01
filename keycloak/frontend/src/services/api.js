import axios from 'axios'
import { getToken, getKeycloak } from './keycloak'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081',
})

// 요청 인터셉터: 토큰 추가
api.interceptors.request.use(async (config) => {
  const token = getToken()
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 응답 인터셉터: 401 발생 시 리프레시 시도 후 재요청
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const originalRequest = error.config
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const kc = getKeycloak()
        if (kc) {
          await kc.updateToken(30)
          originalRequest.headers = originalRequest.headers || {}
          originalRequest.headers.Authorization = `Bearer ${kc.token}`
          return api(originalRequest)
        }
      } catch (e) {
        // 리프레시 실패 시 로그아웃 유도
        console.warn('Token refresh failed on 401')
      }
    }
    return Promise.reject(error)
  }
)

export function getUserInfo() {
  return api.get('/api/protected/user').then((r) => r.data)
}

export function validateToken() {
  return api.get('/api/validate').then((r) => r.data)
}

export default api
