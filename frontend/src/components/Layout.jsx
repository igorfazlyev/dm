import { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { 
  Menu, X, Home, FileText, Building2, Package, 
  ClipboardList, LogOut, User
} from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import t from '../lib/translations';

const Layout = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // Patient navigation
  const patientNav = [
    { name: t.nav.dashboard, path: '/patient', icon: Home },
    { name: t.nav.myStudies, path: '/patient/studies', icon: FileText },
    { name: t.nav.treatmentPlans, path: '/patient/plans', icon: ClipboardList },
    { name: t.nav.offerRequests, path: '/patient/offers', icon: Package },
    { name: t.nav.profile, path: '/patient/profile', icon: User },
  ];

  // Clinic navigation
  const clinicNav = [
    { name: t.nav.dashboard, path: '/clinic', icon: Home },
    { name: t.nav.profile, path: '/clinic/profile', icon: Building2 },
    { name: t.nav.pricelist, path: '/clinic/pricelist', icon: FileText },
    { name: t.nav.offers, path: '/clinic/offers', icon: Package },
    { name: t.nav.orders, path: '/clinic/orders', icon: ClipboardList },
  ];

  const navigation = user?.role === 'patient' ? patientNav : clinicNav;

  const NavLink = ({ item }) => {
    const isActive = location.pathname === item.path;
    const Icon = item.icon;

    return (
      <Link
        to={item.path}
        onClick={() => setSidebarOpen(false)}
        className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
          isActive
            ? 'bg-primary-50 text-primary-700 font-medium'
            : 'text-gray-700 hover:bg-gray-100'
        }`}
      >
        <Icon className="w-5 h-5" />
        <span>{item.name}</span>
      </Link>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Mobile header */}
      <div className="lg:hidden fixed top-0 left-0 right-0 bg-white border-b border-gray-200 z-40">
        <div className="flex items-center justify-between px-4 py-3">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-2 rounded-lg hover:bg-gray-100"
          >
            {sidebarOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
          <h1 className="text-lg font-semibold text-gray-900">
            {t.auth.dentalMarketplace}
          </h1>
          <div className="w-10" />
        </div>
      </div>

      {/* Sidebar overlay for mobile */}
      {sidebarOpen && (
        <div
          className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-40"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`fixed top-0 left-0 bottom-0 w-64 bg-white border-r border-gray-200 z-50 transform transition-transform duration-200 ease-in-out ${
          sidebarOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0`}
      >
        <div className="flex flex-col h-full">
          {/* Logo */}
          <div className="px-6 py-4 border-b border-gray-200">
            <h1 className="text-xl font-bold text-primary-600">
              {t.auth.dentalMarketplace}
            </h1>
            <p className="text-sm text-gray-500 mt-1">
              {user?.role === 'patient' ? t.nav.patientPortal : t.nav.clinicPortal}
            </p>
          </div>

          {/* Navigation */}
          <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
            {navigation.map((item) => (
              <NavLink key={item.path} item={item} />
            ))}
          </nav>

          {/* User info & logout */}
          <div className="px-3 py-4 border-t border-gray-200 space-y-2">
            <div className="px-4 py-2">
              <p className="text-sm font-medium text-gray-900">{user?.email}</p>
              <p className="text-xs text-gray-500 capitalize">
                {user?.role === 'patient' ? t.auth.patient : t.auth.clinicManager}
              </p>
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center gap-3 w-full px-4 py-3 text-gray-700 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <LogOut className="w-5 h-5" />
              <span>{t.auth.logout}</span>
            </button>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="lg:pl-64 pt-16 lg:pt-0">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {children}
        </div>
      </main>
    </div>
  );
};

export default Layout;
