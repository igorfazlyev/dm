import { useState, useEffect } from 'react';
import { Package, FileText } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const ClinicOffers = () => {
  const [offers, setOffers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOffers();
  }, []);

  const fetchOffers = async () => {
    try {
      const response = await api.get('/clinic/offers');
      setOffers(response.data);
    } catch (error) {
      console.error('Failed to fetch offers:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status) => {
    const colors = {
      pending: 'bg-yellow-100 text-yellow-800',
      accepted: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
    };
    return colors[status] || colors.pending;
  };

  if (loading) {
    return <div className="text-center py-12">{t.clinic.offers.loadingOffers}</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{t.clinic.offers.title}</h1>
        <p className="text-gray-600 mt-2">{t.clinic.offers.subtitle}</p>
      </div>

      {offers.length === 0 ? (
        <div className="card text-center py-12">
          <Package className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.clinic.offers.noOffers}</h3>
          <p className="text-gray-600">{t.clinic.offers.noOffersText}</p>
        </div>
      ) : (
        <div className="space-y-4">
          {offers.map((offer) => (
            <div key={offer.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-start gap-3">
                  <div className="p-2 bg-primary-100 rounded-lg">
                    <FileText className="w-5 h-5 text-primary-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">
                      {t.clinic.offers.offer} #{offer.id.slice(0, 8)}
                    </h3>
                    <p className="text-sm text-gray-500 mt-1">
                      {t.common.created}: {new Date(offer.created_at).toLocaleDateString('ru-RU')}
                    </p>
                  </div>
                </div>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(offer.status)}`}>
                  {t.offerStatus[offer.status] || offer.status}
                </span>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="text-xs text-gray-500">{t.clinic.offers.totalPrice}</p>
                  <p className="font-semibold text-gray-900">₽{offer.total_price.toLocaleString('ru-RU')}</p>
                </div>
                {offer.discount_percent > 0 && (
                  <div>
                    <p className="text-xs text-gray-500">{t.clinic.offers.discount}</p>
                    <p className="font-semibold text-green-600">{offer.discount_percent}%</p>
                  </div>
                )}
                {offer.estimated_days && (
                  <div>
                    <p className="text-xs text-gray-500">{t.clinic.offers.duration}</p>
                    <p className="font-semibold text-gray-900">{offer.estimated_days} {t.clinic.offers.days}</p>
                  </div>
                )}
                {offer.has_installment && (
                  <div>
                    <p className="text-xs text-gray-500">{t.clinic.offers.payment}</p>
                    <p className="font-semibold text-gray-900">{t.clinic.offers.installment}</p>
                  </div>
                )}
              </div>

              {offer.special_offer && (
                <div className="mt-3 p-3 bg-yellow-50 rounded-lg">
                  <p className="text-sm text-yellow-800">🎁 {offer.special_offer}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ClinicOffers;
