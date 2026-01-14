import { useState, useEffect } from 'react';
import { ClipboardList, User } from 'lucide-react';
import api from '../../lib/api';
import t from '../../lib/translations';

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      const response = await api.get('/clinic/orders');
      setOrders(response.data);
    } catch (error) {
      console.error('Failed to fetch orders:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateStatus = async (orderId, newStatus) => {
    try {
      await api.patch(`/clinic/orders/${orderId}/status`, { status: newStatus });
      fetchOrders();
    } catch (error) {
      alert(t.errors.failedToUpdateStatus);
    }
  };

  const getStatusColor = (status) => {
    const colors = {
      new: 'bg-blue-100 text-blue-800',
      consultation_scheduled: 'bg-purple-100 text-purple-800',
      in_progress: 'bg-yellow-100 text-yellow-800',
      completed: 'bg-green-100 text-green-800',
      cancelled: 'bg-red-100 text-red-800',
    };
    return colors[status] || colors.new;
  };

  const statusOptions = [
    { value: 'new', label: t.orderStatus.new },
    { value: 'consultation_scheduled', label: t.orderStatus.consultation_scheduled },
    { value: 'in_progress', label: t.orderStatus.in_progress },
    { value: 'completed', label: t.orderStatus.completed },
    { value: 'cancelled', label: t.orderStatus.cancelled },
  ];

  if (loading) {
    return <div className="text-center py-12">{t.clinic.orders.loadingOrders}</div>;
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{t.clinic.orders.title}</h1>
        <p className="text-gray-600 mt-2">{t.clinic.orders.subtitle}</p>
      </div>

      {orders.length === 0 ? (
        <div className="card text-center py-12">
          <ClipboardList className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">{t.clinic.orders.noOrders}</h3>
          <p className="text-gray-600">{t.clinic.orders.noOrdersText}</p>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div key={order.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-start gap-3">
                  <div className="p-2 bg-primary-100 rounded-lg">
                    <User className="w-5 h-5 text-primary-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">
                      {t.clinic.orders.order} #{order.id.slice(0, 8)}
                    </h3>
                    <p className="text-sm text-gray-500 mt-1">
                      {t.common.created}: {new Date(order.created_at).toLocaleDateString('ru-RU')}
                    </p>
                  </div>
                </div>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(order.status)}`}>
                  {t.orderStatus[order.status] || order.status}
                </span>
              </div>

              {order.offer && (
                <div className="p-3 bg-gray-50 rounded-lg mb-4">
                  <p className="text-sm font-medium text-gray-900 mb-2">{t.clinic.orders.offerDetails}</p>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                    <div>
                      <p className="text-gray-500">{t.clinic.offers.totalPrice}</p>
                      <p className="font-semibold text-gray-900">₽{order.offer.total_price.toLocaleString('ru-RU')}</p>
                    </div>
                    {order.offer.discount_percent > 0 && (
                      <div>
                        <p className="text-gray-500">{t.clinic.offers.discount}</p>
                        <p className="font-semibold text-green-600">{order.offer.discount_percent}%</p>
                      </div>
                    )}
                    {order.offer.estimated_days && (
                      <div>
                        <p className="text-gray-500">{t.clinic.offers.duration}</p>
                        <p className="font-semibold text-gray-900">{order.offer.estimated_days} {t.clinic.offers.days}</p>
                      </div>
                    )}
                  </div>
                </div>
              )}

              <div className="pt-4 border-t border-gray-200">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  {t.clinic.orders.updateStatus}
                </label>
                <select
                  value={order.status}
                  onChange={(e) => handleUpdateStatus(order.id, e.target.value)}
                  className="input max-w-xs"
                >
                  {statusOptions.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Orders;
