import { Routes } from "@angular/router";

export const routes: Routes = [
  {
    path: "",
    loadComponent: () =>
      import("./features/home/home.component").then((m) => m.HomeComponent),
  },
  {
    path: "products",
    loadComponent: () =>
      import("./features/products/product-list.component").then(
        (m) => m.ProductListComponent,
      ),
  },
  {
    path: "products/:slug",
    loadComponent: () =>
      import("./features/products/product-detail.component").then(
        (m) => m.ProductDetailComponent,
      ),
  },
  {
    path: "login",
    loadComponent: () =>
      import("./features/auth/login.component").then((m) => m.LoginComponent),
  },
  {
    path: "register",
    loadComponent: () =>
      import("./features/auth/register.component").then(
        (m) => m.RegisterComponent,
      ),
  },
  // TODO: /cart, /checkout, /account/orders, /consultations (Phase 2),
  // /doctor/dashboard, /admin (guarded by role via a functional CanActivate
  // guard reading AuthService.currentUser()).
];
