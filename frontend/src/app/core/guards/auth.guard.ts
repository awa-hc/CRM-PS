import { CanActivateFn, Router } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { map, take, catchError } from 'rxjs/operators';
import { of } from 'rxjs';

export const authGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  
  console.log('🛡️ AuthGuard: Checking authentication for route:', state.url);
  
  const token = authService.getToken();
  console.log('🛡️ AuthGuard: Token exists:', !!token);
  
  if (!token) {
    console.log('🛡️ AuthGuard: No token found - redirecting to login');
    router.navigate(['/auth/login']);
    return false;
  }
  
  const isAuthenticated = authService.isAuthenticated();
  console.log('🛡️ AuthGuard: Local auth state:', isAuthenticated);
  
  // If user is authenticated locally and we have a token, allow access immediately
  // Only verify with backend periodically or when there's doubt
  if (isAuthenticated) {
    console.log('🛡️ AuthGuard: User authenticated locally with valid token - allowing access');
    return true;
  }
  
  // If not authenticated locally but we have a token, verify with backend
  console.log('🛡️ AuthGuard: Token exists but not authenticated locally - verifying with backend');
  return authService.verifyToken().pipe(
    map(response => {
      console.log('🛡️ AuthGuard: Backend verification result:', response.valid);
      if (response.valid) {
        console.log('🛡️ AuthGuard: Token verified - allowing access');
        return true;
      } else {
        console.log('🛡️ AuthGuard: Token invalid - redirecting to login');
        router.navigate(['/auth/login']);
        return false;
      }
    }),
    catchError(error => {
      console.log('🛡️ AuthGuard: Verification error:', error);
      router.navigate(['/auth/login']);
      return of(false);
    })
  );
};

export const adminGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // Check if user has a token
  const token = authService.getToken();
  if (!token) {
    router.navigate(['/auth/login']);
    return false;
  }

  // If user is authenticated and is admin, allow access
  if (authService.isAuthenticated() && authService.isAdmin()) {
    return true;
  }

  // Verify token with backend and check admin role
  return authService.verifyToken().pipe(
    take(1),
    map(response => {
      if (response.valid && response.user?.role === 'admin') {
        return true;
      } else {
        router.navigate(['/app/dashboard']);
        return false;
      }
    }),
    catchError(() => {
      router.navigate(['/auth/login']);
      return of(false);
    })
  );
};

export const managerGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // Check if user has a token
  const token = authService.getToken();
  if (!token) {
    router.navigate(['/auth/login']);
    return false;
  }

  // If user is authenticated and is manager/admin, allow access
  if (authService.isAuthenticated() && authService.isManager()) {
    return true;
  }

  // Verify token with backend and check manager/admin role
  return authService.verifyToken().pipe(
    take(1),
    map(response => {
      if (response.valid && response.user && 
          (response.user.role === 'admin' || response.user.role === 'manager')) {
        return true;
      } else {
        router.navigate(['/app/dashboard']);
        return false;
      }
    }),
    catchError(() => {
      router.navigate(['/auth/login']);
      return of(false);
    })
  );
};