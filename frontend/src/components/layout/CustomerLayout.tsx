import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '@/features/auth/AuthContext';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Layers, Repeat, Receipt, User, LogOut } from 'lucide-react';

export function CustomerLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  const navLinks = [
    { name: 'My Subscriptions', path: '/portal', icon: Repeat },
    { name: 'Invoices & Receipts', path: '/portal/invoices', icon: Receipt },
    { name: 'Profile & Account', path: '/portal/account', icon: User },
  ];

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-between">
      <div>
        <header className="bg-white border-b border-slate-200 sticky top-0 z-10">
          <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
            <div className="flex items-center gap-8">
              <Link to="/portal" className="flex items-center gap-2 text-indigo-600 font-bold text-xl tracking-tight">
                <div className="h-8 w-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white shadow-sm">
                  <Layers className="h-5 w-5" />
                </div>
                <span>Veyra Customer Portal</span>
              </Link>

              <nav className="hidden md:flex items-center space-x-1">
                {navLinks.map((item) => (
                  <NavLink
                    key={item.path}
                    to={item.path}
                    end={item.path === '/portal'}
                    className={({ isActive }) =>
                      `px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                        isActive
                          ? 'text-indigo-600 bg-indigo-50/70 font-semibold'
                          : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
                      }`
                    }
                  >
                    {item.name}
                  </NavLink>
                ))}
              </nav>
            </div>

            <div className="flex items-center gap-3">
              <div className="text-right hidden sm:block">
                <div className="text-xs font-medium text-slate-800">{user?.email}</div>
                <Badge variant="outline" className="text-[10px] py-0 px-1.5">
                  CUSTOMER
                </Badge>
              </div>

              {user?.role === 'ADMIN' && (
                <Link to="/admin">
                  <Button variant="outline" size="sm" className="text-xs">
                    Admin Console
                  </Button>
                </Link>
              )}

              <Button variant="ghost" size="sm" onClick={handleLogout} className="text-slate-600">
                <LogOut className="h-4 w-4 mr-1.5" />
                Sign Out
              </Button>
            </div>
          </div>
        </header>

        <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Outlet />
        </main>
      </div>

      <footer className="border-t border-slate-200 bg-white py-6 text-center text-xs text-slate-500">
        Brightloop Software &bull; Veyra Customer Self-Service &bull; Phase 0 Foundation
      </footer>
    </div>
  );
}
