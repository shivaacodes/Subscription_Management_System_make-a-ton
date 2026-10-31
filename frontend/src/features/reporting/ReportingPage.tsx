import { ModulePlaceholder } from '@/components/ui/ModulePlaceholder';
import { BarChart3 } from 'lucide-react';

export function ReportingPage() {
  return (
    <ModulePlaceholder
      title="Revenue Reporting & Analytics"
      description="Track recurring operational health with distinct separation between MRR, Invoiced Revenue, and Collected Cash."
      icon={BarChart3}
      owner="Member 4 (Quotations & Experience)"
      branch="feat/quotations-portal"
      deliverables={[
        'Monthly Recurring Revenue (MRR) calculated from normalized plan pricing',
        'Invoiced revenue vs actual collected payments separation',
        'Outstanding accounts receivable and overdue invoice aging brackets',
        'Date and status filtering across active customer accounts',
      ]}
    />
  );
}
