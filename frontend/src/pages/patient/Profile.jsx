import { useState, useEffect } from 'react';
import { Save, AlertCircle } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Profile = () => {
  const [profile, setProfile] = useState(null);
  const [formData, setFormData] = useState({
    first_name: '',
    last_name: '',
    phone: '',
    preferred_city: '',
    preferred_district: '',
    preferred_price_segment: '',
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    fetchProfile();
  }, []);

  // const fetchProfile = async () => {
  //   try {
  //     const response = await api.get('/patient/profile');
  //     setProfile(response.data);
  //     setFormData({
  //       first_name: response.data.first_name || '',
  //       last_name: response.data.last_name || '',
  //       phone: response.data.phone || '',
  //       preferred_city: response.data.preferred_city || '',
  //       preferred_district: response.data.preferred_district || '',
  //       preferred_price_segment: response.data.preferred_price_segment || '',
  //     });
  //   } catch (error) {
  //     console.error('Failed to fetch profile:', error);
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  // const handleSubmit = async (e) => {
  //   e.preventDefault();
  //   setError('');
  //   setMessage('');
  //   setSaving(true);

  //   try {
  //     await api.put('/patient/profile', formData);
  //     setMessage(t.patient.profile.profileUpdated);
  //     fetchProfile();
  //   } catch (err) {
  //     setError(err.response?.data?.error || t.errors.failedToUpdate);
  //   } finally {
  //     setSaving(false);
  //   }
  // };

  const fetchProfile = async () => {
  try {
    const response = await api.get('/patient/profile');
    
    // Check if profile exists
    if (response.data.exists === false) {
      setIsCreating(true);
      setProfile(null);
    } else {
      setProfile(response.data);
      setFormData({
        first_name: response.data.first_name || '',
        last_name: response.data.last_name || '',
        date_of_birth: response.data.date_of_birth?.split('T')[0] || '',
        phone: response.data.phone || '',
        preferred_city: response.data.preferred_city || '',
        preferred_district: response.data.preferred_district || '',
        preferred_price_segment: response.data.preferred_price_segment || 'business',
      });
        setIsCreating(false);
      }
    } catch (error) {
      if (error.response?.status === 404) {
        setIsCreating(true);
      }
    } finally {
      setLoading(false);
    }
  };

    const handleSubmit = async (e) => {
      e.preventDefault();
      setError('');
      setMessage('');
      setSaving(true);

      try {
        if (isCreating) {
          await api.post('/patient/profile', formData);
          setMessage(t.patient.profile.profileCreated);
        } else {
          await api.put('/patient/profile', formData);
          setMessage(t.patient.profile.profileUpdated);
        }
        fetchProfile();
      } catch (err) {
        setError(err.response?.data?.error || t.errors.failedToUpdate);
      } finally {
        setSaving(false);
      }
    };


  if (loading) {
    return <div className="text-center py-12">{t.patient.profile.loadingProfile}</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{t.patient.profile.title}</h1>
        <p className="text-gray-600 mt-2">{t.patient.profile.subtitle}</p>
      </div>

      <div className="max-w-2xl">
        <div className="card">
          {message && (
            <div className="bg-green-50 border border-green-200 rounded-lg p-3 flex items-start gap-2 mb-4">
              <Save className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-green-600">{message}</p>
            </div>
          )}

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3 flex items-start gap-2 mb-4">
              <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.patient.profile.firstName}
                </label>
                <input
                  type="text"
                  value={formData.first_name}
                  onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.patient.profile.lastName}
                </label>
                <input
                  type="text"
                  value={formData.last_name}
                  onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                  className="input"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                {t.patient.profile.phone}
              </label>
              <input
                type="tel"
                value={formData.phone}
                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                className="input"
                placeholder="+7 (999) 123-45-67"
              />
            </div>

            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">{t.patient.profile.preferences}</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.patient.profile.preferredCity}
                  </label>
                  <input
                    type="text"
                    value={formData.preferred_city}
                    onChange={(e) => setFormData({ ...formData, preferred_city: e.target.value })}
                    className="input"
                    placeholder="Москва"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.patient.profile.preferredDistrict}
                  </label>
                  <input
                    type="text"
                    value={formData.preferred_district}
                    onChange={(e) => setFormData({ ...formData, preferred_district: e.target.value })}
                    className="input"
                    placeholder="Центральный"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.patient.profile.priceSegment}
                  </label>
                  <select
                    value={formData.preferred_price_segment}
                    onChange={(e) => setFormData({ ...formData, preferred_price_segment: e.target.value })}
                    className="input"
                  >
                    <option value="">{t.common.any}</option>
                    <option value="economy">{t.priceSegments.economy}</option>
                    <option value="business">{t.priceSegments.business}</option>
                    <option value="premium">{t.priceSegments.premium}</option>
                  </select>
                </div>
              </div>
            </div>

            <div className="flex justify-end pt-4">
              <button
                type="submit"
                disabled={saving}
                className="btn-primary flex items-center gap-2"
              >
                <Save className="w-5 h-5" />
                {saving ? t.common.saving : t.patient.profile.saveChanges}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Profile;
