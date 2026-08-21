import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import {
  CreateInvoiceInput,
  Invoice,
} from '../models/invoice';

@Injectable({
  providedIn: 'root',
})
export class InvoiceService {
  private readonly http = inject(HttpClient);

  private readonly apiUrl =
    'http://localhost:8082/invoices';

  list(): Observable<Invoice[]> {
    return this.http.get<Invoice[]>(
      this.apiUrl,
    );
  }

  findById(id: number): Observable<Invoice> {
    return this.http.get<Invoice>(
      `${this.apiUrl}/${id}`,
    );
  }

  create(
    input: CreateInvoiceInput,
  ): Observable<Invoice> {
    return this.http.post<Invoice>(
      this.apiUrl,
      input,
    );
  }

  close(id: number): Observable<Invoice> {
    return this.http.post<Invoice>(
      `${this.apiUrl}/${id}/close`,
      {},
    );
  }
}