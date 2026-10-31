import { useState } from 'react';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert } from '@/components/ui/alert';
import { adminPingApi } from '@/lib/api/auth';
import { apiClient } from '@/lib/api/client';
import { ShieldCheck, Activity, CheckCircle2, Server } from 'lucide-react';
import { HealthResponse } from '@/types';

export function SystemPage() {
  const [healthStatus, setHealthStatus] = useState<HealthResponse | null>(null);
  const [adminPingResult, setAdminPingResult] = useState<any | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loadingHealth, setLoadingHealth] = useState(false);
  const [loadingPing, setLoadingPing] = useState(false);

  const checkHealth = async () => {
    setLoadingHealth(true);
    setError(null);
    try {
      const res = await apiClient<HealthResponse>('/health');
      setHealthStatus(res);
    } catch (err: any) {
      setError(err.message || 'Health check failed');
    } finally {
      setLoadingHealth(false);
    }
  };

  const verifyAdminRole = async () => {
    setLoadingPing(true);
    setError(null);
    try {
      const res = await adminPingApi();
      setAdminPingResult(res);
    } catch (err: any) {
      setError(err.message || 'Admin ping authorization failed');
    } finally {
      setLoadingPing(false);
    }
  };

  return (
    <div className="space-y-6 max-w-4xl">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">System & Authorization Health</h1>
        <p className="text-sm text-slate-500 mt-1">
          Verify live backend API status, database connectivity, and role authorization enforcement.
        </p>
      </div>

      {error && <Alert variant="danger">{error}</Alert>}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Health Check Card */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2 text-indigo-600 mb-1">
              <Server className="h-5 w-5" />
              <CardTitle className="text-base">Backend & Database Status</CardTitle>
            </div>
            <CardDescription>Executes standard GET /api/v1/health probe.</CardDescription>
          </CardHeader>
          <CardContent>
            {healthStatus ? (
              <div className="bg-slate-50 p-4 rounded-md border border-slate-200 space-y-2 text-xs font-mono">
                <div className="flex justify-between">
                  <span className="text-slate-500">Service:</span>
                  <span className="text-emerald-700 font-semibold">{healthStatus.status}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">Database:</span>
                  <span className="text-slate-900 font-semibold">{healthStatus.database}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">Version:</span>
                  <span className="text-slate-900">{healthStatus.version}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-slate-500">Timestamp:</span>
                  <span className="text-slate-700">{healthStatus.timestamp}</span>
                </div>
              </div>
            ) : (
              <p className="text-xs text-slate-500 italic">Click below to query backend health.</p>
            )}
          </CardContent>
          <CardFooter>
            <Button
              variant="outline"
              size="sm"
              onClick={checkHealth}
              isLoading={loadingHealth}
              className="w-full"
            >
              <Activity className="h-4 w-4 mr-2 text-indigo-600" />
              Run Health Check
            </Button>
          </CardFooter>
        </Card>

        {/* Role Authorization Ping Card */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2 text-indigo-600 mb-1">
              <ShieldCheck className="h-5 w-5" />
              <CardTitle className="text-base">Role Authorization Ping</CardTitle>
            </div>
            <CardDescription>Executes GET /api/v1/admin/ping (ADMIN role only).</CardDescription>
          </CardHeader>
          <CardContent>
            {adminPingResult ? (
              <div className="bg-emerald-50/70 border border-emerald-200 p-4 rounded-md space-y-2 text-xs text-emerald-900 font-mono">
                <div className="flex items-center gap-2 font-semibold">
                  <CheckCircle2 className="h-4 w-4 text-emerald-600" />
                  <span>Authorized ADMIN Access</span>
                </div>
                <div>{adminPingResult.message}</div>
                <div className="text-[11px] text-emerald-700">{adminPingResult.timestamp}</div>
              </div>
            ) : (
              <p className="text-xs text-slate-500 italic">Click below to test role-protected endpoint.</p>
            )}
          </CardContent>
          <CardFooter>
            <Button
              variant="primary"
              size="sm"
              onClick={verifyAdminRole}
              isLoading={loadingPing}
              className="w-full"
            >
              <ShieldCheck className="h-4 w-4 mr-2" />
              Test ADMIN Ping
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
