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
  {
    path: "forgot-password",
    loadComponent: () =>
      import("./features/auth/forgot-password.component").then(
        (m) => m.ForgotPasswordComponent,
      ),
  },
  {
    path: "reset-password",
    loadComponent: () =>
      import("./features/auth/forgot-password.component").then(
        (m) => m.ForgotPasswordComponent,
      ),
  },
  {
    path: "cart",
    loadComponent: () =>
      import("./features/cart/cart.component").then((m) => m.CartComponent),
  },
  {
    path: "favorites",
    loadComponent: () =>
      import("./features/favorites/favorites.component").then((m) => m.FavoritesComponent),
  },
  // TODO: /checkout, /account/orders, /consultations (Phase 2),
  // /doctor/dashboard, /admin (guarded by role via a functional CanActivate
  // guard reading AuthService.currentUser()).
];
