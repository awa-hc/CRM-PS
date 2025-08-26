import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, BehaviorSubject, tap, catchError, throwError } from 'rxjs';
import { User, LoginRequest, RegisterRequest, AuthResponse, ChangePasswordRequest, UpdateProfileRequest } from '../models/user.model';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly apiUrl = environment.apiUrl;
  
  private readonly currentUserSubject = new BehaviorSubject<User | null>(null);
  public readonly currentUser$ = this.currentUserSubject.asObservable();
  
  public readonly isAuthenticated = signal<boolean>(false);
  public readonly currentUser = signal<User | null>(null);

  constructor() {
    console.log('🔐 AuthService: Constructor called - initializing service');
    this.loadUserFromStorage();
  }

  private loadUserFromStorage(): void {
    console.log('🔐 AuthService: Loading user from storage');
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    console.log('🔐 AuthService: Token exists:', !!token);
    console.log('🔐 AuthService: User data exists:', !!userData);
    
    if (token && userData) {
      try {
        const user = JSON.parse(userData) as User;
        console.log('🔐 AuthService: Parsed user from storage:', user.email);
        this.currentUser.set(user);
        this.currentUserSubject.next(user);
        this.isAuthenticated.set(true);
        console.log('🔐 AuthService: Auth state set from storage - isAuthenticated:', this.isAuthenticated());
      } catch (error) {
        console.log('🔐 AuthService: Error parsing user data from storage:', error);
        this.logout();
      }
    } else {
      console.log('🔐 AuthService: No valid token/user data in storage');
    }
  }

  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/auth/login`, credentials)
      .pipe(
        tap(response => {
          this.setAuthData(response.token, response.user);
        }),
        catchError(error => {
          console.error('Login error:', error);
          return throwError(() => error);
        })
      );
  }

  register(userData: RegisterRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/auth/register`, userData)
      .pipe(
        tap(response => {
          this.setAuthData(response.token, response.user);
        }),
        catchError(error => {
          console.error('Registration error:', error);
          return throwError(() => error);
        })
      );
  }

  logout(): void {
    console.log('🔐 AuthService: Logging out user');
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    this.currentUser.set(null);
    this.currentUserSubject.next(null);
    this.isAuthenticated.set(false);
    console.log('🔐 AuthService: Auth state cleared - isAuthenticated:', this.isAuthenticated());
    this.router.navigate(['/auth/login']);
  }

  getProfile(): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}/auth/profile`)
      .pipe(
        tap(user => {
          this.currentUser.set(user);
          this.currentUserSubject.next(user);
          localStorage.setItem('user', JSON.stringify(user));
        }),
        catchError(error => {
          console.error('Get profile error:', error);
          if (error.status === 401) {
            this.logout();
          }
          return throwError(() => error);
        })
      );
  }

  updateProfile(profileData: UpdateProfileRequest): Observable<User> {
    return this.http.put<User>(`${this.apiUrl}/auth/profile`, profileData)
      .pipe(
        tap(user => {
          this.currentUser.set(user);
          this.currentUserSubject.next(user);
          localStorage.setItem('user', JSON.stringify(user));
        }),
        catchError(error => {
          console.error('Update profile error:', error);
          return throwError(() => error);
        })
      );
  }

  changePassword(passwordData: ChangePasswordRequest): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.apiUrl}/auth/change-password`, passwordData)
      .pipe(
        catchError(error => {
          console.error('Change password error:', error);
          return throwError(() => error);
        })
      );
  }

  verifyToken(): Observable<{ valid: boolean; user?: User }> {
    return this.http.get<{ valid: boolean; user?: User }>(`${this.apiUrl}/auth/verify`)
      .pipe(
        tap(response => {
          if (response.valid && response.user) {
            this.currentUser.set(response.user);
            this.currentUserSubject.next(response.user);
            this.isAuthenticated.set(true);
            localStorage.setItem('user', JSON.stringify(response.user));
          } else {
            this.logout();
          }
        }),
        catchError(error => {
          console.error('Token verification error:', error);
          this.logout();
          return throwError(() => error);
        })
      );
  }

  getToken(): string | null {
    return localStorage.getItem('token');
  }

  private setAuthData(token: string, user: User): void {
    console.log('🔐 AuthService: Setting auth data for user:', user.email);
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    this.currentUser.set(user);
    this.currentUserSubject.next(user);
    this.isAuthenticated.set(true);
    console.log('🔐 AuthService: Auth state updated - isAuthenticated:', this.isAuthenticated());
  }

  hasRole(role: string): boolean {
    const user = this.currentUser();
    return user ? user.role === role : false;
  }

  isAdmin(): boolean {
    return this.hasRole('admin');
  }

  isManager(): boolean {
    return this.hasRole('manager') || this.isAdmin();
  }
}