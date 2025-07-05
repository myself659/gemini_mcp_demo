'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { get } from '@/lib/api';

interface Order {
  id: number;
  user_id: number;
  product_id: number;
  amount: number;
  status: string;
  created_at: string;
  paid_at: string | null;
  product_name?: string; // Assuming backend can return product name with order
  product_description?: string;
  product_file_key?: string;
}

export default function OrderDetailPage() {
  const { id } = useParams();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchOrder() {
      if (!id) return;

      const { data, error: apiError } = await get<Order>(`/orders/${id}`);
      if (apiError) {
        setError(apiError);
      } else if (data) {
        setOrder(data);
      }
      setLoading(false);
    }
    fetchOrder();
  }, [id]);

  const handleDownload = async () => {
    if (!order) return;
    try {
      const { data, error: apiError } = await get<{ download_url: string }>(`/downloads/order/${order.id}`);
      if (apiError) {
        alert(`Error generating download link: ${apiError}`);
      } else if (data && data.download_url) {
        window.location.href = data.download_url;
      }
    } catch (err) {
      alert('Failed to initiate download.');
      console.error(err);
    }
  };

  if (loading) {
    return <div className="text-center py-8">Loading order details...</div>;
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">Error: {error}</div>;
  }

  if (!order) {
    return <div className="text-center py-8">Order not found.</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-center mb-8">Order Details</h1>
      <div className="bg-white p-8 rounded-lg shadow-md max-w-2xl mx-auto">
        <p className="text-lg mb-2"><strong>Order ID:</strong> {order.id}</p>
        <p className="text-lg mb-2"><strong>Product:</strong> {order.product_name || 'N/A'}</p>
        <p className="text-lg mb-2"><strong>Description:</strong> {order.product_description || 'N/A'}</p>
        <p className="text-lg mb-2"><strong>Amount:</strong> ${order.amount.toFixed(2)}</p>
        <p className="text-lg mb-2"><strong>Status:</strong> {order.status}</p>
        <p className="text-lg mb-2"><strong>Order Date:</strong> {new Date(order.created_at).toLocaleString()}</p>
        {order.paid_at && (
          <p className="text-lg mb-2"><strong>Paid Date:</strong> {new Date(order.paid_at).toLocaleString()}</p>
        )}
        {order.status === 'completed' && order.product_file_key && (
          <button
            onClick={handleDownload}
            className="mt-6 bg-green-500 hover:bg-green-600 text-white font-bold py-3 px-6 rounded-lg shadow-lg transition duration-300 ease-in-out"
          >
            Download Product
          </button>
        )}
      </div>
    </div>
  );
}
