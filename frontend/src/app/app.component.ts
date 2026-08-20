import { Component } from '@angular/core';
import { RouterOutlet, RouterLink } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterLink],
  template: `
    <header>
      <a routerLink="/">Homeopathy Store</a>
      <nav>
        <a routerLink="/products">Shop</a>
        <a routerLink="/login">Login</a>
        <a routerLink="/cart">Cart</a>
      </nav>
    </header>
    <main>
      <router-outlet></router-outlet>
    </main>
  `
})
export class AppComponent {}
