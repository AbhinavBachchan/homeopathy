export interface Product {
  id: string;
  name: string;
  slug: string;
  potency: string;
  form: string;
  size_quantity: string;
  manufacturer: string;
  therapeutic_category: string;
  indications: string[];
  contraindications: string[];
  schedule: 'OTC' | 'H';
  hsn_code: string;
  sku: string;
  price: number; // paise
  mrp: number;
  stock_qty: number;
  images: string[];
  is_active: boolean;
}

export interface ProductFilters {
  potency?: string;
  brand?: string;
  category?: string;
  q?: string;
}
