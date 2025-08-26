import { Routes } from '@angular/router';

export const quotesRoutes: Routes = [
  {
    path: '',
    loadComponent: () => import('./quotes-list/quotes-list.component').then(m => m.QuotesListComponent)
  },
  {
    path: 'new',
    loadComponent: () => import('./quote-form/quote-form.component').then(m => m.QuoteFormComponent)
  },
  {
    path: ':id',
    loadComponent: () => import('./quote-detail/quote-detail.component').then(m => m.QuoteDetailComponent)
  },
  {
    path: ':id/edit',
    loadComponent: () => import('./quote-form/quote-form.component').then(m => m.QuoteFormComponent)
  }
];