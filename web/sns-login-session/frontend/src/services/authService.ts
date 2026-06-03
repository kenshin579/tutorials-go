import axios from 'axios';
import type { User } from '../types/auth';

const API_URL = import.meta.env.VITE_API_URL || '';

// withCredentials: 쿠키를 교차 출처 요청에 포함
const api = axios.create({ baseURL: API_URL, withCredentials: true });

export async function getGoogleAuthURL(): Promise<string> {
  const { data } = await api.get<{ url: string }>('/api/auth/google/url');
  return data.url;
}

export async function getMe(): Promise<User> {
  const { data } = await api.get<User>('/api/user/me');
  return data;
}

export async function logout(): Promise<void> {
  await api.post('/api/auth/logout');
}
