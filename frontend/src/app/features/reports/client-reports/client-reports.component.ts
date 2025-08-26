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
import { MatTableModule } from '@angular/material/table';
import { ReportsService, ClientReport, ReportFilters } from '../../../core/services/reports.service';
import { ClientService, Client } from '../../../core/services/client.service';

@Component({
  selector: 'app-client-reports',
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
    MatProgressBarModule,
    MatTableModule
  ],
  templateUrl: './client-reports.component.html',
  styleUrls: ['./client-reports.component.css']
})
export class ClientReportsComponent implements OnInit {
  filtersForm: FormGroup;
  loading = false;
  report: ClientReport | null = null;
  clients: Client[] = [];
  performanceColumns: string[] = ['name', 'projectCount', 'totalValue', 'averageProjectValue', 'completedProjects', 'completionRate', 'lastProject'];

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
      clientType: ['']
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
    if (formValue.clientType) {
      filters.clientType = formValue.clientType;
    }

    this.loadClientReport(filters);
  }

  clearFilters(): void {
    this.filtersForm.reset();
    this.setDefaultFilters();
    this.applyFilters();
  }

  loadClientReport(filters: ReportFilters): void {
    this.loading = true;
    this.reportsService.getClientReport(filters).subscribe({
      next: (report) => {
        this.report = report;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading client report:', error);
        this.loading = false;
        this.snackBar.open('Error al cargar el reporte de clientes', 'Cerrar', {
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
    if (formValue.clientType) {
      filters.clientType = formValue.clientType;
    }

    this.reportsService.exportClientReport(filters, format).subscribe({
      next: (blob) => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `reporte-clientes.${format === 'pdf' ? 'pdf' : 'xlsx'}`;
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

  getTypeLabel(type: string): string {
    const labels: { [key: string]: string } = {
      'individual': 'Individual',
      'company': 'Empresa'
    };
    return labels[type] || type;
  }

  getTypeClass(type: string): string {
    return type;
  }

  getTypePercentage(count: number): number {
    if (!this.report || this.report.totalClients === 0) return 0;
    return (count / this.report.totalClients) * 100;
  }

  getLocationPercentage(count: number): number {
    if (!this.report || this.report.totalClients === 0) return 0;
    return (count / this.report.totalClients) * 100;
  }

  getCompletionColor(rate: number): 'primary' | 'accent' | 'warn' {
    if (rate >= 80) return 'primary';
    if (rate >= 60) return 'accent';
    return 'warn';
  }

  getActivityIcon(type: string): string {
    const icons: { [key: string]: string } = {
      'new_client': 'person_add',
      'new_project': 'add_circle',
      'quote_sent': 'send',
      'project_completed': 'check_circle'
    };
    return icons[type] || 'info';
  }

  getActivityClass(type: string): string {
    return type.replace('_', '-');
  }
}