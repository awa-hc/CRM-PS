import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { ClientService } from '../../../core/services/client.service';
import { Client } from '../../../core/models/client.model';

@Component({
  selector: 'app-client-detail',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatDividerModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './client-detail.component.html',
  styleUrls: ['./client-detail.component.css']
})
export class ClientDetailComponent implements OnInit {
  client: Client | null = null;
  loading = false;
  clientId: number | null = null;

  constructor(
    private clientService: ClientService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      if (params['id']) {
        this.clientId = +params['id'];
        this.loadClient();
      }
    });
  }

  loadClient(): void {
    if (!this.clientId) return;
    
    this.loading = true;
    this.clientService.getClient(this.clientId).subscribe({
      next: (client) => {
        this.client = client;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading client:', error);
        this.snackBar.open('Error al cargar el cliente', 'Cerrar', {
          duration: 3000
        });
        this.loading = false;
      }
    });
  }

  editClient(): void {
    if (this.clientId) {
      this.router.navigate(['/app/clients', this.clientId, 'edit']);
    }
  }

  goBack(): void {
    this.router.navigate(['/app/clients']);
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
}