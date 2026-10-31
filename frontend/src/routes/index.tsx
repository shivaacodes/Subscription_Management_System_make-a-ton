import { createBrowserRouter, Navigate } from 'react-router-dom';
import { PublicLayout } from '@/components/layout/PublicLayout';
import { AppLayout } from '@/components/layout/AppLayout';
import { CustomerLayout } from '@/components/layout/CustomerLayout';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';
import { useAuth } from '@/features/auth/AuthContext';

import { LoginPage } from '@/features/auth/LoginPage';
import { RegisterPage } from '@/features/auth/RegisterPage';
import { ForgotPasswordPage } from '@/features/auth/ForgotPasswordPage';

import { DashboardPage } from '@/features/dashboard/DashboardPage';
import { SubscriptionsPage } from '@/features/subscriptions/SubscriptionsPage';
import { QuotationsPage } from '@/features/quotations/QuotationsPage';
import { CatalogPage } from '@/features/catalog/CatalogPage';
import { PlansPage } from '@/features/plans/PlansPage';
import { BillingPage } from '@/features/billing/BillingPage';
import { PaymentsPage } from '@/features/payments/PaymentsPage';
import { ReportingPage } from '@/features/reporting/ReportingPage';
import { SystemPage } from '@/features/system/SystemPage';

import { CustomerSubscriptionsPage } from '@/features/portal/CustomerSubscriptionsPage';
import { CustomerInvoicesPage } from '@/features/portal/CustomerInvoicesPage';
import { CustomerAccountPage } from '@/features/portal/CustomerAccountPage';

function RootRedirect() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen bg-slate-50 flex items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-indigo-600 border-t-transparent" />
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (user.role === 'CUSTOMER') {
    return <Navigate to="/portal" replace />;
  }

  return <Navigate to="/admin" replace />;
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <RootRedirect />,
  },
  {
    element: <PublicLayout />,
    children: [
      { path: '/login', element: <LoginPage /> },
      { path: '/register', element: <RegisterPage /> },
      { path: '/forgot-password', element: <ForgotPasswordPage /> },
    ],
  },
  {
    path: '/admin',
    element: (
      <ProtectedRoute allowedRoles={['ADMIN', 'STAFF']}>
        <AppLayout />
      </ProtectedRoute>
    ),
    children: [
      { index: true, element: <DashboardPage /> },
      { path: 'subscriptions', element: <SubscriptionsPage /> },
      { path: 'quotations', element: <QuotationsPage /> },
      {
        path: 'catalog',
        element: (
          <ProtectedRoute allowedRoles={['ADMIN']}>
            <CatalogPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'plans',
        element: (
          <ProtectedRoute allowedRoles={['ADMIN']}>
            <PlansPage />
          </ProtectedRoute>
        ),
      },
      { path: 'billing', element: <BillingPage /> },
      { path: 'payments', element: <PaymentsPage /> },
      {
        path: 'reporting',
        element: (
          <ProtectedRoute allowedRoles={['ADMIN']}>
            <ReportingPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'system',
        element: (
          <ProtectedRoute allowedRoles={['ADMIN']}>
            <SystemPage />
          </ProtectedRoute>
        ),
      },
    ],
  },
  {
    path: '/portal',
    element: (
      <ProtectedRoute allowedRoles={['CUSTOMER', 'ADMIN']}>
        <CustomerLayout />
      </ProtectedRoute>
    ),
    children: [
      { index: true, element: <CustomerSubscriptionsPage /> },
      { path: 'invoices', element: <CustomerInvoicesPage /> },
      { path: 'account', element: <CustomerAccountPage /> },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
]);
