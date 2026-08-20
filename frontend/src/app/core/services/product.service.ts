import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { Product, ProductFilters } from '../models/product.model';

interface ApiEnvelope<T> {
  success: boolean;
  data: T;
}

@Injectable({ providedIn: 'root' })
export class ProductService {
  private http = inject(HttpClient);
  private baseUrl = `${environment.apiUrl}/products`;

  list(filters: ProductFilters = {}): Observable<Product[]> {
    let params = new HttpParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value) params = params.set(key, value);
    });
    return this.http
      .get<ApiEnvelope<Product[]>>(this.baseUrl, { params })
      .pipe(map((res) => res.data));
  }

  getBySlug(slug: string): Observable<Product> {
    return this.http
      .get<ApiEnvelope<Product>>(`${this.baseUrl}/${slug}`)
      .pipe(map((res) => res.data));
  }
}
