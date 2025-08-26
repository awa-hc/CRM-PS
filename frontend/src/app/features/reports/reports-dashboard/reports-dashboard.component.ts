import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatListModule } from '@angular/material/list';
import { MatChipsModule } from '@angular/material/chips';
import { ReportsService, DashboardStats } from '../../../core/services/reports.service';

@Component({
  selector: 'app-reports-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatGridListModule,
    MatListModule,
    MatChipsModule
  ],
  templateUrl: './reports-dashboard.component.html',
  styleUrls: ['./reports-dashboard.component.css']
})
export class ReportsDashboardComponent implements OnInit {
  loading = true;
  stats: DashboardStats | null = null;

  constructor(private reportsService: ReportsService) {}

  ngOnInit(): void {
    this.loadDashboardStats();
  }

  loadDashboardStats(): void {
    this.loading = true;
    this.reportsService.getDashboardStats().subscribe({
      next: (stats) => {
        this.stats = stats;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading dashboard stats:', error);
        this.loading = false;
      }
    });
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-ES', {
      style: 'currency',
      currency: 'EUR'
    }).format(amount);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }

  getActivityIcon(type: string): string {
    switch (type) {
      case 'client': return 'person';
      case 'project': return 'construction';
      case 'quote': return 'description';
      default: return 'info';
    }
  }

  getActivityTypeLabel(type: string): string {
    switch (type) {
      case 'client': return 'Cliente';
      case 'project': return 'Proyecto';
      case 'quote': return 'Cotización';
      default: return 'Actividad';
    }
  }

  getStatusColor(status: string): string {
    switch (status.toLowerCase()) {
      case 'active':
      case 'activo':
      case 'completed':
      case 'completado':
        return 'primary';
      case 'pending':
      case 'pendiente':
        return 'accent';
      case 'cancelled':
      case 'cancelado':
        return 'warn';
      default:
        return '';
    }
  }

  getDeadlineClass(daysLeft: number): string {
    if (daysLeft <= 3) return 'deadline-urgent';
    if (daysLeft <= 7) return 'deadline-warning';
    return 'deadline-normal';
  }

  getDeadlineChipClass(daysLeft: number): string {
    if (daysLeft <= 3) return 'chip-urgent';
    if (daysLeft <= 7) return 'chip-warning';
    return 'chip-normal';
  }
}