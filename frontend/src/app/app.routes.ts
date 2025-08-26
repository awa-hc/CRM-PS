import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'auth/login',
    pathMatch: 'full'
  },
  {
    path: 'auth',
    loadChildren: () => import('./features/auth/auth.routes').then(m => m.authRoutes)
  },
  {
    path: 'app',
    loadComponent: () => import('./shared/layout/main-layout.component').then(m => m.MainLayoutComponent),
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () => import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent)
      },
      {
        path: 'clients',
        loadChildren: () => import('./features/clients/clients.routes').then(m => m.clientsRoutes)
      },
      {
        path: 'projects',
        loadChildren: () => import('./features/projects/projects.routes').then(m => m.projectsRoutes)
      },
      {
        path: 'quotes',
        loadChildren: () => import('./features/quotes/quotes.routes').then(m => m.quotesRoutes)
      },
      {
        path: 'reports',
        loadChildren: () => import('./features/reports/reports.routes').then(m => m.reportsRoutes)
      },
      {
        path: '',
        redirectTo: 'dashboard',
        pathMatch: 'full'
      }
    ]
  },
  // {
  //   path: '**',
  //   redirectTo: 'auth/login'
  // }
];
