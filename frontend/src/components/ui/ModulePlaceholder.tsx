import { Badge } from '@/components/ui/badge';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { GitBranch, User, CheckCircle2, LucideIcon } from 'lucide-react';

interface ModulePlaceholderProps {
  title: string;
  description: string;
  icon: LucideIcon;
  owner: string;
  branch: string;
  deliverables: string[];
}

export function ModulePlaceholder({
  title,
  description,
  icon: Icon,
  owner,
  branch,
  deliverables,
}: ModulePlaceholderProps) {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">{title}</h1>
          <p className="text-sm text-slate-500 mt-1">{description}</p>
        </div>
      </div>

      <Card className="border-dashed border-2 border-slate-200 bg-slate-50/50">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-lg bg-indigo-100 text-indigo-700 flex items-center justify-center">
              <Icon className="h-5 w-5" />
            </div>
            <div>
              <CardTitle className="text-lg">Module Scheduled for Feature Branch</CardTitle>
              <CardDescription>
                Phase 0 establishes the database models, tables, and API contracts for this module.
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="bg-white p-4 rounded-md border border-slate-200">
              <div className="flex items-center gap-2 text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                <User className="h-3.5 w-3.5 text-indigo-600" /> Assigned Member
              </div>
              <div className="text-sm font-medium text-slate-900">{owner}</div>
            </div>

            <div className="bg-white p-4 rounded-md border border-slate-200">
              <div className="flex items-center gap-2 text-xs font-semibold text-slate-500 uppercase tracking-wider mb-1">
                <GitBranch className="h-3.5 w-3.5 text-indigo-600" /> Target Feature Branch
              </div>
              <Badge variant="indigo" className="font-mono text-xs">
                {branch}
              </Badge>
            </div>
          </div>

          <div className="bg-white p-4 rounded-md border border-slate-200">
            <div className="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-3">
              Phase 1 Deliverables
            </div>
            <ul className="space-y-2">
              {deliverables.map((item, idx) => (
                <li key={idx} className="flex items-start gap-2 text-sm text-slate-700">
                  <CheckCircle2 className="h-4 w-4 text-slate-400 mt-0.5 flex-shrink-0" />
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
