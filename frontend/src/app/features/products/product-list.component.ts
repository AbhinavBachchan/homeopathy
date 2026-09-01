import { Component, OnInit, computed, inject, signal } from "@angular/core";
import { CommonModule } from "@angular/common";
import { Router, RouterLink } from "@angular/router";
import { ProductService } from "../../core/services/product.service";
import { AuthService } from "../../core/services/auth.service";
import { Product } from "../../core/models/product.model";
import { MOCK_PRODUCTS } from "../../core/mock/mock-products";
import { ActiveFilters, ProductFilterComponent } from "./product-filter/product-filter.component";
import { FavoriteButtonComponent } from "../../shared/favorite-button/favorite-button.component";
import { PaginationComponent } from "../../shared/pagination/pagination.component";

@Component({
  selector: "app-product-list",
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    ProductFilterComponent,
    FavoriteButtonComponent,
    PaginationComponent,
  ],
  templateUrl: "./product-list.component.html",
  styleUrl: "./product-list.component.css",
})
export class ProductListComponent implements OnInit {
  private productService = inject(ProductService);
  private authService = inject(AuthService);
  private router = inject(Router);

  allProducts = signal<Product[]>(MOCK_PRODUCTS);
  searchQuery = signal<string>("");
  activeFilters = signal<ActiveFilters>({});
  loading = signal<boolean>(false);

  // Favorites tracking (Set of product IDs)
  favorites = signal<Set<string>>(new Set<string>());

  // Added to cart feedback tracking (Set of product IDs)
  addedToCart = signal<Set<string>>(new Set<string>());

  // Max 9 items per page as requested
  pageSize = signal<number>(9);
  currentPage = signal<number>(1);

  // Filtered products based on search query and side-filter checkboxes
  displayedProducts = computed(() => {
    let result = this.allProducts();
    const q = this.searchQuery().toLowerCase().trim();
    const filters = this.activeFilters();

    // 1. Filter by text search query (e.g., searching "abc" returns matching ABC products)
    if (q) {
      result = result.filter(
        (p) =>
          p.name?.toLowerCase().includes(q) ||
          p.manufacturer?.toLowerCase().includes(q) ||
          p.potency?.toLowerCase().includes(q) ||
          p.therapeutic_category?.toLowerCase().includes(q) ||
          p.form?.toLowerCase().includes(q) ||
          p.indications?.some((ind) => ind.toLowerCase().includes(q))
      );
    }

    // 2. Filter by homeopathy / category checkboxes
    const homeopathyFilters = filters['homeopathy'] || filters['category'];
    if (homeopathyFilters && homeopathyFilters.length > 0) {
      result = result.filter((p) =>
        homeopathyFilters.some((f) =>
          p.form?.toLowerCase().includes(f.toLowerCase()) ||
          p.therapeutic_category?.toLowerCase().includes(f.toLowerCase()) ||
          p.name?.toLowerCase().includes(f.toLowerCase())
        )
      );
    }

    // 3. Filter by potency checkboxes
    const potencyFilters = filters['potency'];
    if (potencyFilters && potencyFilters.length > 0) {
      result = result.filter((p) =>
        potencyFilters.some((pot) =>
          p.potency?.toLowerCase().includes(pot.toLowerCase())
        )
      );
    }

    // 4. Filter by brand/manufacturer checkboxes
    const brandFilters = filters['brand'];
    if (brandFilters && brandFilters.length > 0) {
      result = result.filter((p) =>
        brandFilters.some((b) =>
          p.manufacturer?.toLowerCase().includes(b.toLowerCase())
        )
      );
    }

    return result;
  });

  // Current page products slice (max 9 items)
  paginatedProducts = computed(() => {
    const products = this.displayedProducts();
    const page = this.currentPage();
    const size = this.pageSize();
    const start = (page - 1) * size;
    return products.slice(start, start + size);
  });

  ngOnInit(): void {
    this.fetch();
  }

  fetch(query?: string): void {
    this.loading.set(true);
    this.productService.list({ q: query }).subscribe({
      next: (products) => {
        if (products && products.length > 0) {
          this.allProducts.set(products);
        } else {
          this.allProducts.set(MOCK_PRODUCTS);
        }
        this.loading.set(false);
      },
      error: () => {
        this.allProducts.set(MOCK_PRODUCTS);
        this.loading.set(false);
      },
    });
  }

  onSearch(event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.searchQuery.set(value);
    this.currentPage.set(1); // Reset to page 1 on search
  }

  onFilterChange(filters: ActiveFilters): void {
    this.activeFilters.set(filters);
    this.currentPage.set(1); // Reset to page 1 on filter change
  }

  // Handle favorite toggling from FavoriteButtonComponent
  onFavoriteChange(event: { productId: string; isFavorite: boolean }): void {
    const updated = new Set(this.favorites());
    if (event.isFavorite) {
      updated.add(event.productId);
    } else {
      updated.delete(event.productId);
    }
    this.favorites.set(updated);
  }

  isFavorite(productId: string): boolean {
    return this.favorites().has(productId);
  }

  // Add to cart with auth check
  addToCart(product: Product, event: Event): void {
    event.stopPropagation();
    event.preventDefault();

    // If not logged in, redirect to login page
    if (!this.authService.isLoggedIn()) {
      this.router.navigate(['/login'], {
        queryParams: { returnUrl: this.router.url },
      });
      return;
    }

    // Add to cart feedback
    const updated = new Set(this.addedToCart());
    updated.add(product.id);
    this.addedToCart.set(updated);

    setTimeout(() => {
      const resetSet = new Set(this.addedToCart());
      resetSet.delete(product.id);
      this.addedToCart.set(resetSet);
    }, 2000);
  }

  isAddedToCart(productId: string): boolean {
    return this.addedToCart().has(productId);
  }

  // Page change from PaginationComponent
  onPageChange(newPage: number): void {
    this.currentPage.set(newPage);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}
