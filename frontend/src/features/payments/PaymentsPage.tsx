import { ModulePlaceholder } from '@/components/ui/ModulePlaceholder';
import { CreditCard } from 'lucide-react';

export function PaymentsPage() {
  return (
    <ModulePlaceholder
      title="Payments & Receivables"
      description="Record multi-part payments against invoices, manage balances, and monitor overdue accounts."
      icon={CreditCard}
      owner="Member 3 (Billing & Payments)"
      branch="feat/billing-engine"
      deliverables={[
        'Multiple payment records linked to individual invoices',
        'Derived invoice payment state: PAID, PARTIALLY_PAID, UNPAID, OVERDUE',
        'Database constraint preventing payment sums exceeding invoice total balance',
        'Receivables tracking and reconciliation reporting',
      ]}
    />
  );
}
