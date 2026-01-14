import { useState, useEffect } from 'react';
import { Building2, Save, AlertCircle } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const ClinicProfile = () => {
  const [profile, setProfile] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    legal_name: '',
    license_number: '',
    year_established: new Date().getFullYear(),
    city: '',
    district: '',
    address: '',
    phone: '',
    email: '',
    website: '',
    price_segment: 'business',
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const response = await api.get('/clinic/profile');
      setProfile(response.data);
      setFormData({
        name: response.data.name || '',
        legal_name: response.data.legal_name || '',
        license_number: response.data.license_number || '',
        year_established: response.data.year_established || new Date().getFullYear(),
        city: response.data.city || '',
        district: response.data.district || '',
        address: response.data.address || '',
        phone: response.data.phone || '',
        email: response.data.email || '',
        website: response.data.website || '',
        price_segment: response.data.price_segment || 'business',
      });
      setIsCreating(false);
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
        await api.post('/clinic/profile', formData);
        setMessage(t.clinic.profile.profileCreated);
      } else {
        await api.put('/clinic/profile', formData);
        setMessage(t.clinic.profile.profileUpdated);
      }
      fetchProfile();
    } catch (err) {
      setError(err.response?.data?.error || t.errors.failedToUpdate);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="text-center py-12">{t.clinic.profile.loadingProfile}</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{t.clinic.profile.title}</h1>
        <p className="text-gray-600 mt-2">
          {isCreating ? t.clinic.profile.createSubtitle : t.clinic.profile.manageSubtitle}
        </p>
      </div>

      <div className="max-w-3xl">
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
            <div>
              <h3 className="text-lg font-medium text-gray-900 mb-4">{t.clinic.profile.basicInfo}</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.profile.clinicName} *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="input"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.profile.legalName} *
                  </label>
                  <input
                    type="text"
                    value={formData.legal_name}
                    onChange={(e) => setFormData({ ...formData, legal_name: e.target.value })}
                    className="input"
                    required
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t.clinic.profile.licenseNumber} *
                    </label>
                    <input
                      type="text"
                      value={formData.license_number}
                      onChange={(e) => setFormData({ ...formData, license_number: e.target.value })}
                      className="input"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t.clinic.profile.yearEstablished} *
                    </label>
                    <input
                      type="number"
                      value={formData.year_established}
                      onChange={(e) => setFormData({ ...formData, year_established: parseInt(e.target.value) })}
                      className="input"
                      min="1900"
                      max={new Date().getFullYear()}
                      required
                    />
                  </div>
                </div>
              </div>
            </div>

            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">{t.clinic.profile.location}</h3>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t.clinic.profile.city} *
                    </label>
                    <input
                      type="text"
                      value={formData.city}
                      onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                      className="input"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t.clinic.profile.district} *
                    </label>
                    <input
                      type="text"
                      value={formData.district}
                      onChange={(e) => setFormData({ ...formData, district: e.target.value })}
                      className="input"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.profile.address} *
                  </label>
                  <input
                    type="text"
                    value={formData.address}
                    onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                    className="input"
                    required
                  />
                </div>
              </div>
            </div>

            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">{t.clinic.profile.contactInfo}</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.profile.phone} *
                  </label>
                  <input
                    type="tel"
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                    className="input"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.profile.email} *
                  </label>
                  <input
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    className="input"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.profile.website}
                  </label>
                  <input
                    type="url"
                    value={formData.website}
                    onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                    className="input"
                    placeholder="https://example.com"
                  />
                </div>
              </div>
            </div>

            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">{t.clinic.profile.marketPositioning}</h3>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.clinic.profile.priceSegment} *
                </label>
                <select
                  value={formData.price_segment}
                  onChange={(e) => setFormData({ ...formData, price_segment: e.target.value })}
                  className="input"
                  required
                >
                  <option value="economy">{t.priceSegments.economy}</option>
                  <option value="business">{t.priceSegments.business}</option>
                  <option value="premium">{t.priceSegments.premium}</option>
                </select>
              </div>
            </div>

            <div className="flex justify-end pt-4">
              <button
                type="submit"
                disabled={saving}
                className="btn-primary flex items-center gap-2"
              >
                <Save className="w-5 h-5" />
                {saving ? t.common.saving : isCreating ? t.clinic.profile.createProfile : t.patient.profile.saveChanges}
              </button>
            </div>
          </form>

          {profile && !profile.is_active && (
            <div className="mt-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
              <div className="flex items-start gap-2">
                <AlertCircle className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-yellow-800">{t.clinic.profile.pendingActivation}</p>
                  <p className="text-sm text-yellow-700 mt-1">
                    {t.clinic.profile.pendingActivationText}
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ClinicProfile;
