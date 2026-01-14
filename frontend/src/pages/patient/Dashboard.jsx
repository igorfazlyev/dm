import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { FileText, Package, ClipboardList, Plus } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const PatientDashboard = () => {
  const [stats, setStats] = useState({
    studies: 0,
    plans: 0,
    offerRequests: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [studiesRes, offersRes] = await Promise.all([
          api.get('/patient/studies'),
          api.get('/patient/offer-requests'),
        ]);

        setStats({
          studies: studiesRes.data.length,
          plans: studiesRes.data.filter(s => s.plan_versions?.length > 0).length,
          offerRequests: offersRes.data.length,
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
        <h1 className="text-3xl font-bold text-gray-900">{t.patient.dashboard.title}</h1>
        <p className="text-gray-600 mt-2">{t.patient.dashboard.welcome}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <StatCard
          icon={FileText}
          label={t.patient.dashboard.myStudies}
          value={stats.studies}
          linkTo="/patient/studies"
          color="bg-blue-500"
        />
        <StatCard
          icon={ClipboardList}
          label={t.patient.dashboard.treatmentPlans}
          value={stats.plans}
          linkTo="/patient/plans"
          color="bg-green-500"
        />
        <StatCard
          icon={Package}
          label={t.patient.dashboard.offerRequests}
          value={stats.offerRequests}
          linkTo="/patient/offers"
          color="bg-purple-500"
        />
      </div>

      <div className="card">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">{t.patient.dashboard.quickActions}</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <Link
            to="/patient/studies"
            className="flex items-center gap-3 p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-400 hover:bg-primary-50 transition-colors"
          >
            <Plus className="w-5 h-5 text-primary-600" />
            <span className="font-medium text-gray-700">{t.patient.dashboard.createNewStudy}</span>
          </Link>
          <Link
            to="/patient/offers"
            className="flex items-center gap-3 p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-400 hover:bg-primary-50 transition-colors"
          >
            <Package className="w-5 h-5 text-primary-600" />
            <span className="font-medium text-gray-700">{t.patient.dashboard.requestOffers}</span>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default PatientDashboard;
