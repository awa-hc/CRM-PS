import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDividerModule } from '@angular/material/divider';
import { Quote, QuoteFilters } from '../../../core/models/quote.model';
import { QuoteService } from '../../../core/services/quote.service';
import { ClientService } from '../../../core/services/client.service';
import { ProjectService } from '../../../core/services/project.service';
import { Client } from '../../../core/models/client.model';
import { Project } from '../../../core/models/project.model';

@Component({
  selector: 'app-quotes-list',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    MatTableModule,
    MatPaginatorModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatFormFieldModule,
    MatSelectModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatCardModule,
    MatMenuModule,
    MatTooltipModule,
    MatDividerModule
  ],
  templateUrl: './quotes-list.component.html',
  styleUrls: ['./quotes-list.component.css']
})
export class QuotesListComponent implements OnInit {
  quotes: Quote[] = [];
  clients: Client[] = [];
  projects: Project[] = [];
  loading = false;
  totalQuotes = 0;
  currentPage = 1;
  pageSize = 10;
  
  filters: QuoteFilters = {
    search: '',
    clientId: undefined,
    projectId: undefined,
    status: ''
  };

  displayedColumns: string[] = [
    'quoteNumber',
    'title', 
    'client',
    'project',
    'status',
    'total',
    'validUntil',
    'actions'
  ];

  constructor(
    private quoteService: QuoteService,
    private clientService: ClientService,
    private projectService: ProjectService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadClients();
    this.loadProjects();
    this.loadQuotes();
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

  loadProjects(): void {
    this.projectService.getProjects().subscribe({
      next: (response) => {
        this.projects = response.projects;
      },
      error: (error) => {
        console.error('Error loading projects:', error);
      }
    });
  }

  loadQuotes(): void {
    this.loading = true;
    
    const cleanFilters = { ...this.filters };
    if (!cleanFilters.search) delete cleanFilters.search;
    if (!cleanFilters.clientId) delete cleanFilters.clientId;
    if (!cleanFilters.projectId) delete cleanFilters.projectId;
    if (!cleanFilters.status) delete cleanFilters.status;

    this.quoteService.getQuotes(cleanFilters, this.currentPage, this.pageSize).subscribe({
      next: (response) => {
        this.quotes = response.quotes;
        this.totalQuotes = response.total;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading quotes:', error);
        this.loading = false;
        this.snackBar.open('Error al cargar las cotizaciones', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  onFilterChange(): void {
    this.currentPage = 1;
    this.loadQuotes();
  }

  onPageChange(event: PageEvent): void {
    this.currentPage = event.pageIndex + 1;
    this.pageSize = event.pageSize;
    this.loadQuotes();
  }

  clearFilters(): void {
    this.filters = {
      search: '',
      clientId: undefined,
      projectId: undefined,
      status: ''
    };
    this.currentPage = 1;
    this.loadQuotes();
  }

  sendQuote(quote: Quote): void {
    this.quoteService.sendQuote(quote.id).subscribe({
      next: () => {
        this.snackBar.open('Cotización enviada exitosamente', 'Cerrar', {
          duration: 3000
        });
        this.loadQuotes();
      },
      error: (error) => {
        console.error('Error sending quote:', error);
        this.snackBar.open('Error al enviar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  duplicateQuote(quote: Quote): void {
    this.quoteService.duplicateQuote(quote.id).subscribe({
      next: () => {
        this.snackBar.open('Cotización duplicada exitosamente', 'Cerrar', {
          duration: 3000
        });
        this.loadQuotes();
      },
      error: (error) => {
        console.error('Error duplicating quote:', error);
        this.snackBar.open('Error al duplicar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  downloadPDF(quote: Quote): void {
    this.quoteService.generatePDF(quote.id).subscribe({
      next: (blob) => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `cotizacion-${quote.quoteNumber}.pdf`;
        link.click();
        window.URL.revokeObjectURL(url);
      },
      error: (error) => {
        console.error('Error downloading PDF:', error);
        this.snackBar.open('Error al descargar el PDF', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  deleteQuote(quote: Quote): void {
    if (confirm(`¿Está seguro de que desea eliminar la cotización ${quote.quoteNumber}?`)) {
      this.quoteService.deleteQuote(quote.id).subscribe({
        next: () => {
          this.snackBar.open('Cotización eliminada exitosamente', 'Cerrar', {
            duration: 3000
          });
          this.loadQuotes();
        },
        error: (error) => {
          console.error('Error deleting quote:', error);
          this.snackBar.open('Error al eliminar la cotización', 'Cerrar', {
            duration: 3000
          });
        }
      });
    }
  }

  getStatusLabel(status: string): string {
    const statusLabels: { [key: string]: string } = {
      'draft': 'Borrador',
      'sent': 'Enviada',
      'accepted': 'Aceptada',
      'rejected': 'Rechazada',
      'expired': 'Expirada'
    };
    return statusLabels[status] || status;
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('es-ES');
  }

  isExpired(validUntil: string): boolean {
    return new Date(validUntil) < new Date();
  }
}