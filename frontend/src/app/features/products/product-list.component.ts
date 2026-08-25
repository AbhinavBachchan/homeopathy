import { Component, OnInit, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { RouterLink } from "@angular/router";
import { ProductService } from "../../core/services/product.service";
import { Product } from "../../core/models/product.model";

@Component({
  selector: "app-product-list",
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: "./product-list.component.html",
  styleUrl: "./product-list.component.css",
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
      error: () => this.loading.set(false),
    });
  }

  onSearch(event: Event): void {
    this.fetch((event.target as HTMLInputElement).value);
  }
}
