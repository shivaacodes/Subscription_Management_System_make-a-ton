export type Role = 'ADMIN' | 'STAFF' | 'CUSTOMER';

export interface User {
  id: string;
  email: string;
  role: Role;
  createdAt: string;
  updatedAt: string;
}

export interface CustomerProfile {
  id: string;
  userId: string;
  companyName?: string;
  phone?: string;
  billingAddress?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ApiErrorResponse {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
}

export interface AuthResponse {
  user: User;
  message?: string;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
  database: string;
}
