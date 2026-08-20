import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { ProductService } from '../../core/services/product.service';
import { Product } from '../../core/models/product.model';

@Component({
  selector: 'app-product-list',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="filters">
      <input
        type="text"
        placeholder="Search by remedy, symptom, brand..."
        (input)="onSearch($event)"
      />
    </div>

    <div class="grid" *ngIf="!loading(); else loadingTpl">
      <a
        class="card"
        *ngFor="let p of products()"
        [routerLink]="['/products', p.slug]"
      >
        <h3>{{ p.name }}</h3>
        <p>{{ p.potency }} &middot; {{ p.form }} &middot; {{ p.manufacturer }}</p>
        <p class="price">₹{{ p.price / 100 | number: '1.2-2' }}</p>
        <span class="badge" *ngIf="p.schedule === 'H'">Prescription required</span>
      </a>
    </div>
    <ng-template #loadingTpl>Loading products...</ng-template>
  `
})
export class ProductListComponent implements OnInit {
  private productService = inject(ProductService);

  products = signal<Product[]>([]);
  loading = signal<boolean>(true);

  ngOnInit(): void {
    this.fetch();
  }

  fetch(query?: string): void {
    this.loading.set(true);
    this.productService.list({ q: query }).subscribe({
      next: (products) => {
        this.products.set(products);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
  }

  onSearch(event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.fetch(value);
  }
}
