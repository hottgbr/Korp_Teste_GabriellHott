export interface InvoiceItem {
  id: number;
  invoiceId: number;
  productId: number;
  quantity: number;
}

export interface CreateInvoiceItemInput {
  productId: number;
  quantity: number;
}