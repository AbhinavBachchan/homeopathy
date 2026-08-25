import { CommonModule } from "@angular/common";
import { Component, inject, signal } from "@angular/core";
import { HttpErrorResponse } from "@angular/common/http";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { Router } from "@angular/router";
import { AuthService } from "../../core/services/auth.service";
import { TitleCasePipe } from "@angular/common";

@Component({
  selector: "app-register",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TitleCasePipe],
  templateUrl: "./register.component.html",
  styleUrl: "./register.component.css",
})
export class RegisterComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  submitting = signal(false);
  error = signal<string | null>(null);

  form = this.fb.group({
    name: ["", [Validators.required]],
    email: ["", [Validators.required, Validators.email]],
    phone: [""],
    password: ["", [Validators.required, Validators.minLength(6)]],
  });

  onSubmit(): void {
    if (this.form.invalid) return;

    this.submitting.set(true);
    this.error.set(null);

    this.authService
      .register({
        name: this.form.value.name!,
        email: this.form.value.email!,
        phone: this.form.value.phone || undefined,
        password: this.form.value.password!,
      })
      .subscribe({
        next: () => {
          this.submitting.set(false);
          this.router.navigate(["/"]);
        },
        error: (err: HttpErrorResponse) => {
          this.submitting.set(false);
          this.error.set(this.getApiErrorMessage(err) ?? "Could not create account");
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
