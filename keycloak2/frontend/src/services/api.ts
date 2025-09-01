import axios from 'axios';
import keycloak from './keycloak';

const api = axios.create({
  baseURL: 'http://localhost:8081/api'
});

// Request interceptor to add token to headers
api.interceptors.request.use((config) => {
  if (keycloak.token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${keycloak.token}`;
  }
  return config;
});

// Response interceptor to handle token expiration
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, redirect to login
      keycloak.login();
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
