import { ApiErrorResponse } from '@/types';

export class ApiError extends Error {
  code: string;
  status: number;
  details?: Record<string, string>;

  constructor(status: number, errorData: ApiErrorResponse['error']) {
    super(errorData.message);
    this.name = 'ApiError';
    this.code = errorData.code;
    this.status = status;
    this.details = errorData.details;
  }
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1';

export async function apiClient<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${BASE_URL}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`;

  const headers = new Headers(options.headers || {});
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }
  headers.set('X-Requested-With', 'XMLHttpRequest');

  const response = await fetch(url, {
    ...options,
    headers,
    credentials: 'include', // essential for HttpOnly session cookie transmission
  });

  if (response.status === 204) {
    return {} as T;
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    const errorEnvelope: ApiErrorResponse['error'] = data?.error || {
      code: 'HTTP_ERROR',
      message: response.statusText || 'An unexpected server error occurred',
    };
    throw new ApiError(response.status, errorEnvelope);
  }

  return data as T;
}
