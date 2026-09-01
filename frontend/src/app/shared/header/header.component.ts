import { CommonModule } from "@angular/common";
import { Component, inject } from "@angular/core";
import { Router, RouterLink } from "@angular/router";
import { AuthService } from "../../core/services/auth.service";
import { TitleCasePipe } from "@angular/common";

@Component({
  selector: "app-header",
  standalone: true,
  imports: [CommonModule, RouterLink, TitleCasePipe],
  templateUrl: "./header.component.html",
  styleUrl: "./header.component.css",
})
export class HeaderComponent {
  auth = inject(AuthService);
  private router = inject(Router);

  // Default cart and favorites count (to be connected to services/API in future)
  cartCount = 0;
  favoritesCount = 0;

  logout(): void {
    this.auth.logout();
    this.router.navigate(["/"]);
  }
}
