import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { ClientService, ClientFilters } from '../../../core/services/client.service';
import { Client } from '../../../core/models/client.model';

@Component({
  selector: 'app-clients-list',
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
    MatSnackBarModule,
    MatDialogModule,
    MatCardModule
  ],
  templateUrl: './clients-list.component.html',
  styleUrls: ['./clients-list.component.css']
})
export class ClientsListComponent implements OnInit {
  clients: Client[] = [];
  loading = false;
  totalClients = 0;
  currentPage = 1;
  pageSize = 10;
  
  filters: ClientFilters = {
    search: '',
    contactType: undefined,
    isActive: undefined,
    page: 1,
    limit: 10
  };

  displayedColumns: string[] = ['name', 'contact', 'type', 'status', 'location', 'actions'];

  constructor(
    private clientService: ClientService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadClients();
  }

  loadClients(): void {
    this.loading = true;
    this.filters.page = this.currentPage;
    this.filters.limit = this.pageSize;

    this.clientService.getClients(this.filters).subscribe({
      next: (response) => {
        this.clients = response.clients;
        this.totalClients = response.total;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading clients:', error);
        this.snackBar.open('Error al cargar los clientes', 'Cerrar', {
          duration: 3000
        });
        this.loading = false;
      }
    });
  }

  onPageChange(event: PageEvent): void {
    this.currentPage = event.pageIndex + 1;
    this.pageSize = event.pageSize;
    this.loadClients();
  }

  onFilterChange(): void {
    this.currentPage = 1;
    this.loadClients();
  }

  clearFilters(): void {
    this.filters = {
      search: '',
      contactType: undefined,
      isActive: undefined,
      page: 1,
      limit: this.pageSize
    };
    this.currentPage = 1;
    this.loadClients();
  }

  deleteClient(client: Client): void {
    if (confirm(`¿Está seguro de que desea eliminar el cliente "${client.name}"?`)) {
      this.clientService.deleteClient(client.id).subscribe({
        next: () => {
          this.snackBar.open('Cliente eliminado correctamente', 'Cerrar', {
            duration: 3000
          });
          this.loadClients();
        },
        error: (error) => {
          console.error('Error deleting client:', error);
          this.snackBar.open('Error al eliminar el cliente', 'Cerrar', {
            duration: 3000
          });
        }
      });
    }
  }

  getCompaniesCount(): number {
    return this.clients.filter(client => client.contactType === 'company').length;
  }

  getIndividualsCount(): number {
    return this.clients.filter(client => client.contactType === 'individual').length;
  }

  getActiveClientsCount(): number {
    return this.clients.filter(client => client.isActive).length;
  }
}