import { useState, useEffect } from 'react';
import { ClipboardList, Circle } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Plans = () => {
  const [studies, setStudies] = useState([]);
  const [selectedStudy, setSelectedStudy] = useState(null);
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStudies();
  }, []);

  const fetchStudies = async () => {
    try {
      const response = await api.get('/patient/studies');
      setStudies(response.data);
      if (response.data.length > 0) {
        setSelectedStudy(response.data[0].id);
        fetchPlansForStudy(response.data[0].id);
      }
    } catch (error) {
      console.error('Failed to fetch studies:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchPlansForStudy = async (studyId) => {
    try {
      const response = await api.get(`/plans/study/${studyId}`);
      setPlans(response.data);
    } catch (error) {
      console.error('Failed to fetch plans:', error);
      setPlans([]);
    }
  };

  const handleStudyChange = (studyId) => {
    setSelectedStudy(studyId);
    fetchPlansForStudy(studyId);
  };

  const getSpecialtyColor = (specialty) => {
    const colors = {
      therapy: 'bg-blue-100 text-blue-800',
      orthopedics: 'bg-purple-100 text-purple-800',
      surgery: 'bg-red-100 text-red-800',
      hygiene: 'bg-green-100 text-green-800',
      periodontics: 'bg-yellow-100 text-yellow-800',
    };
    return colors[specialty] || 'bg-gray-100 text-gray-800';
  };

  if (loading) {
    return <div className="text-center py-12">{t.common.loading}</div>;
  }

  if (studies.length === 0) {
    return (
      <div className="card text-center py-12">
        <ClipboardList className="w-12 h-12 text-gray-400 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-gray-900 mb-2">{t.patient.plans.noStudies}</h3>
        <p className="text-gray-600">{t.patient.plans.noStudiesText}</p>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{t.patient.plans.title}</h1>
        <p className="text-gray-600 mt-2">{t.patient.plans.subtitle}</p>
      </div>

      <div className="card mb-6">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          {t.patient.plans.selectStudy}
        </label>
        <select
          value={selectedStudy || ''}
          onChange={(e) => handleStudyChange(e.target.value)}
          className="input max-w-md"
        >
          {studies.map((study) => (
            <option key={study.id} value={study.id}>
              {study.modality} - {new Date(study.created_at).toLocaleDateString('ru-RU')}
            </option>
          ))}
        </select>
      </div>

      {plans.length === 0 ? (
        <div className="card text-center py-12">
          <ClipboardList className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.patient.plans.noPlans}</h3>
          <p className="text-gray-600">{t.patient.plans.noPlansText}</p>
        </div>
      ) : (
        <div className="space-y-6">
          {plans.map((plan) => (
            <div key={plan.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">
                    {t.patient.plans.treatmentPlan} {t.patient.plans.version}{plan.version}
                  </h3>
                  <p className="text-sm text-gray-500 mt-1">
                    {t.common.created}: {new Date(plan.created_at).toLocaleDateString('ru-RU')}
                  </p>
                  <span className="inline-block px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800 mt-2">
                    {plan.source}
                  </span>
                </div>
              </div>

              {plan.plan_items && plan.plan_items.length > 0 ? (
                <div className="space-y-3">
                  <h4 className="font-medium text-gray-900 mb-3">
                    {t.patient.plans.procedures} ({plan.plan_items.length})
                  </h4>
                  <div className="space-y-2">
                    {plan.plan_items.map((item) => (
                      <div
                        key={item.id}
                        className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg"
                      >
                        <div className="p-2 bg-white rounded-lg">
                          <Circle className="w-4 h-4 text-gray-600" />
                        </div>
                        <div className="flex-1">
                          <div className="flex items-start justify-between gap-2">
                            <div>
                              <p className="font-medium text-gray-900">{item.procedure_name}</p>
                              {item.diagnosis && (
                                <p className="text-sm text-gray-600 mt-1">{item.diagnosis}</p>
                              )}
                            </div>
                            <span className={`px-2 py-1 rounded-full text-xs font-medium whitespace-nowrap ${getSpecialtyColor(item.specialty)}`}>
                              {t.specialties[item.specialty] || item.specialty}
                            </span>
                          </div>
                          <div className="flex gap-4 mt-2 text-sm text-gray-500">
                            {item.tooth_number && (
                              <span>{t.patient.plans.tooth}: #{item.tooth_number}</span>
                            )}
                            <span>{t.patient.plans.qty}: {item.quantity}</span>
                            {item.procedure_code && (
                              <span>{t.patient.plans.code}: {item.procedure_code}</span>
                            )}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : (
                <p className="text-gray-500 text-center py-4">{t.patient.plans.noProcedures}</p>
              )}

              <div className="mt-4 pt-4 border-t border-gray-200">
                <button className="btn-primary">
                  {t.patient.plans.requestOffers}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Plans;
