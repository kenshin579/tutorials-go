import { useEffect, useState } from 'react';
import apiClient from '../api/client';
import { useAuth } from '../auth/AuthContext';
import PermissionGate from '../components/PermissionGate';

interface ProductCreator {
  id: number;
  name: string;
  email: string;
}

interface Product {
  id: number;
  name: string;
  price: number;
  status: string;
  created_by: number;
  creator: ProductCreator;
  created_at: string;
  updated_at: string;
}

interface ModalState {
  open: boolean;
  product: Product | null;
}

export default function ProductsPage() {
  const { user } = useAuth();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [modal, setModal] = useState<ModalState>({ open: false, product: null });
  const [form, setForm] = useState({ name: '', price: '', status: 'active' });
  const [saving, setSaving] = useState(false);

  const fetchProducts = async () => {
    try {
      const res = await apiClient.get('/products');
      setProducts(res.data ?? []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const openCreate = () => {
    setForm({ name: '', price: '', status: 'active' });
    setModal({ open: true, product: null });
  };

  const openEdit = (product: Product) => {
    setForm({ name: product.name, price: String(product.price), status: product.status });
    setModal({ open: true, product });
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      if (modal.product) {
        await apiClient.put(`/products/${modal.product.id}`, {
          name: form.name,
          price: parseFloat(form.price),
        });
      } else {
        await apiClient.post('/products', {
          name: form.name,
          price: parseFloat(form.price),
        });
      }
      setModal({ open: false, product: null });
      fetchProducts();
    } catch {
      // ignore
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this product?')) return;
    try {
      await apiClient.delete(`/products/${id}`);
      fetchProducts();
    } catch {
      // ignore
    }
  };

  const handleStatusToggle = async (product: Product) => {
    const newStatus = product.status === 'active' ? 'inactive' : 'active';
    try {
      await apiClient.patch(`/products/${product.id}/status`, { status: newStatus });
      fetchProducts();
    } catch {
      // ignore
    }
  };

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('ko-KR', { style: 'currency', currency: 'KRW' }).format(price);

  if (loading) {
    return <div className="text-gray-500">Loading...</div>;
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Products</h2>
        <PermissionGate permission="products:create">
          <button
            onClick={openCreate}
            className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
          >
            + Add Product
          </button>
        </PermissionGate>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Name</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Price</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Creator</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {products.map((product) => (
              <tr key={product.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm text-gray-800">{product.name}</td>
                <td className="px-6 py-4 text-sm text-gray-800">{formatPrice(product.price)}</td>
                <td className="px-6 py-4">
                  <span
                    className={`px-2.5 py-0.5 rounded-full text-xs font-semibold cursor-pointer ${
                      product.status === 'active'
                        ? 'bg-green-100 text-green-700'
                        : 'bg-gray-100 text-gray-500'
                    }`}
                    onClick={() => {
                      if (user?.permissions?.includes('products:status:update')) {
                        handleStatusToggle(product);
                      }
                    }}
                  >
                    {product.status}
                  </span>
                </td>
                <td className="px-6 py-4 text-sm text-gray-600">{product.creator?.name ?? '-'}</td>
                <td className="px-6 py-4 space-x-2">
                  <PermissionGate permission="products:update" ownerId={product.created_by}>
                    <button
                      onClick={() => openEdit(product)}
                      className="px-3 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition-colors"
                    >
                      Edit
                    </button>
                  </PermissionGate>
                  <PermissionGate permission="products:delete">
                    <button
                      onClick={() => handleDelete(product.id)}
                      className="px-3 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition-colors"
                    >
                      Delete
                    </button>
                  </PermissionGate>
                </td>
              </tr>
            ))}
            {products.length === 0 && (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-sm text-gray-500">
                  No products found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Modal */}
      {modal.open && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">
              {modal.product ? 'Edit Product' : 'Add Product'}
            </h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Price</label>
                <input
                  type="number"
                  value={form.price}
                  onChange={(e) => setForm({ ...form, price: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              {modal.product && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                  <div className="flex gap-4">
                    <label className="flex items-center gap-2 text-sm">
                      <input
                        type="radio"
                        name="status"
                        checked={form.status === 'active'}
                        onChange={() => setForm({ ...form, status: 'active' })}
                      />
                      active
                    </label>
                    <label className="flex items-center gap-2 text-sm">
                      <input
                        type="radio"
                        name="status"
                        checked={form.status === 'inactive'}
                        onChange={() => setForm({ ...form, status: 'inactive' })}
                      />
                      inactive
                    </label>
                  </div>
                </div>
              )}
            </div>
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={() => setModal({ open: false, product: null })}
                className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                disabled={saving}
                className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                {saving ? 'Saving...' : 'Save'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
