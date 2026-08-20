import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()">
      <label>
        Email
        <input type="email" formControlName="email" />
      </label>
      <label>
        Password
        <input type="password" formControlName="password" />
      </label>
      <button type="submit" [disabled]="form.invalid || submitting()">
        Log in
      </button>
      <p class="error" *ngIf="error()">{{ error() }}</p>
    </form>
  `
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  submitting = signal(false);
  error = signal<string | null>(null);

  form = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required]]
  });

  onSubmit(): void {
    if (this.form.invalid) return;
    this.submitting.set(true);
    this.error.set(null);

    this.authService
      .login({
        email: this.form.value.email!,
        password: this.form.value.password!
      })
      .subscribe({
        next: () => {
          this.submitting.set(false);
          this.router.navigate(['/']);
        },
        error: () => {
          this.submitting.set(false);
          this.error.set('Invalid email or password');
        }
      });
  }
}
