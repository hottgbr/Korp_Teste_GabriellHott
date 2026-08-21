import { Routes } from '@angular/router';

import { Products } from './pages/products/products';
import { ProductCreate } from './pages/product-create/product-create';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'products',
    pathMatch: 'full',
  },
  {
    path: 'products',
    component: Products,
  },
  {
    path: 'products/new',
    component: ProductCreate,
  },
];