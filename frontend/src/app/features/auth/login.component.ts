import { Component, inject, signal } from "@angular/core";
import { HttpErrorResponse } from "@angular/common/http";
import { CommonModule } from "@angular/common";
import { ReactiveFormsModule, FormBuilder, Validators } from "@angular/forms";
import { Router, RouterLink } from "@angular/router";
import { AuthService } from "../../core/services/auth.service";
import { TitleCasePipe } from "@angular/common";

@Component({
  selector: "app-login",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TitleCasePipe, RouterLink],
  templateUrl: "./login.component.html",
  styleUrl: "./login.component.css",
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  submitting = signal(false);
  error = signal<string | null>(null);
  form = this.fb.group({
    email: ["", [Validators.required, Validators.email]],
    password: ["", [Validators.required]],
  });

  onSubmit(): void {
    if (this.form.invalid) return;
    this.submitting.set(true);
    this.error.set(null);
    this.authService
      .login({
        email: this.form.value.email!,
        password: this.form.value.password!,
      })
      .subscribe({
        next: () => {
          this.submitting.set(false);
          this.router.navigate(["/"]);
        },
        error: (err: HttpErrorResponse) => {
          this.submitting.set(false);
          this.error.set(this.getApiErrorMessage(err) ?? "Invalid email or password");
        },
      });
  }

  private getApiErrorMessage(err: HttpErrorResponse): string | null {
    const body = err.error;

    if (typeof body === "string" && body.trim()) {
      try {
        const parsed = JSON.parse(body) as { error?: string };
        return parsed.error ?? body;
      } catch {
        return body;
      }
    }

    if (body && typeof body === "object" && "error" in body) {
      const message = (body as { error?: unknown }).error;
      return typeof message === "string" && message.trim() ? message : null;
    }

    return null;
  }
}
