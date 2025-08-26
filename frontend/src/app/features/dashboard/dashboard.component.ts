import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatListModule } from '@angular/material/list';
import { MatChipsModule } from '@angular/material/chips';
import { Router } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { DashboardService, DashboardStats, RecentActivity } from '../../core/services/dashboard.service';
import { User } from '../../core/models/user.model';
import { catchError, finalize, of } from 'rxjs';

interface DashboardCard {
  title: string;
  value: number | string;
  icon: string;
  color: string;
  route?: string;
  trend?: {
    value: number;
    isPositive: boolean;
  };
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatListModule,
    MatChipsModule
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly dashboardService = inject(DashboardService);
  private readonly router = inject(Router);

  public readonly isLoading = signal(true);
  public readonly currentUser = signal<User | null>(null);
  public readonly dashboardStats = signal<DashboardStats | null>(null);
  public readonly recentActivity = signal<RecentActivity[]>([]);
  public readonly dashboardCards = signal<DashboardCard[]>([]);

  ngOnInit(): void {
    this.currentUser.set(this.authService.currentUser());
    this.loadDashboardData();
  }

  private loadDashboardData(): void {
    this.isLoading.set(true);
    
    // Cargar estadísticas del dashboard
    this.dashboardService.getDashboardStats()
      .pipe(
        catchError(error => {
          console.error('Error loading dashboard stats:', error);
          // Datos de fallback en caso de error
          return of({
            clients: { total: 0, active: 0 },
            projects: { total: 0, active: 0, completed: 0 },
            quotes: { total: 0, pending: 0, accepted: 0, accepted_value: 0, total_value: 0 },
            materials: { total: 0, low_stock: 0, inventory_value: 0 },
            monthly: { clients: 0, projects: 0, quotes: 0 }
          } as DashboardStats);
        }),
        finalize(() => this.isLoading.set(false))
      )
      .subscribe(stats => {
        this.dashboardStats.set(stats);
        this.updateDashboardCards(stats);
      });

    // Cargar actividad reciente
    this.dashboardService.getRecentActivity(5)
      .pipe(
        catchError(error => {
          console.error('Error loading recent activity:', error);
          return of([]);
        })
      )
      .subscribe(activity => {
        this.recentActivity.set(activity);
      });
  }

  private updateDashboardCards(stats: DashboardStats): void {
    this.dashboardCards.set([
      {
        title: 'Clientes Activos',
        value: stats.clients.active,
        icon: 'people',
        color: 'primary',
        route: '/clients'
      },
      {
        title: 'Proyectos en Curso',
        value: stats.projects.active,
        icon: 'construction',
        color: 'accent',
        route: '/projects'
      },
      {
        title: 'Cotizaciones Pendientes',
        value: stats.quotes.pending,
        icon: 'description',
        color: 'warn',
        route: '/quotes'
      },
      {
        title: 'Total de Clientes',
        value: stats.clients.total,
        icon: 'group',
        color: 'success',
        route: '/clients'
      }
    ]);
  }

  navigateToRoute(route?: string): void {
    if (route) {
      this.router.navigate([route]);
    }
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-BO', {
      style: 'currency',
      currency: 'BOB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 1) {
      return 'Ayer';
    } else if (diffDays < 7) {
      return `Hace ${diffDays} días`;
    } else {
      return date.toLocaleDateString('es-BO', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      });
    }
  }

  getActivityIcon(type: string): string {
    switch (type) {
      case 'client':
        return 'person';
      case 'project':
        return 'construction';
      case 'quote':
        return 'description';
      default:
        return 'info';
    }
  }

  getActivityIconClass(type: string): string {
    return `activity-icon-${type}`;
  }

  getActivityTypeLabel(type: string): string {
    switch (type) {
      case 'client':
        return 'Cliente';
      case 'project':
        return 'Proyecto';
      case 'quote':
        return 'Cotización';
      default:
        return 'Actividad';
    }
  }

  getActivityChipClass(type: string): string {
    return `activity-chip-${type}`;
  }
}