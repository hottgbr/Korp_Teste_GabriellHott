import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { CreateProductInput, Product } from '../models/product';

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  private readonly http = inject(HttpClient);

  private readonly apiUrl = 'http://localhost:8081/products';

  list(): Observable<Product[]> {
    return this.http.get<Product[]>(this.apiUrl);
  }

  create(input: CreateProductInput): Observable<Product> {
    return this.http.post<Product>(
      this.apiUrl,
      input,
    );
  }

  findById(id: number): Observable<Product> {
    return this.http.get<Product>(
      `${this.apiUrl}/${id}`,
    );
  }
}