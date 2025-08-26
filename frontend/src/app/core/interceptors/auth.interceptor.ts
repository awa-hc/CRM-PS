import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { catchError, throwError } from 'rxjs';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  
  console.log('🌐 Interceptor: Processing request to:', req.url);
  
  // Skip adding token for login and register endpoints
  if (req.url.includes('/auth/login') || req.url.includes('/auth/register')) {
    console.log('🌐 Interceptor: Skipping token for auth endpoint');
    return next(req);
  }
  
  const token = authService.getToken();
  console.log('🌐 Interceptor: Token exists:', !!token);
  
  if (token) {
    const authReq = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
    
    console.log('🌐 Interceptor: Adding Authorization header');
    
    return next(authReq).pipe(
      catchError(error => {
        console.log('🌐 Interceptor: Request error:', error.status, error.message);
        if (error.status === 401) {
          console.log('🌐 Interceptor: 401 error detected');
          // Only logout if this is not a verify token request
          // Let the AuthGuard handle token verification failures
          if (!req.url.includes('/auth/verify')) {
            console.log('🌐 Interceptor: 401 on non-verify endpoint - calling logout');
            authService.logout();
            router.navigate(['/auth/login']);
          } else {
            console.log('🌐 Interceptor: 401 on verify endpoint - letting AuthGuard handle it');
          }
        }
        return throwError(() => error);
      })
    );
  }
  
  console.log('🌐 Interceptor: No token, proceeding without Authorization header');
  return next(req);
};