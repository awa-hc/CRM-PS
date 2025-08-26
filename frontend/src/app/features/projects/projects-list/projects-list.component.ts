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
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { ProjectService, ProjectFilters } from '../../../core/services/project.service';
import { Project } from '../../../core/models/project.model';
import { ClientService } from '../../../core/services/client.service';
import { Client } from '../../../core/models/client.model';

@Component({
  selector: 'app-projects-list',
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
    MatCardModule,
    MatProgressBarModule
  ],
  templateUrl: './projects-list.component.html',
  styleUrls: ['./projects-list.component.css']
})
export class ProjectsListComponent implements OnInit {
  projects: Project[] = [];
  clients: Client[] = [];
  loading = false;
  totalProjects = 0;
  pageSize = 10;
  currentPage = 0;

  filters: ProjectFilters = {
    search: '',
    status: '',
    priority: '',
    type: '',
    client_id: undefined,
    page: 1,
    limit: 10
  };

  displayedColumns: string[] = ['code', 'name', 'status', 'priority', 'progress', 'budget', 'actions'];

  constructor(
    private projectService: ProjectService,
    private clientService: ClientService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadClients();
    this.loadProjects();
  }

  loadClients(): void {
    this.clientService.getClients({ limit: 1000 }).subscribe({
      next: (response) => {
        this.clients = response.clients;
      },
      error: (error) => {
        console.error('Error loading clients:', error);
      }
    });
  }

  loadProjects(): void {
    this.loading = true;
    this.projectService.getProjects(this.filters).subscribe({
      next: (response) => {
        this.projects = response.projects;
        this.totalProjects = response.total;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading projects:', error);
        this.snackBar.open('Error al cargar proyectos', 'Cerrar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  onFilterChange(): void {
    this.filters.page = 1;
    this.currentPage = 0;
    this.loadProjects();
  }

  onPageChange(event: PageEvent): void {
    this.filters.page = event.pageIndex + 1;
    this.filters.limit = event.pageSize;
    this.pageSize = event.pageSize;
    this.currentPage = event.pageIndex;
    this.loadProjects();
  }

  clearFilters(): void {
    this.filters = {
      search: '',
      status: '',
      priority: '',
      type: '',
      client_id: undefined,
      page: 1,
      limit: 10
    };
    this.currentPage = 0;
    this.loadProjects();
  }

  deleteProject(project: Project): void {
    if (confirm(`¿Está seguro de que desea eliminar el proyecto "${project.name}"?`)) {
      this.projectService.deleteProject(project.id).subscribe({
        next: () => {
          this.snackBar.open('Proyecto eliminado exitosamente', 'Cerrar', { duration: 3000 });
          this.loadProjects();
        },
        error: (error) => {
          console.error('Error deleting project:', error);
          this.snackBar.open('Error al eliminar proyecto', 'Cerrar', { duration: 3000 });
        }
      });
    }
  }

  getStatusLabel(status: string): string {
    const labels: { [key: string]: string } = {
      'planning': 'Planificación',
      'active': 'Activo',
      'completed': 'Completado',
      'cancelled': 'Cancelado',
      'on_hold': 'En Espera'
    };
    return labels[status] || status;
  }

  getPriorityLabel(priority: string): string {
    const labels: { [key: string]: string } = {
      'low': 'Baja',
      'medium': 'Media',
      'high': 'Alta',
      'urgent': 'Urgente'
    };
    return labels[priority] || priority;
  }
}