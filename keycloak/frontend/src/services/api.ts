import axios from 'axios';
import authService from './authService';

const api = axios.create({
  baseURL: 'http://localhost:8081/api'
});

// Request interceptor to add token to headers
api.interceptors.request.use(async (config) => {
  const token = authService.getAccessToken();
  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor to handle token expiration
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, try to refresh
      const refreshed = await authService.refreshAccessToken();
      if (refreshed) {
        // Retry the original request with new token
        const originalRequest = error.config;
        const newToken = authService.getAccessToken();
        if (newToken) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
          return api.request(originalRequest);
        }
      }
      
      // If refresh failed, redirect to login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export interface UserInfo {
  name: string;
  email: string;
}

export const getUserInfo = (): Promise<{ data: UserInfo }> => {
  return api.get('/user');
};

export default api;
