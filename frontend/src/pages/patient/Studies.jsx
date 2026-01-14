import { useState, useEffect } from 'react';
import { Plus, FileText, Calendar, AlertCircle } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Studies = () => {
  const [studies, setStudies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [formData, setFormData] = useState({
    modality: 'CBCT',
    study_date: new Date().toISOString().split('T')[0],
  });
  const [error, setError] = useState('');
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    fetchStudies();
  }, []);

  const fetchStudies = async () => {
    try {
      const response = await api.get('/patient/studies');
      setStudies(response.data);
    } catch (error) {
      console.error('Failed to fetch studies:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateStudy = async (e) => {
    e.preventDefault();
    setError('');
    setCreating(true);

    try {
      await api.post('/studies', formData);
      setShowCreateModal(false);
      setFormData({
        modality: 'CBCT',
        study_date: new Date().toISOString().split('T')[0],
      });
      fetchStudies();
    } catch (err) {
      setError(err.response?.data?.error || t.errors.failedToCreate);
    } finally {
      setCreating(false);
    }
  };

  const getStatusColor = (status) => {
    const colors = {
      created: 'bg-gray-100 text-gray-800',
      uploading: 'bg-blue-100 text-blue-800',
      processing: 'bg-yellow-100 text-yellow-800',
      completed: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
    };
    return colors[status] || colors.created;
  };

  if (loading) {
    return <div className="text-center py-12">{t.common.loading}</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{t.patient.studies.title}</h1>
          <p className="text-gray-600 mt-2">{t.patient.studies.subtitle}</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          {t.patient.studies.newStudy}
        </button>
      </div>

      {studies.length === 0 ? (
        <div className="card text-center py-12">
          <FileText className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.patient.studies.noStudies}</h3>
          <p className="text-gray-600 mb-4">{t.patient.studies.noStudiesText}</p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="btn-primary inline-flex items-center gap-2"
          >
            <Plus className="w-5 h-5" />
            {t.patient.studies.createStudy}
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {studies.map((study) => (
            <div key={study.id} className="card hover:shadow-md transition-shadow">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-primary-100 rounded-lg">
                    <FileText className="w-5 h-5 text-primary-600" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{study.modality || 'CBCT'}</p>
                    <p className="text-sm text-gray-500">
                      {study.study_date ? new Date(study.study_date).toLocaleDateString('ru-RU') : 'Нет даты'}
                    </p>
                  </div>
                </div>
                <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(study.status)}`}>
                  {t.studyStatus[study.status] || study.status}
                </span>
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2 text-gray-600">
                  <Calendar className="w-4 h-4" />
                  {t.patient.studies.createdDate}: {new Date(study.created_at).toLocaleDateString('ru-RU')}
                </div>
              </div>

              <div className="mt-4 pt-4 border-t border-gray-200">
                <button className="text-primary-600 hover:text-primary-700 font-medium text-sm">
                  {t.common.viewDetails}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Create Study Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">{t.patient.studies.createNewStudy}</h2>

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-3 flex items-start gap-2 mb-4">
                <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                <p className="text-sm text-red-600">{error}</p>
              </div>
            )}

            <form onSubmit={handleCreateStudy} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.patient.studies.modality}
                </label>
                <select
                  value={formData.modality}
                  onChange={(e) => setFormData({ ...formData, modality: e.target.value })}
                  className="input"
                >
                  <option value="CBCT">CBCT</option>
                  <option value="CT">CT</option>
                  <option value="Panoramic">Panoramic</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.patient.studies.studyDate}
                </label>
                <input
                  type="date"
                  value={formData.study_date}
                  onChange={(e) => setFormData({ ...formData, study_date: e.target.value })}
                  className="input"
                  required
                />
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowCreateModal(false);
                    setError('');
                  }}
                  className="btn-secondary flex-1"
                  disabled={creating}
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  className="btn-primary flex-1"
                  disabled={creating}
                >
                  {creating ? t.patient.studies.creating : t.patient.studies.createStudy}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Studies;
