export interface Product {
  id: number;
  code: string;
  description: string;
  stock: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateProductInput {
  code: string;
  description: string;
  stock: number;
}