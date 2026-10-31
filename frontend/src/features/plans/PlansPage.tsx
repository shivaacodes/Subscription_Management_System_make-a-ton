import { ModulePlaceholder } from '@/components/ui/ModulePlaceholder';
import { CalendarDays } from 'lucide-react';

export function PlansPage() {
  return (
    <ModulePlaceholder
      title="Subscription Plans & Policies"
      description="Configure recurring billing frequencies, availability windows, pause permissions, and renewal policies."
      icon={CalendarDays}
      owner="Member 2 (Catalog & Configuration)"
      branch="feat/catalog-plans"
      deliverables={[
        'Billing frequencies: DAILY, WEEKLY, MONTHLY, and YEARLY intervals',
        'End-of-month and leap-year recurrence calculations',
        'Plan availability date boundaries vs subscription lifecycle durations',
        'Policy flags: auto-close, allow early close, allow pause, allow renewal',
      ]}
    />
  );
}
