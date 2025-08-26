import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, ActivatedRoute, Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatTableModule } from '@angular/material/table';
import { MatMenuModule } from '@angular/material/menu';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { Quote } from '../../../core/models/quote.model';
import { QuoteService } from '../../../core/services/quote.service';
import { ClientService } from '../../../core/services/client.service';
import { ProjectService } from '../../../core/services/project.service';
import { Client } from '../../../core/models/client.model';
import { Project } from '../../../core/models/project.model';

@Component({
  selector: 'app-quote-detail',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatDividerModule,
    MatTableModule,
    MatMenuModule
  ],
  templateUrl: './quote-detail.component.html',
  styleUrls: ['./quote-detail.component.css']
})
export class QuoteDetailComponent implements OnInit {
  quote?: Quote;
  client?: Client;
  project?: Project;
  loading = false;
  displayedColumns: string[] = ['description', 'quantity', 'unitPrice', 'total'];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private quoteService: QuoteService,
    private clientService: ClientService,
    private projectService: ProjectService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.loadQuote();
  }

  loadQuote(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      this.router.navigate(['/app/quotes']);
      return;
    }

    this.loading = true;
    this.quoteService.getQuote(+id).subscribe({
      next: (quote) => {
        this.quote = quote;
        this.loadClient();
        if (quote.projectId) {
          this.loadProject();
        }
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading quote:', error);
        this.loading = false;
        this.snackBar.open('Error al cargar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  loadClient(): void {
    if (!this.quote?.clientId) return;
    
    this.clientService.getClient(this.quote.clientId).subscribe({
      next: (client) => {
        this.client = client;
      },
      error: (error) => {
        console.error('Error loading client:', error);
      }
    });
  }

  loadProject(): void {
    if (!this.quote?.projectId) return;
    
    this.projectService.getProject(this.quote.projectId).subscribe({
      next: (project) => {
        this.project = project;
      },
      error: (error) => {
        console.error('Error loading project:', error);
      }
    });
  }

  editQuote(): void {
    if (this.quote) {
      this.router.navigate(['/app/quotes', this.quote.id, 'edit']);
    }
  }

  sendQuote(): void {
    if (!this.quote) return;
    
    this.quoteService.sendQuote(this.quote.id).subscribe({
      next: () => {
        this.snackBar.open('Cotización enviada exitosamente', 'Cerrar', {
          duration: 3000
        });
        this.loadQuote(); // Reload to update status
      },
      error: (error) => {
        console.error('Error sending quote:', error);
        this.snackBar.open('Error al enviar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  duplicateQuote(): void {
    if (!this.quote) return;
    
    this.quoteService.duplicateQuote(this.quote.id).subscribe({
      next: (newQuote) => {
        this.snackBar.open('Cotización duplicada exitosamente', 'Cerrar', {
          duration: 3000
        });
        this.router.navigate(['/app/quotes', newQuote.id]);
      },
      error: (error) => {
        console.error('Error duplicating quote:', error);
        this.snackBar.open('Error al duplicar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  downloadPDF(): void {
    if (!this.quote) return;
    
    this.quoteService.generatePDF(this.quote.id).subscribe({
      next: (blob) => {
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `cotizacion-${this.quote!.quoteNumber}.pdf`;
        link.click();
        window.URL.revokeObjectURL(url);
      },
      error: (error) => {
        console.error('Error generating PDF:', error);
        this.snackBar.open('Error al generar el PDF', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  acceptQuote(): void {
    if (!this.quote) return;
    
    this.quoteService.acceptQuote(this.quote.id).subscribe({
      next: () => {
        this.snackBar.open('Cotización aceptada', 'Cerrar', {
          duration: 3000
        });
        this.loadQuote(); // Reload to update status
      },
      error: (error) => {
        console.error('Error accepting quote:', error);
        this.snackBar.open('Error al aceptar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  rejectQuote(): void {
    if (!this.quote) return;
    
    this.quoteService.rejectQuote(this.quote.id).subscribe({
      next: () => {
        this.snackBar.open('Cotización rechazada', 'Cerrar', {
          duration: 3000
        });
        this.loadQuote(); // Reload to update status
      },
      error: (error) => {
        console.error('Error rejecting quote:', error);
        this.snackBar.open('Error al rechazar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/app/quotes']);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
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

  isExpired(validUntil: string): boolean {
    return new Date(validUntil) < new Date();
  }
}