import { Component, EventEmitter, Input, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-favorite-button',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './favorite-button.component.html',
  styleUrl: './favorite-button.component.css',
})
export class FavoriteButtonComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  @Input({ required: true }) productId!: string;
  @Input() isFavorite: boolean = false;
  @Output() favoriteChange = new EventEmitter<{ productId: string; isFavorite: boolean }>();

  toggleFavorite(event: Event): void {
    event.stopPropagation();
    event.preventDefault();

    // If user is not logged in, redirect to login page
    if (!this.authService.isLoggedIn()) {
      this.router.navigate(['/login'], {
        queryParams: { returnUrl: this.router.url },
      });
      return;
    }

    // Toggle favorite state
    this.isFavorite = !this.isFavorite;
    this.favoriteChange.emit({
      productId: this.productId,
      isFavorite: this.isFavorite,
    });
  }
}
