import { Component, OnInit, inject, signal } from "@angular/core";
import { HttpErrorResponse } from "@angular/common/http";
import { CommonModule, TitleCasePipe } from "@angular/common";
import {
  ReactiveFormsModule,
  FormBuilder,
  Validators,
  AbstractControl,
  ValidationErrors,
} from "@angular/forms";
import { ActivatedRoute, Router, RouterLink } from "@angular/router";
import { AuthService } from "../../core/services/auth.service";

function passwordMatchValidator(
  control: AbstractControl,
): ValidationErrors | null {
  const password = control.get("password")?.value;
  const confirmPassword = control.get("confirmPassword")?.value;
  if (password && confirmPassword && password !== confirmPassword) {
    return { passwordMismatch: true };
  }
  return null;
}

@Component({
  selector: "app-forgot-password",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink, TitleCasePipe],
  templateUrl: "./forgot-password.component.html",
  styleUrl: "./forgot-password.component.css",
})
export class ForgotPasswordComponent implements OnInit {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  step = signal<"request" | "reset" | "completed">("request");
  submitting = signal(false);
  error = signal<string | null>(null);
  successMessage = signal<string | null>(null);

  requestForm = this.fb.group({
    email: ["", [Validators.required, Validators.email]],
  });

  resetForm = this.fb.group(
    {
      email: ["", [Validators.required, Validators.email]],
      token: ["", [Validators.required]],
      password: ["", [Validators.required, Validators.minLength(8)]],
      confirmPassword: ["", [Validators.required]],
    },
    { validators: passwordMatchValidator },
  );

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      const token = params["token"];
      const email = params["email"];
      if (email) {
        this.requestForm.patchValue({ email });
        this.resetForm.patchValue({ email });
      }
      if (token) {
        this.resetForm.patchValue({ token });
        this.step.set("reset");
      }
    });
  }

  onRequestReset(): void {
    if (this.requestForm.invalid) return;
    this.submitting.set(true);
    this.error.set(null);
    this.successMessage.set(null);

    const email = this.requestForm.value.email!;

    this.authService.forgotPassword({ email }).subscribe({
      next: (res) => {
        this.submitting.set(false);
        this.resetForm.patchValue({ email });
        this.successMessage.set(
          res.data?.message ??
            "A 6 digit verification code has been sent to the registered email. Kindly enter the code.",
        );
        this.step.set("reset");
      },
      error: (err: HttpErrorResponse) => {
        this.submitting.set(false);
        this.error.set(this.getApiErrorMessage(err) ?? "Email is not registered");
      },
    });
  }

  onResetPassword(): void {
    if (this.resetForm.invalid) return;
    this.submitting.set(true);
    this.error.set(null);
    this.successMessage.set(null);

    const { email, token, password } = this.resetForm.value;

    this.authService
      .resetPassword({
        email: email!,
        token: token!,
        password: password!,
      })
      .subscribe({
        next: (res) => {
          this.submitting.set(false);
          this.successMessage.set(
            res.data?.message ?? "Your password has been successfully reset.",
          );
          this.step.set("completed");
        },
        error: (err: HttpErrorResponse) => {
          this.submitting.set(false);
          this.error.set(this.getApiErrorMessage(err) ?? "Failed to reset password");
        },
      });
  }

  goToResetStep(): void {
    this.error.set(null);
    this.step.set("reset");
  }

  goToRequestStep(): void {
    this.error.set(null);
    this.step.set("request");
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
