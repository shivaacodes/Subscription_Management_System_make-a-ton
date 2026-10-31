import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Repeat, ShieldAlert } from 'lucide-react';

export function CustomerSubscriptionsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">My Subscriptions</h1>
        <p className="text-sm text-slate-500 mt-1">
          Review active Brightloop Software subscriptions, renewal dates, and permitted add-ons.
        </p>
      </div>

      <Card className="border-dashed border-2 border-slate-200 bg-white">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center">
              <Repeat className="h-5 w-5" />
            </div>
            <div>
              <CardTitle className="text-base">No Active Subscriptions Found</CardTitle>
              <CardDescription>
                You do not have any active subscriptions associated with your account yet.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="p-4 rounded-md bg-slate-50 border border-slate-100 text-xs text-slate-600 flex items-start gap-2">
            <ShieldAlert className="h-4 w-4 text-indigo-500 mt-0.5 flex-shrink-0" />
            <div>
              Subscriptions are created when an approved Brightloop Software quotation is accepted or when an administrator provisions a plan for your organization.
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
