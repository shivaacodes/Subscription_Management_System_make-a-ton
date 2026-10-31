import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '@/features/auth/AuthContext';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Layers,
  LayoutDashboard,
  Package,
  CalendarDays,
  FileText,
  Repeat,
  Receipt,
  CreditCard,
  BarChart3,
  ShieldCheck,
  LogOut,
  User as UserIcon,
} from 'lucide-react';

export function AppLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const navItems = [
    { name: 'Dashboard', path: '/admin', icon: LayoutDashboard, roles: ['ADMIN', 'STAFF'] },
    { name: 'Subscriptions', path: '/admin/subscriptions', icon: Repeat, roles: ['ADMIN', 'STAFF'] },
    { name: 'Quotations', path: '/admin/quotations', icon: FileText, roles: ['ADMIN', 'STAFF'] },
    { name: 'Products & Variants', path: '/admin/catalog', icon: Package, roles: ['ADMIN'] },
    { name: 'Subscription Plans', path: '/admin/plans', icon: CalendarDays, roles: ['ADMIN'] },
    { name: 'Billing Engine', path: '/admin/billing', icon: Receipt, roles: ['ADMIN', 'STAFF'] },
    { name: 'Payments', path: '/admin/payments', icon: CreditCard, roles: ['ADMIN', 'STAFF'] },
    { name: 'Revenue Reports', path: '/admin/reporting', icon: BarChart3, roles: ['ADMIN'] },
    { name: 'System Check', path: '/admin/system', icon: ShieldCheck, roles: ['ADMIN'] },
  ];

  const allowedNav = navItems.filter((item) => user && item.roles.includes(user.role));

  return (
    <div className="min-h-screen bg-slate-50 flex">
      {/* Sidebar */}
      <aside className="w-64 bg-slate-900 text-slate-300 flex flex-col justify-between border-r border-slate-800">
        <div>
          {/* Brand header */}
          <div className="h-16 px-6 flex items-center gap-3 border-b border-slate-800">
            <div className="h-8 w-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white font-bold shadow">
              <Layers className="h-5 w-5" />
            </div>
            <div>
              <span className="font-bold text-white tracking-tight text-lg">Veyra</span>
              <span className="text-[10px] block text-indigo-400 font-mono leading-none">OPERATIONS</span>
            </div>
          </div>

          {/* Navigation */}
          <nav className="p-4 space-y-1">
            <div className="px-3 py-2 text-[11px] font-semibold tracking-wider text-slate-400 uppercase">
              Management
            </div>
            {allowedNav.map((item) => {
              const Icon = item.icon;
              return (
                <NavLink
                  key={item.path}
                  to={item.path}
                  end={item.path === '/admin'}
                  className={({ isActive }) =>
                    `flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                      isActive
                        ? 'bg-indigo-600 text-white shadow-sm'
                        : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                    }`
                  }
                >
                  <Icon className="h-4 w-4" />
                  <span>{item.name}</span>
                </NavLink>
              );
            })}
          </nav>
        </div>

        {/* User footer */}
        <div className="p-4 border-t border-slate-800">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 overflow-hidden">
              <div className="h-8 w-8 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center text-slate-300">
                <UserIcon className="h-4 w-4" />
              </div>
              <div className="overflow-hidden">
                <p className="text-xs font-medium text-white truncate">{user?.email}</p>
                <Badge variant={user?.role === 'ADMIN' ? 'indigo' : 'default'} className="mt-0.5 text-[10px] py-0 px-1.5">
                  {user?.role}
                </Badge>
              </div>
            </div>
            <button
              onClick={handleLogout}
              title="Sign Out"
              className="p-1.5 text-slate-400 hover:text-rose-400 rounded hover:bg-slate-800 transition-colors"
            >
              <LogOut className="h-4 w-4" />
            </button>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        <header className="h-16 bg-white border-b border-slate-200 px-8 flex items-center justify-between">
          <div className="flex items-center gap-2 text-sm text-slate-500">
            <span>Workspace</span>
            <span>/</span>
            <span className="font-medium text-slate-800">Brightloop Software</span>
          </div>

          <div className="flex items-center gap-3">
            <span className="text-xs text-slate-500 font-mono">Phase 0 Foundation</span>
            {user?.role === 'ADMIN' && (
              <Link to="/portal">
                <Button variant="outline" size="sm" className="text-xs">
                  Switch to Customer Portal
                </Button>
              </Link>
            )}
          </div>
        </header>

        <main className="flex-1 p-8 overflow-y-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
