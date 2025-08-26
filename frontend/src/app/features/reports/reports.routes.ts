import { Routes } from '@angular/router';
import { adminGuard, managerGuard } from '../../core/guards/auth.guard';

export const reportsRoutes: Routes = [
  {
    path: '',
    loadComponent: () => import('./reports-dashboard/reports-dashboard.component').then(m => m.ReportsDashboardComponent),
    canActivate: [managerGuard]
  },
  {
    path: 'financial',
    loadComponent: () => import('./financial-reports/financial-reports.component').then(m => m.FinancialReportsComponent),
    canActivate: [managerGuard]
  },
  {
    path: 'projects',
    loadComponent: () => import('./project-reports/project-reports.component').then(m => m.ProjectReportsComponent),
    canActivate: [managerGuard]
  },
  {
    path: 'clients',
    loadComponent: () => import('./client-reports/client-reports.component').then(m => m.ClientReportsComponent),
    canActivate: [managerGuard]
  }
];