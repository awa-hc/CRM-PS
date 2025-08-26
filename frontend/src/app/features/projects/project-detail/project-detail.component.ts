import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDividerModule } from '@angular/material/divider';
import { MatTabsModule } from '@angular/material/tabs';
import { ProjectService } from '../../../core/services/project.service';
import { ClientService } from '../../../core/services/client.service';
import { Project } from '../../../core/models/project.model';
import { Client } from '../../../core/models/client.model';

@Component({
  selector: 'app-project-detail',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatDividerModule,
    MatTabsModule
  ],
  templateUrl: './project-detail.component.html',
  styleUrls: ['./project-detail.component.css']
})
export class ProjectDetailComponent implements OnInit {
  project?: Project;
  client?: Client;
  loading = true;

  constructor(
    private projectService: ProjectService,
    private clientService: ClientService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.loadProject();
  }

  loadProject(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      this.router.navigate(['/app/projects']);
      return;
    }

    const projectId = parseInt(id, 10);
    this.projectService.getProject(projectId).subscribe({
      next: (project) => {
        this.project = project;
        this.loadClient(project.client_id);
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading project:', error);
        this.snackBar.open('Error al cargar proyecto', 'Cerrar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  loadClient(clientId: number): void {
    this.clientService.getClient(clientId).subscribe({
      next: (client) => {
        this.client = client;
      },
      error: (error) => {
        console.error('Error loading client:', error);
      }
    });
  }

  editProject(): void {
    if (this.project) {
      this.router.navigate(['/app/projects', this.project.id, 'edit']);
    }
  }

  goBack(): void {
    this.router.navigate(['/app/projects']);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-ES', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount);
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

  getTypeLabel(type: string): string {
    const labels: { [key: string]: string } = {
      'residential': 'Residencial',
      'commercial': 'Comercial',
      'industrial': 'Industrial',
      'renovation': 'Renovación',
      'maintenance': 'Mantenimiento'
    };
    return labels[type] || type;
  }

  getVarianceIcon(): string {
    if (!this.project?.actualCost) return 'up';
    return this.project.actualCost > this.project.budget ? 'up' : 'down';
  }

  getVarianceClass(): string {
    if (!this.project?.actualCost) return '';
    return this.project.actualCost > this.project.budget ? 'negative' : 'positive';
  }

  getVarianceText(): string {
    if (!this.project?.actualCost) return 'N/A';
    const variance = this.project.actualCost - this.project.budget;
    const percentage = (variance / this.project.budget) * 100;
    const sign = variance > 0 ? '+' : '';
    const formattedVariance = this.formatCurrency(Math.abs(variance));
    const formattedPercentage = percentage.toFixed(1);
    return sign + formattedVariance + ' (' + sign + formattedPercentage + '%)';
  }
}