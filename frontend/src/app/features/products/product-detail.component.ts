import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { ProductService } from '../../core/services/product.service';
import { Product } from '../../core/models/product.model';

@Component({
  selector: 'app-product-detail',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div *ngIf="product() as p">
      <h1>{{ p.name }} ({{ p.potency }})</h1>
      <p>{{ p.form }} &middot; {{ p.size_quantity }} &middot; {{ p.manufacturer }}</p>
      <p class="price">₹{{ p.price / 100 | number: '1.2-2' }}</p>

      <h4>Indications</h4>
      <ul><li *ngFor="let i of p.indications">{{ i }}</li></ul>

      <h4>Contraindications</h4>
      <ul><li *ngFor="let c of p.contraindications">{{ c }}</li></ul>

      <div class="rx-notice" *ngIf="p.schedule === 'H'">
        This is a Schedule H medicine. You'll need to upload a valid
        prescription before checkout.
      </div>

      <button [disabled]="p.stock_qty === 0">
        {{ p.stock_qty > 0 ? 'Add to cart' : 'Out of stock' }}
      </button>
    </div>
  `
})
export class ProductDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private productService = inject(ProductService);

  product = signal<Product | null>(null);

  ngOnInit(): void {
    const slug = this.route.snapshot.paramMap.get('slug');
    if (slug) {
      this.productService.getBySlug(slug).subscribe((p) => this.product.set(p));
    }
  }
}
