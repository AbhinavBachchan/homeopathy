import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthResponse, User } from '../models/user.model';

const TOKEN_KEY = 'hp_token';
const USER_KEY = 'hp_user';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private http = inject(HttpClient);
  private baseUrl = `${environment.apiUrl}/auth`;

  // Signal-based current-user state, readable from any standalone component.
  currentUser = signal<User | null>(this.getStoredUser());

  register(payload: { email: string; password: string; name: string; phone?: string }): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.baseUrl}/register`, payload).pipe(
      tap((res) => this.persistSession(res))
    );
  }

  login(payload: { email: string; password: string }): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.baseUrl}/login`, payload).pipe(
      tap((res) => this.persistSession(res))
    );
  }

  logout(): void {
    sessionStorage.removeItem(TOKEN_KEY);
    sessionStorage.removeItem(USER_KEY);
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    this.currentUser.set(null);
  }

  getToken(): string | null {
    const sessionToken = sessionStorage.getItem(TOKEN_KEY);
    if (sessionToken) return sessionToken;

    const legacyToken = localStorage.getItem(TOKEN_KEY);
    if (legacyToken) {
      sessionStorage.setItem(TOKEN_KEY, legacyToken);
      localStorage.removeItem(TOKEN_KEY);
    }

    return legacyToken;
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  private persistSession(res: AuthResponse): void {
    sessionStorage.setItem(TOKEN_KEY, res.data.token);
    sessionStorage.setItem(USER_KEY, JSON.stringify(res.data.user));
    this.currentUser.set(res.data.user);
  }

  private getStoredUser(): User | null {
    const storedUser = sessionStorage.getItem(USER_KEY);
    if (!storedUser) return null;

    try {
      return JSON.parse(storedUser) as User;
    } catch {
      sessionStorage.removeItem(USER_KEY);
      return null;
    }
  }

  forgotPassword(payload: { email: string }): Observable<{ success: boolean; data: { message: string } }> {
    return this.http.post<{ success: boolean; data: { message: string } }>(
      `${this.baseUrl}/forgot-password`,
      payload
    );
  }

  resetPassword(payload: { email: string; token: string; password: string }): Observable<{ success: boolean; data: { message: string } }> {
    return this.http.post<{ success: boolean; data: { message: string } }>(
      `${this.baseUrl}/reset-password`,
      payload
    );
  }

  // TODO: loginWithOtp() hitting MSG91-backed /auth/otp/* endpoints, and
  // Google OAuth redirect flow, once those backend handlers exist.
}
