import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { ReportsService, ProjectReport, ReportFilters } from '../../../core/services/reports.service';
import { ClientService, Client } from '../../../core/services/client.service';

@Component({
  selector: 'app-project-reports',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatMenuModule,
    MatDividerModule,
    MatChipsModule,
    MatProgressBarModule
  ],
  templateUrl: './project-reports.component.html',
  styleUrls: ['./project-reports.component.css']
})
export class ProjectReportsComponent implements OnInit {
  filtersForm: FormGroup;
  loading = false;
  report: ProjectReport | null = null;
  clients: Client[] = [];

  constructor(
    private fb: FormBuilder,
    private reportsService: ReportsService,
    private clientService: ClientService,
    private snackBar: MatSnackBar
  ) {
    this.filtersForm = this.fb.group({
      startDate: [''],
      endDate: [''],
      clientId: [''],
      status: ['']
    });
  }

  ngOnInit(): void {
    this.loadClients();
    this.setDefaultFilters();
    this.applyFilters();
  }

  loadClients(): void {
    this.clientService.getClients().subscribe({
      next: (response) => {
        this.clients = response.clients;
      },
      error: (error) => {
        console.error('Error loading clients:', error);
      }
    });
  }

  setDefaultFilters(): void {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setMonth(startDate.getMonth() - 12); // Last 12 months

    this.filtersForm.patchValue({
      startDate: startDate,
      endDate: endDate
    });
  }

  applyFilters(): void {
    const formValue = this.filtersForm.value;
    const filters: ReportFilters = {};

    if (formValue.startDate) {
      filters.startDate = formValue.startDate.toISOString().split('T')[0];
    }
    if (formValue.endDate) {
      filters.endDate = formValue.endDate.toISOString().split('T')[0];
    }
    if (formValue.clientId) {
      filters.clientId = formValue.clientId;
    }
    if (formValue.status) {
      filters.status = formValue.status;
    }

    this.loadProjectReport(filters);
  }

  clearFilters(): void {
    this.filtersForm.reset();
    this.setDefaultFilters();
    this.applyFilters();
  }

  loadProjectReport(filters: ReportFilters): void {
    this.loading = true;
    this.reportsService.getProjectReport(filters).subscribe({
      next: (report) => {
        this.report = report;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading project report:', error);
        this.loading = false;
        this.snackBar.open('Error al cargar el reporte de proyectos', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  exportReport(format: 'pdf' | 'excel'): void {
    if (!this.report) return;

    const formValue = this.filtersForm.value;
    const filters: ReportFilters = {};

    if (formValue.startDate) {
      filters.startDate = formValue.startDate.toISOString().split('T')[0];
    }
    if (formValue.endDate) {
      filters.endDate = formValue.endDate.toISOString().split('T')[0];
    }
    if (formValue.clientId) {
      filters.clientId = formValue.clientId;
    }
    if (formValue.status) {
      filters.status = formValue.status;
    }

    this.reportsService.exportProjectReport(filters, format).subscribe({
      next: (blob) => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `reporte-proyectos.${format === 'pdf' ? 'pdf' : 'xlsx'}`;
        link.click();
        window.URL.revokeObjectURL(url);
        
        this.snackBar.open(`Reporte exportado como ${format.toUpperCase()}`, 'Cerrar', {
          duration: 3000
        });
      },
      error: (error) => {
        console.error('Error exporting report:', error);
        this.snackBar.open('Error al exportar el reporte', 'Cerrar', {
          duration: 3000
        });
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
    return new Date(date).toLocaleDateString('es-ES');
  }

  getStatusLabel(status: string): string {
    const labels: { [key: string]: string } = {
      'planning': 'Planificación',
      'in_progress': 'En Progreso',
      'on_hold': 'En Pausa',
      'completed': 'Completado',
      'cancelled': 'Cancelado'
    };
    return labels[status] || status;
  }

  getStatusClass(status: string): string {
    return status;
  }

  getStatusColor(status: string): 'primary' | 'accent' | 'warn' {
    switch (status) {
      case 'completed':
        return 'primary';
      case 'in_progress':
        return 'accent';
      case 'cancelled':
        return 'warn';
      default:
        return 'primary';
    }
  }

  getStatusPercentage(count: number): number {
    if (!this.report || this.report.totalProjects === 0) return 0;
    return (count / this.report.totalProjects) * 100;
  }

  getBudgetUtilization(): number {
    if (!this.report || this.report.totalBudget === 0) return 0;
    return (this.report.totalSpent / this.report.totalBudget) * 100;
  }

  getBudgetColor(): 'primary' | 'accent' | 'warn' {
    const utilization = this.getBudgetUtilization();
    if (utilization > 90) return 'warn';
    if (utilization > 75) return 'accent';
    return 'primary';
  }
}