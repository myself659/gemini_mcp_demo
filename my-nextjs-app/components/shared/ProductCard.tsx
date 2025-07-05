import Link from "next/link";

interface ProductCardProps {
  id: number;
  name: string;
  price: number;
  coverImageURL?: string;
}

export default function ProductCard({ id, name, price, coverImageURL }: ProductCardProps) {
  return (
    <Link href={`/products/${id}`} className="block border rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200">
      <div className="relative w-full h-48 bg-gray-200 rounded-t-lg overflow-hidden">
        {coverImageURL ? (
          <img src={coverImageURL} alt={name} className="w-full h-full object-cover" />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-500">
            No Image
          </div>
        )}
      </div>
      <div className="p-4">
        <h3 className="text-lg font-semibold text-gray-800 truncate">{name}</h3>
        <p className="mt-1 text-gray-600">${price.toFixed(2)}</p>
      </div>
    </Link>
  );
}
