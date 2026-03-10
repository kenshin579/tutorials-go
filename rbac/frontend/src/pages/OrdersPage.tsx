import { useEffect, useState } from 'react';
import apiClient from '../api/client';
import { useAuth } from '../auth/AuthContext';
import { usePermission } from '../auth/usePermission';
import PermissionGate from '../components/PermissionGate';

interface OrderProduct {
  id: number;
  name: string;
  price: number;
  status: string;
}

interface OrderUser {
  id: number;
  name: string;
  email: string;
}

interface Order {
  id: number;
  product_id: number;
  product: OrderProduct;
  quantity: number;
  total_price: number;
  status: string;
  ordered_by: number;
  orderer: OrderUser;
  created_at: string;
}

interface ProductOption {
  id: number;
  name: string;
  price: number;
}

const statusBadge: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-700',
  confirmed: 'bg-blue-100 text-blue-700',
  shipped: 'bg-purple-100 text-purple-700',
  completed: 'bg-green-100 text-green-700',
  cancelled: 'bg-red-100 text-red-500',
};

export default function OrdersPage() {
  const { user } = useAuth();
  const { hasRole } = usePermission();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [createModal, setCreateModal] = useState(false);
  const [products, setProducts] = useState<ProductOption[]>([]);
  const [selectedProduct, setSelectedProduct] = useState<number>(0);
  const [quantity, setQuantity] = useState(1);
  const [saving, setSaving] = useState(false);

  const fetchOrders = async () => {
    try {
      const res = await apiClient.get('/orders');
      setOrders(res.data ?? []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  const openCreateModal = async () => {
    try {
      const res = await apiClient.get('/products');
      const active = (res.data ?? []).filter((p: ProductOption & { status: string }) => p.status === 'active');
      setProducts(active);
      if (active.length > 0) setSelectedProduct(active[0].id);
    } catch {
      // ignore
    }
    setQuantity(1);
    setCreateModal(true);
  };

  const handleCreateOrder = async () => {
    setSaving(true);
    try {
      await apiClient.post('/orders', {
        product_id: selectedProduct,
        quantity,
      });
      setCreateModal(false);
      fetchOrders();
    } catch {
      // ignore
    } finally {
      setSaving(false);
    }
  };

  const handleStatusUpdate = async (orderId: number, status: string) => {
    try {
      await apiClient.patch(`/orders/${orderId}/status`, { status });
      fetchOrders();
    } catch {
      // ignore
    }
  };

  const handleCancel = async (orderId: number) => {
    try {
      await apiClient.patch(`/orders/${orderId}/cancel`);
      fetchOrders();
    } catch {
      // ignore
    }
  };

  const renderActions = (order: Order) => {
    const buttons: React.ReactElement[] = [];
    const isAdmin = hasRole('admin');
    const isManager = hasRole('manager');
    const isUser = !isAdmin && !isManager;
    const isOwnOrder = order.ordered_by === user?.id;

    if (order.status === 'pending') {
      if (isAdmin) {
        buttons.push(
          <button key="confirm" onClick={() => handleStatusUpdate(order.id, 'confirmed')}
            className="px-2.5 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition-colors">
            Confirm
          </button>
        );
      }
      if (isAdmin || isManager || (isUser && isOwnOrder)) {
        buttons.push(
          <button key="cancel" onClick={() => handleCancel(order.id)}
            className="px-2.5 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition-colors">
            Cancel
          </button>
        );
      }
    }

    if (order.status === 'confirmed') {
      if (isAdmin || isManager) {
        buttons.push(
          <button key="ship" onClick={() => handleStatusUpdate(order.id, 'shipped')}
            className="px-2.5 py-1 text-xs font-medium text-purple-600 bg-purple-50 rounded hover:bg-purple-100 transition-colors">
            Ship
          </button>
        );
      }
      if (isAdmin || isManager) {
        buttons.push(
          <button key="cancel" onClick={() => handleCancel(order.id)}
            className="px-2.5 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition-colors">
            Cancel
          </button>
        );
      }
    }

    if (order.status === 'shipped' && isAdmin) {
      buttons.push(
        <button key="complete" onClick={() => handleStatusUpdate(order.id, 'completed')}
          className="px-2.5 py-1 text-xs font-medium text-green-600 bg-green-50 rounded hover:bg-green-100 transition-colors">
          Complete
        </button>
      );
    }

    return buttons.length > 0 ? buttons : <span className="text-xs text-gray-400">-</span>;
  };

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('ko-KR', { style: 'currency', currency: 'KRW' }).format(price);

  const selectedProductData = products.find((p) => p.id === selectedProduct);
  const totalPrice = selectedProductData ? selectedProductData.price * quantity : 0;

  if (loading) {
    return <div className="text-gray-500">Loading...</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Orders</h2>
        <PermissionGate permission="orders:create">
          <button
            onClick={openCreateModal}
            className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            + Create Order
          </button>
        </PermissionGate>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Order #</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Product</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Qty</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Total</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Orderer</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.map((order) => (
              <tr key={order.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm text-gray-800 font-medium">#{String(order.id).padStart(3, '0')}</td>
                <td className="px-6 py-4 text-sm text-gray-800">{order.product?.name ?? '-'}</td>
                <td className="px-6 py-4 text-sm text-gray-800">{order.quantity}</td>
                <td className="px-6 py-4 text-sm text-gray-800">{formatPrice(order.total_price)}</td>
                <td className="px-6 py-4">
                  <span className={`px-2.5 py-0.5 rounded-full text-xs font-semibold ${statusBadge[order.status] ?? 'bg-gray-100 text-gray-700'}`}>
                    {order.status}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm text-gray-600">{order.orderer?.name ?? '-'}</td>
                <td className="px-6 py-4 space-x-1.5">{renderActions(order)}</td>
              </tr>
            ))}
            {orders.length === 0 && (
              <tr>
                <td colSpan={7} className="px-6 py-8 text-center text-sm text-gray-500">
                  No orders found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Create Order Modal */}
      {createModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Create Order</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Product</label>
                <select
                  value={selectedProduct}
                  onChange={(e) => setSelectedProduct(Number(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {products.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name} - {formatPrice(p.price)}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Quantity</label>
                <input
                  type="number"
                  min={1}
                  value={quantity}
                  onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div className="pt-2 border-t border-gray-200">
                <div className="flex justify-between text-sm">
                  <span className="font-medium text-gray-600">Total:</span>
                  <span className="font-bold text-gray-800">{formatPrice(totalPrice)}</span>
                </div>
              </div>
            </div>
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={() => setCreateModal(false)}
                className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleCreateOrder}
                disabled={saving || !selectedProduct}
                className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                {saving ? 'Creating...' : 'Create Order'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
