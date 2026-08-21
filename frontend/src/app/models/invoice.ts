import { InvoiceStatus } from '../enums/invoice-status';
import { CreateInvoiceItemInput, InvoiceItem } from './invoice-item';

export interface Invoice {
  id: number;
  number: number;
  status: InvoiceStatus;
  items: InvoiceItem[] | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateInvoiceInput {
  items: CreateInvoiceItemInput[];
}