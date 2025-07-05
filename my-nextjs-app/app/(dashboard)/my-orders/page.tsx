'use client';

import { useEffect, useState } from 'react';
import { get } from '@/lib/api';
import Link from 'next/link';

interface Order {
  id: number;
  user_id: number;
  product_id: number;
  amount: number;
  status: string;
  created_at: string;
  paid_at: string | null;
  product_name?: string; // Assuming backend can return product name with order
}

export default function MyOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchOrders() {
      const { data, error: apiError } = await get<Order[]>('/orders');
      if (apiError) {
        setError(apiError);
      } else if (data) {
        setOrders(data);
      }
      setLoading(false);
    }
    fetchOrders();
  }, []);

  if (loading) {
    return <div className="text-center py-8">Loading orders...</div>;
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">Error: {error}</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-center mb-8">My Orders</h1>
      {orders.length === 0 ? (
        <p className="text-center text-gray-600">You have no orders yet.</p>
      ) : (
        <div className="grid grid-cols-1 gap-6">
          {orders.map((order) => (
            <div key={order.id} className="bg-white p-6 rounded-lg shadow-md">
              <h2 className="text-xl font-semibold mb-2">Order ID: {order.id}</h2>
              <p className="text-gray-700">Product: {order.product_name || 'N/A'}</p>
              <p className="text-gray-700">Amount: ${order.amount.toFixed(2)}</p>
              <p className="text-gray-700">Status: {order.status}</p>
              <p className="text-gray-700">Order Date: {new Date(order.created_at).toLocaleDateString()}</p>
              {order.paid_at && (
                <p className="text-gray-700">Paid Date: {new Date(order.paid_at).toLocaleDateString()}</p>
              )}
              <Link href={`/my-orders/${order.id}`} className="text-blue-500 hover:underline mt-4 inline-block">
                View Details
              </Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
