import { useAuth } from '@/features/auth/AuthContext';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { User, Mail, Shield, Calendar } from 'lucide-react';

export function CustomerAccountPage() {
  const { user } = useAuth();

  return (
    <div className="space-y-6 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">Account Details</h1>
        <p className="text-sm text-slate-500 mt-1">
          Review your customer identity and authentication profile.
        </p>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-full bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              <User className="h-5 w-5" />
            </div>
            <div>
              <CardTitle className="text-base">{user?.email}</CardTitle>
              <CardDescription>Customer Account Reference</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="p-3 bg-slate-50 rounded-md border border-slate-200">
              <div className="flex items-center gap-2 text-xs text-slate-500 font-semibold uppercase mb-1">
                <Mail className="h-3.5 w-3.5" /> Email Address
              </div>
              <div className="text-sm font-medium text-slate-800">{user?.email}</div>
            </div>

            <div className="p-3 bg-slate-50 rounded-md border border-slate-200">
              <div className="flex items-center gap-2 text-xs text-slate-500 font-semibold uppercase mb-1">
                <Shield className="h-3.5 w-3.5" /> Assigned Role
              </div>
              <Badge variant="indigo">{user?.role}</Badge>
            </div>

            <div className="p-3 bg-slate-50 rounded-md border border-slate-200">
              <div className="flex items-center gap-2 text-xs text-slate-500 font-semibold uppercase mb-1">
                <Calendar className="h-3.5 w-3.5" /> Member Since
              </div>
              <div className="text-sm font-medium text-slate-800">
                {user?.createdAt ? new Date(user.createdAt).toLocaleDateString() : 'N/A'}
              </div>
            </div>

            <div className="p-3 bg-slate-50 rounded-md border border-slate-200">
              <div className="flex items-center gap-2 text-xs text-slate-500 font-semibold uppercase mb-1">
                Account ID
              </div>
              <div className="text-xs font-mono text-slate-700 truncate">{user?.id}</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
