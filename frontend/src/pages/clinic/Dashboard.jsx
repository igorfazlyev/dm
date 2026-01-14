import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Package, ClipboardList, FileText, Building2, DollarSign } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const ClinicDashboard = () => {
  const [stats, setStats] = useState({
    offers: 0,
    orders: 0,
    pricelistItems: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [offersRes, ordersRes, pricelistRes] = await Promise.all([
          api.get('/clinic/offers'),
          api.get('/clinic/orders'),
          api.get('/clinic/pricelist'),
        ]);

        setStats({
          offers: offersRes.data.length,
          orders: ordersRes.data.length,
          pricelistItems: pricelistRes.data.length,
        });
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const StatCard = ({ icon: Icon, label, value, linkTo, color }) => (
    <Link to={linkTo} className="card hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600">{label}</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
        </div>
        <div className={`p-3 rounded-lg ${color}`}>
          <Icon className="w-6 h-6 text-white" />
        </div>
      </div>
    </Link>
  );

  if (loading) {
    return <div className="text-center py-12">{t.common.loading}</div>;
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">{t.clinic.dashboard.title}</h1>
        <p className="text-gray-600 mt-2">{t.clinic.dashboard.subtitle}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <StatCard
          icon={Package}
          label={t.clinic.dashboard.myOffers}
          value={stats.offers}
          linkTo="/clinic/offers"
          color="bg-purple-500"
        />
        <StatCard
          icon={ClipboardList}
          label={t.clinic.dashboard.orders}
          value={stats.orders}
          linkTo="/clinic/orders"
          color="bg-green-500"
        />
        <StatCard
          icon={DollarSign}
          label={t.clinic.dashboard.pricelistItems}
          value={stats.pricelistItems}
          linkTo="/clinic/pricelist"
          color="bg-blue-500"
        />
      </div>

      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">{t.clinic.dashboard.quickActions}</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <Link
            to="/clinic/profile"
            className="flex items-center gap-3 p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-400 hover:bg-primary-50 transition-colors"
          >
            <Building2 className="w-5 h-5 text-primary-600" />
            <span className="font-medium text-gray-700">{t.clinic.dashboard.updateProfile}</span>
          </Link>
          <Link
            to="/clinic/pricelist"
            className="flex items-center gap-3 p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-400 hover:bg-primary-50 transition-colors"
          >
            <FileText className="w-5 h-5 text-primary-600" />
            <span className="font-medium text-gray-700">{t.clinic.dashboard.managePricelist}</span>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default ClinicDashboard;
