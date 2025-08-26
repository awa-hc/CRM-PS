import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, ActivatedRoute, Router } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup, FormArray, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDividerModule } from '@angular/material/divider';
import { MatTableModule } from '@angular/material/table';
import { Quote, CreateQuoteRequest, UpdateQuoteRequest } from '../../../core/models/quote.model';
import { QuoteService } from '../../../core/services/quote.service';
import { ClientService } from '../../../core/services/client.service';
import { ProjectService } from '../../../core/services/project.service';
import { Client } from '../../../core/models/client.model';
import { Project } from '../../../core/models/project.model';

@Component({
  selector: 'app-quote-form',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatFormFieldModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatProgressSpinnerModule,
    MatDividerModule,
    MatTableModule
  ],
  templateUrl: './quote-form.component.html',
  styleUrls: ['./quote-form.component.css']
})
export class QuoteFormComponent implements OnInit {
  quoteForm: FormGroup;
  clients: Client[] = [];
  projects: Project[] = [];
  filteredProjects: Project[] = [];
  loading = false;
  saving = false;
  isEditMode = false;
  quoteId?: number;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private quoteService: QuoteService,
    private clientService: ClientService,
    private projectService: ProjectService,
    private snackBar: MatSnackBar
  ) {
    this.quoteForm = this.createForm();
  }

  ngOnInit(): void {
    this.loadClients();
    this.loadProjects();
    this.checkEditMode();
  }

  createForm(): FormGroup {
    return this.fb.group({
      clientId: ['', Validators.required],
      projectId: [''],
      title: ['', Validators.required],
      description: [''],
      validUntil: [''],
      taxRate: [0, [Validators.min(0), Validators.max(100)]],
      discount: [0, Validators.min(0)],
      notes: [''],
      terms: [''],
      items: this.fb.array([this.createItemForm()])
    });
  }

  createItemForm(): FormGroup {
    return this.fb.group({
      description: ['', Validators.required],
      quantity: [1, [Validators.required, Validators.min(0.01)]],
      unit: ['pcs'],
      unitPrice: [0, [Validators.required, Validators.min(0)]],
      notes: ['']
    });
  }

  get items(): FormArray {
    return this.quoteForm.get('items') as FormArray;
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
        this.filteredProjects = this.projects;
      },
      error: (error) => {
        console.error('Error loading projects:', error);
      }
    });
  }

  checkEditMode(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id && id !== 'new') {
      this.isEditMode = true;
      this.quoteId = +id;
      this.loadQuote();
    }
  }

  loadQuote(): void {
    if (!this.quoteId) return;
    
    this.loading = true;
    this.quoteService.getQuote(this.quoteId).subscribe({
      next: (quote) => {
        this.populateForm(quote);
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

  populateForm(quote: Quote): void {
    // Clear existing items
    while (this.items.length > 0) {
      this.items.removeAt(0);
    }

    // Add items from quote
    if (quote.items && quote.items.length > 0) {
      quote.items.forEach(item => {
        this.items.push(this.fb.group({
          description: [item.description, Validators.required],
          quantity: [item.quantity, [Validators.required, Validators.min(0.01)]],
          unit: [item.unit],
          unitPrice: [item.unitPrice, [Validators.required, Validators.min(0)]],
          notes: [item.notes]
        }));
      });
    } else {
      this.items.push(this.createItemForm());
    }

    // Populate form
    this.quoteForm.patchValue({
      clientId: quote.clientId,
      projectId: quote.projectId || '',
      title: quote.title,
      description: quote.description,
      validUntil: quote.validUntil ? new Date(quote.validUntil) : '',
      taxRate: quote.taxRate,
      discount: quote.discount,
      notes: quote.notes,
      terms: quote.terms
    });

    // Filter projects by client
    this.onClientChange(quote.clientId);
  }

  onClientChange(clientId: number): void {
    if (clientId) {
      this.filteredProjects = this.projects.filter(p => p.client_id === clientId);
    } else {
      this.filteredProjects = this.projects;
    }
    
    // Clear project selection if it doesn't belong to selected client
    const currentProjectId = this.quoteForm.get('projectId')?.value;
    if (currentProjectId && !this.filteredProjects.find(p => p.id === currentProjectId)) {
      this.quoteForm.patchValue({ projectId: '' });
    }
  }

  addItem(): void {
    this.items.push(this.createItemForm());
  }

  removeItem(index: number): void {
    if (this.items.length > 1) {
      this.items.removeAt(index);
    }
  }

  calculateItemTotal(index: number): void {
    // This will trigger the getItemTotal calculation
  }

  getItemTotal(index: number): number {
    const item = this.items.at(index);
    const quantity = item.get('quantity')?.value || 0;
    const unitPrice = item.get('unitPrice')?.value || 0;
    return quantity * unitPrice;
  }

  getSubtotal(): number {
    let subtotal = 0;
    for (let i = 0; i < this.items.length; i++) {
      subtotal += this.getItemTotal(i);
    }
    return subtotal;
  }

  getTaxAmount(): number {
    const subtotal = this.getSubtotal();
    const discount = this.quoteForm.get('discount')?.value || 0;
    const taxRate = this.quoteForm.get('taxRate')?.value || 0;
    return (subtotal - discount) * (taxRate / 100);
  }

  getTotal(): number {
    const subtotal = this.getSubtotal();
    const discount = this.quoteForm.get('discount')?.value || 0;
    const taxAmount = this.getTaxAmount();
    return subtotal - discount + taxAmount;
  }

  onSubmit(): void {
    if (this.quoteForm.valid) {
      this.saving = true;
      
      const formValue = this.quoteForm.value;
      const quoteData = {
        ...formValue,
        validUntil: formValue.validUntil ? formValue.validUntil.toISOString() : null,
        projectId: formValue.projectId || null
      };

      if (this.isEditMode && this.quoteId) {
        this.updateQuote(quoteData);
      } else {
        this.createQuote(quoteData);
      }
    } else {
      this.markFormGroupTouched();
    }
  }

  createQuote(quoteData: CreateQuoteRequest): void {
    this.quoteService.createQuote(quoteData).subscribe({
      next: (quote) => {
        this.saving = false;
        this.snackBar.open('Cotización creada exitosamente', 'Cerrar', {
          duration: 3000
        });
        this.router.navigate(['/quotes', quote.id]);
      },
      error: (error) => {
        console.error('Error creating quote:', error);
        this.saving = false;
        this.snackBar.open('Error al crear la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  updateQuote(quoteData: UpdateQuoteRequest): void {
    if (!this.quoteId) return;
    
    this.quoteService.updateQuote(this.quoteId, quoteData).subscribe({
      next: (quote) => {
        this.saving = false;
        this.snackBar.open('Cotización actualizada exitosamente', 'Cerrar', {
          duration: 3000
        });
        this.router.navigate(['/quotes', quote.id]);
      },
      error: (error) => {
        console.error('Error updating quote:', error);
        this.saving = false;
        this.snackBar.open('Error al actualizar la cotización', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  markFormGroupTouched(): void {
    Object.keys(this.quoteForm.controls).forEach(key => {
      const control = this.quoteForm.get(key);
      control?.markAsTouched();
      
      if (control instanceof FormArray) {
        control.controls.forEach(item => {
          Object.keys((item as FormGroup).controls).forEach(itemKey => {
            item.get(itemKey)?.markAsTouched();
          });
        });
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/quotes']);
  }
}