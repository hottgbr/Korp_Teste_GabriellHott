import { Routes } from '@angular/router';

import { Products } from './pages/products/products';
import { ProductCreate } from './pages/product-create/product-create';
import { Invoices } from './pages/invoices/invoices';
import { InvoiceCreate } from './pages/invoice-create/invoice-create';
import { InvoiceDetails } from './pages/invoice-details/invoice-details';

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
  {
    path: 'invoices',
    component: Invoices,
  },
  {
    path: 'invoices/new',
    component: InvoiceCreate,
  },
  {
    path: 'invoices/:id',
    component: InvoiceDetails,
  },
];