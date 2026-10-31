import { ModulePlaceholder } from '@/components/ui/ModulePlaceholder';
import { Repeat } from 'lucide-react';

export function SubscriptionsPage() {
  return (
    <ModulePlaceholder
      title="Subscription Management"
      description="Manage recurring subscription lifecycles, renewal intervals, customer assignments, and add-on services."
      icon={Repeat}
      owner="Member 1 (Team Lead)"
      branch="feat/subscription-lifecycle"
      deliverables={[
        'Subscription state machine transitions (DRAFT -> QUOTATION -> CONFIRMED -> ACTIVE -> CLOSED)',
        'Pause and cancellation policy validation against plan settings',
        'Add-on activation/deactivation taking effect at next billing cycle',
        'Immutable subscription event timeline (subscription_events logging)',
      ]}
    />
  );
}
