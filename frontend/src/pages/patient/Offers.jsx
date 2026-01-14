import { useState, useEffect } from 'react';
import { Package, Building2, Check } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Offers = () => {
  const [offerRequests, setOfferRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [acceptingOffer, setAcceptingOffer] = useState(null);

  useEffect(() => {
    fetchOfferRequests();
  }, []);

  const fetchOfferRequests = async () => {
    try {
      const response = await api.get('/patient/offer-requests');
      setOfferRequests(response.data);
    } catch (error) {
      console.error('Failed to fetch offer requests:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAcceptOffer = async (offerId) => {
    setAcceptingOffer(offerId);
    try {
      await api.post(`/patient/offers/${offerId}/accept`);
      fetchOfferRequests();
    } catch (error) {
      alert(t.errors.failedToAcceptOffer + ': ' + (error.response?.data?.error || t.errors.somethingWentWrong));
    } finally {
      setAcceptingOffer(null);
    }
  };

  const getStatusColor = (status) => {
    const colors = {
      open: 'bg-green-100 text-green-800',
      closed: 'bg-gray-100 text-gray-800',
    };
    return colors[status] || colors.open;
  };

  if (loading) {
    return <div className="text-center py-12">{t.common.loading}</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{t.patient.offers.title}</h1>
        <p className="text-gray-600 mt-2">{t.patient.offers.subtitle}</p>
      </div>

      {offerRequests.length === 0 ? (
        <div className="card text-center py-12">
          <Package className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.patient.offers.noRequests}</h3>
          <p className="text-gray-600">{t.patient.offers.noRequestsText}</p>
        </div>
      ) : (
        <div className="space-y-6">
          {offerRequests.map((request) => (
            <div key={request.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">
                    {t.patient.offers.offerRequest}
                  </h3>
                  <p className="text-sm text-gray-500 mt-1">
                    {t.common.created}: {new Date(request.created_at).toLocaleDateString('ru-RU')}
                  </p>
                </div>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(request.status)}`}>
                  {t.offerStatus[request.status] || request.status}
                </span>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4 p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="text-xs text-gray-500">{t.patient.offers.city}</p>
                  <p className="font-medium text-gray-900">{request.preferred_city || t.common.any}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">{t.patient.offers.district}</p>
                  <p className="font-medium text-gray-900">{request.preferred_district || t.common.any}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">{t.patient.offers.priceSegment}</p>
                  <p className="font-medium text-gray-900 capitalize">
                    {request.price_segment ? t.priceSegments[request.price_segment] : t.common.any}
                  </p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">{t.patient.offers.items}</p>
                  <p className="font-medium text-gray-900">{request.selected_item_ids?.length || 0}</p>
                </div>
              </div>

              {request.offers && request.offers.length > 0 ? (
                <div className="space-y-3">
                  <h4 className="font-medium text-gray-900">
                    {t.patient.offers.offersReceived} ({request.offers.length})
                  </h4>
                  <div className="space-y-3">
                    {request.offers.map((offer) => (
                      <div
                        key={offer.id}
                        className="border border-gray-200 rounded-lg p-4 hover:border-primary-300 transition-colors"
                      >
                        <div className="flex items-start justify-between gap-4">
                          <div className="flex items-start gap-3 flex-1">
                            <div className="p-2 bg-primary-100 rounded-lg">
                              <Building2 className="w-5 h-5 text-primary-600" />
                            </div>
                            <div className="flex-1">
                              <h5 className="font-medium text-gray-900">
                                {offer.clinic?.name || 'Клиника'}
                              </h5>
                              <p className="text-sm text-gray-600 mt-1">
                                {offer.clinic?.district}, {offer.clinic?.city}
                              </p>
                              
                              <div className="grid grid-cols-2 gap-3 mt-3">
                                <div>
                                  <p className="text-xs text-gray-500">{t.patient.offers.totalPrice}</p>
                                  <p className="font-semibold text-gray-900">
                                    ₽{offer.total_price.toLocaleString('ru-RU')}
                                  </p>
                                </div>
                                {offer.discount_percent > 0 && (
                                  <div>
                                    <p className="text-xs text-gray-500">{t.patient.offers.discount}</p>
                                    <p className="font-semibold text-green-600">
                                      {offer.discount_percent}%
                                    </p>
                                  </div>
                                )}
                                {offer.estimated_days && (
                                  <div>
                                    <p className="text-xs text-gray-500">{t.patient.offers.duration}</p>
                                    <p className="font-medium text-gray-900">
                                      {offer.estimated_days} {t.patient.offers.days}
                                    </p>
                                  </div>
                                )}
                                {offer.has_installment && (
                                  <div>
                                    <p className="text-xs text-gray-500">{t.patient.offers.payment}</p>
                                    <p className="font-medium text-gray-900">{t.patient.offers.installment}</p>
                                  </div>
                                )}
                              </div>

                              {offer.special_offer && (
                                <div className="mt-3 p-2 bg-yellow-50 rounded text-sm text-yellow-800">
                                  🎁 {offer.special_offer}
                                </div>
                              )}
                            </div>
                          </div>

                          <div className="flex flex-col gap-2">
                            {offer.status === 'pending' && (
                              <button
                                onClick={() => handleAcceptOffer(offer.id)}
                                disabled={acceptingOffer === offer.id}
                                className="btn-primary flex items-center gap-2 whitespace-nowrap"
                              >
                                <Check className="w-4 h-4" />
                                {acceptingOffer === offer.id ? t.patient.offers.accepting : t.patient.offers.accept}
                              </button>
                            )}
                            {offer.status === 'accepted' && (
                              <span className="px-3 py-2 bg-green-100 text-green-800 rounded-lg text-sm font-medium">
                                {t.patient.offers.accepted}
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              ) : (
                <div className="text-center py-8 bg-gray-50 rounded-lg">
                  <Package className="w-10 h-10 text-gray-400 mx-auto mb-2" />
                  <p className="text-gray-600">{t.patient.offers.noOffers}</p>
                  <p className="text-sm text-gray-500 mt-1">{t.patient.offers.noOffersText}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Offers;
