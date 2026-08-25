import { Component, OnInit, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { ActivatedRoute } from "@angular/router";
import { ProductService } from "../../core/services/product.service";
import { Product } from "../../core/models/product.model";

@Component({
  selector: "app-product-detail",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./product-detail.component.html",
  styleUrl: "./product-detail.component.css",
})
export class ProductDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private productService = inject(ProductService);
  product = signal<Product | null>(null);

  ngOnInit(): void {
    const slug = this.route.snapshot.paramMap.get("slug");
    if (slug)
      this.productService.getBySlug(slug).subscribe((p) => this.product.set(p));
  }
}
