import { apiClient } from './client';
import { AuthResponse, User } from '@/types';

export interface RegisterPayload {
  email: string;
  password: string;
  companyName?: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export async function loginApi(payload: LoginPayload): Promise<AuthResponse> {
  return apiClient<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function registerApi(payload: RegisterPayload): Promise<AuthResponse> {
  return apiClient<AuthResponse>('/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function logoutApi(): Promise<{ message: string }> {
  return apiClient<{ message: string }>('/auth/logout', {
    method: 'POST',
  });
}

export async function getMeApi(): Promise<{ user: User }> {
  return apiClient<{ user: User }>('/auth/me', {
    method: 'GET',
  });
}

export async function adminPingApi(): Promise<{ message: string; timestamp: string }> {
  return apiClient<{ message: string; timestamp: string }>('/admin/ping', {
    method: 'GET',
  });
}
