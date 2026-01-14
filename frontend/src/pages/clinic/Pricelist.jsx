import { useState, useEffect } from 'react';
import { Plus, DollarSign, Trash2, AlertCircle } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Pricelist = () => {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [formData, setFormData] = useState({
    specialty: 'therapy',
    procedure_code: '',
    procedure_name: '',
    price_from: '',
    price_to: '',
  });
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    fetchPricelist();
  }, []);

  const fetchPricelist = async () => {
    try {
      const response = await api.get('/clinic/pricelist');
      setItems(response.data);
    } catch (error) {
      console.error('Failed to fetch pricelist:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddItem = async (e) => {
    e.preventDefault();
    setError('');
    setSaving(true);

    try {
      await api.post('/clinic/pricelist', {
        ...formData,
        price_from: parseFloat(formData.price_from),
        price_to: parseFloat(formData.price_to),
      });
      setShowAddModal(false);
      setFormData({
        specialty: 'therapy',
        procedure_code: '',
        procedure_name: '',
        price_from: '',
        price_to: '',
      });
      fetchPricelist();
    } catch (err) {
      setError(err.response?.data?.error || t.errors.failedToCreate);
    } finally {
      setSaving(false);
    }
  };

  const handleDeleteItem = async (id) => {
    if (!confirm(t.clinic.pricelist.deleteConfirm)) return;

    try {
      await api.delete(`/clinic/pricelist/${id}`);
      fetchPricelist();
    } catch (error) {
      alert(t.errors.failedToDelete);
    }
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
    return <div className="text-center py-12">{t.clinic.pricelist.loadingPricelist}</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{t.clinic.pricelist.title}</h1>
          <p className="text-gray-600 mt-2">{t.clinic.pricelist.subtitle}</p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-5 h-5" />
          {t.clinic.pricelist.addItem}
        </button>
      </div>

      {items.length === 0 ? (
        <div className="card text-center py-12">
          <DollarSign className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.clinic.pricelist.noItems}</h3>
          <p className="text-gray-600 mb-4">{t.clinic.pricelist.noItemsText}</p>
          <button
            onClick={() => setShowAddModal(true)}
            className="btn-primary inline-flex items-center gap-2"
          >
            <Plus className="w-5 h-5" />
            {t.clinic.pricelist.addFirstItem}
          </button>
        </div>
      ) : (
        <div className="card">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">{t.clinic.pricelist.specialty}</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">{t.clinic.pricelist.procedure}</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">{t.clinic.pricelist.code}</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">{t.clinic.pricelist.priceRange}</th>
                  <th className="text-left py-3 px-4 text-sm font-medium text-gray-700">{t.common.status}</th>
                  <th className="text-right py-3 px-4 text-sm font-medium text-gray-700">{t.common.actions}</th>
                </tr>
              </thead>
              <tbody>
                {items.map((item) => (
                  <tr key={item.id} className="border-b border-gray-100 hover:bg-gray-50">
                    <td className="py-3 px-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${getSpecialtyColor(item.specialty)}`}>
                        {t.specialties[item.specialty] || item.specialty}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-sm text-gray-900">{item.procedure_name}</td>
                    <td className="py-3 px-4 text-sm text-gray-600">{item.procedure_code}</td>
                    <td className="py-3 px-4 text-sm font-medium text-gray-900">
                      ₽{item.price_from.toLocaleString('ru-RU')} - ₽{item.price_to.toLocaleString('ru-RU')}
                    </td>
                    <td className="py-3 px-4">
                      <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                        item.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                      }`}>
                        {item.is_active ? t.clinic.pricelist.active : t.clinic.pricelist.inactive}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-right">
                      <button
                        onClick={() => handleDeleteItem(item.id)}
                        className="text-red-600 hover:text-red-700 p-2 rounded-lg hover:bg-red-50"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Add Item Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">{t.clinic.pricelist.addPricelistItem}</h2>

            {error && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-3 flex items-start gap-2 mb-4">
                <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
                <p className="text-sm text-red-600">{error}</p>
              </div>
            )}

            <form onSubmit={handleAddItem} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.clinic.pricelist.specialty}
                </label>
                <select
                  value={formData.specialty}
                  onChange={(e) => setFormData({ ...formData, specialty: e.target.value })}
                  className="input"
                  required
                >
                  <option value="therapy">{t.specialties.therapy}</option>
                  <option value="orthopedics">{t.specialties.orthopedics}</option>
                  <option value="surgery">{t.specialties.surgery}</option>
                  <option value="hygiene">{t.specialties.hygiene}</option>
                  <option value="periodontics">{t.specialties.periodontics}</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.clinic.pricelist.procedureName}
                </label>
                <input
                  type="text"
                  value={formData.procedure_name}
                  onChange={(e) => setFormData({ ...formData, procedure_name: e.target.value })}
                  className="input"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t.clinic.pricelist.procedureCode}
                </label>
                <input
                  type="text"
                  value={formData.procedure_code}
                  onChange={(e) => setFormData({ ...formData, procedure_code: e.target.value })}
                  className="input"
                  placeholder="например, D2391"
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.pricelist.priceFrom}
                  </label>
                  <input
                    type="number"
                    value={formData.price_from}
                    onChange={(e) => setFormData({ ...formData, price_from: e.target.value })}
                    className="input"
                    min="0"
                    step="0.01"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {t.clinic.pricelist.priceTo}
                  </label>
                  <input
                    type="number"
                    value={formData.price_to}
                    onChange={(e) => setFormData({ ...formData, price_to: e.target.value })}
                    className="input"
                    min="0"
                    step="0.01"
                    required
                  />
                </div>
              </div>

              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => {
                    setShowAddModal(false);
                    setError('');
                  }}
                  className="btn-secondary flex-1"
                  disabled={saving}
                >
                  {t.common.cancel}
                </button>
                <button
                  type="submit"
                  className="btn-primary flex-1"
                  disabled={saving}
                >
                  {saving ? t.clinic.pricelist.adding : t.clinic.pricelist.addItem}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Pricelist;
