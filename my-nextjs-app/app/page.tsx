"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { get } from "@/lib/api";
import ProductCard from "@/components/shared/ProductCard";

interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  cover_image_url: string;
  file_key: string;
  created_at: string;
}

export default function HomePage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchProducts() {
      const { data, error: apiError } = await get<Product[]>("/products");
      if (apiError) {
        setError(apiError);
      } else if (data) {
        setProducts(data);
      }
      setLoading(false);
    }
    fetchProducts();
  }, []);

  if (loading) {
    return <div className="text-center py-8">Loading products...</div>;
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">Error: {error}</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold text-center mb-8">Our Products</h1>
      <div className="flex justify-center space-x-4 mb-8">
        <Link href="/login" className="text-blue-500 hover:underline">Login</Link>
        <Link href="/register" className="text-blue-500 hover:underline">Register</Link>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((product) => (
          <ProductCard
            key={product.id}
            id={product.id}
            name={product.name}
            price={product.price}
            coverImageURL={product.cover_image_url}
          />
        ))}
      </div>
      {products.length === 0 && (
        <p className="text-center text-gray-600 mt-8">No products available yet.</p>
      )}
    </div>
  );
}