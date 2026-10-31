import { useEffect, useState } from 'react';
import { getSubscriptions, pauseSubscription, resumeSubscription, cancelSubscription } from '@/lib/api/modules';
import { Subscription } from '@/types';
import { formatCurrency } from '@/lib/utils';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Alert } from '@/components/ui/alert';
import { Repeat, Play, Pause, XCircle, ShieldCheck, RefreshCw } from 'lucide-react';

export function CustomerSubscriptionsPage() {
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [message, setMessage] = useState<string | null>(null);

  const loadSubscriptions = async () => {
    setIsLoading(true);
    try {
      const res = await getSubscriptions();
      setSubscriptions(res.subscriptions || []);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadSubscriptions();
  }, []);

  const handlePause = async (id: string) => {
    try {
      await pauseSubscription(id);
      setMessage('Subscription paused successfully');
      loadSubscriptions();
    } catch (err: any) {
      alert(err.message || 'Cannot pause subscription according to plan policy');
    }
  };

  const handleResume = async (id: string) => {
    try {
      await resumeSubscription(id);
      setMessage('Subscription resumed successfully');
      loadSubscriptions();
    } catch (err: any) {
      alert(err.message || 'Failed resuming subscription');
    }
  };

  const handleCancel = async (id: string) => {
    if (!confirm('Are you sure you wish to cancel this subscription?')) return;
    try {
      await cancelSubscription(id);
      setMessage('Subscription cancelled');
      loadSubscriptions();
    } catch (err: any) {
      alert(err.message || 'Cannot cancel subscription according to plan policy');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">My Subscriptions</h1>
          <p className="text-sm text-slate-500 mt-1">
            Manage your active Brightloop Software subscriptions, renewal dates, and permitted add-on services.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={loadSubscriptions} isLoading={isLoading}>
          <RefreshCw className="h-4 w-4 mr-1 text-slate-500" /> Refresh
        </Button>
      </div>

      {message && (
        <Alert variant="success" className="cursor-pointer" onClick={() => setMessage(null)}>
          {message}
        </Alert>
      )}

      {subscriptions.length > 0 ? (
        <div className="grid grid-cols-1 gap-6">
          {subscriptions.map((sub) => (
            <Card key={sub.id} className="border-slate-200">
              <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b bg-slate-50/50 pb-4">
                <div className="flex items-center gap-3">
                  <div className="h-10 w-10 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
                    <Repeat className="h-5 w-5" />
                  </div>
                  <div>
                    <CardTitle className="text-base text-slate-900">{sub.planName || 'Software Subscription'}</CardTitle>
                    <span className="text-xs font-mono text-slate-500">{sub.subscriptionNumber}</span>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <Badge
                    variant={
                      sub.status === 'ACTIVE'
                        ? 'success'
                        : sub.status === 'PAUSED'
                        ? 'warning'
                        : 'default'
                    }
                  >
                    {sub.status}
                  </Badge>
                  <div className="text-right">
                    <span className="text-xs text-slate-400 block uppercase font-semibold">Total Amount</span>
                    <span className="text-base font-bold text-slate-900">{formatCurrency(sub.totalAmountCents || 0)}</span>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="p-6 space-y-6">
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
                  <div className="p-3 bg-slate-50 rounded border">
                    <span className="text-slate-500 font-semibold uppercase block mb-1">Contract Start Date</span>
                    <div className="font-medium text-slate-900">{new Date(sub.startDate).toLocaleDateString()}</div>
                  </div>
                  <div className="p-3 bg-slate-50 rounded border">
                    <span className="text-slate-500 font-semibold uppercase block mb-1">Next Renewal Date</span>
                    <div className="font-medium text-slate-900">
                      {sub.nextBillingDate ? new Date(sub.nextBillingDate).toLocaleDateString() : '—'}
                    </div>
                  </div>
                  <div className="p-3 bg-slate-50 rounded border">
                    <span className="text-slate-500 font-semibold uppercase block mb-1">Payment Terms</span>
                    <div className="font-medium text-slate-900">{sub.paymentTerms}</div>
                  </div>
                </div>

                {/* Permitted Self-service actions */}
                <div className="border-t pt-4 flex items-center justify-between">
                  <div className="text-xs text-slate-500 flex items-center gap-1.5">
                    <ShieldCheck className="h-4 w-4 text-indigo-600" />
                    <span>Plan policies govern pause and cancellation permissions</span>
                  </div>
                  <div className="flex items-center gap-2">
                    {sub.status === 'ACTIVE' && (
                      <Button variant="secondary" size="sm" onClick={() => handlePause(sub.id)} className="text-xs">
                        <Pause className="h-3 w-3 mr-1" /> Pause Subscription
                      </Button>
                    )}
                    {sub.status === 'PAUSED' && (
                      <Button variant="primary" size="sm" onClick={() => handleResume(sub.id)} className="text-xs">
                        <Play className="h-3 w-3 mr-1" /> Resume Subscription
                      </Button>
                    )}
                    {sub.status !== 'CANCELLED' && sub.status !== 'CLOSED' && (
                      <Button variant="outline" size="sm" onClick={() => handleCancel(sub.id)} className="text-xs text-rose-600">
                        <XCircle className="h-3 w-3 mr-1" /> Cancel
                      </Button>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <Card className="border-dashed border-2 p-8 text-center bg-white">
          <Repeat className="h-10 w-10 text-slate-300 mx-auto mb-2" />
          <h3 className="text-sm font-semibold text-slate-800">No Active Subscriptions</h3>
          <p className="text-xs text-slate-500 mt-1 max-w-sm mx-auto">
            You do not currently have any active subscriptions. When a Brightloop Software proposal is confirmed, it will appear here.
          </p>
        </Card>
      )}
    </div>
  );
}
