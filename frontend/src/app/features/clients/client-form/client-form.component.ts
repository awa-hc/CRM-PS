import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ClientService } from '../../../core/services/client.service';
import { Client, CreateClientRequest, UpdateClientRequest } from '../../../core/models/client.model';

@Component({
  selector: 'app-client-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './client-form.component.html',
  styleUrls: ['./client-form.component.css']
})
export class ClientFormComponent implements OnInit {
  clientForm: FormGroup;
  isEditMode = false;
  clientId: number | null = null;
  loading = false;
  saving = false;

  constructor(
    private fb: FormBuilder,
    private clientService: ClientService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    this.clientForm = this.createForm();
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      if (params['id']) {
        this.isEditMode = true;
        this.clientId = +params['id'];
        this.loadClient();
      }
    });

    // Watch contactType changes to show/hide company fields
    this.clientForm.get('contactType')?.valueChanges.subscribe(value => {
      if (value === 'company') {
        this.clientForm.get('company')?.setValidators([Validators.required]);
      } else {
        this.clientForm.get('company')?.clearValidators();
        this.clientForm.get('company')?.setValue('');
        this.clientForm.get('taxId')?.setValue('');
      }
      this.clientForm.get('company')?.updateValueAndValidity();
    });
  }

  createForm(): FormGroup {
    return this.fb.group({
      contactType: ['individual', [Validators.required]],
      firstName: [''],
      lastName: [''],
      companyName: [''],
      email: ['', [Validators.required, Validators.email]],
      phone: ['', [Validators.required]],
      address: [''],
      city: [''],
      state: [''],
      zipCode: [''],
      company: [''],
      taxId: [''],
      notes: ['']
    });
  }

  loadClient(): void {
    if (!this.clientId) return;
    
    this.loading = true;
    this.clientService.getClient(this.clientId).subscribe({
      next: (client) => {
        // Transform client data to match form structure
        const formData = {
          contactType: client.contactType || 'individual',
          firstName: this.extractFirstName(client.name, client.contactType),
          lastName: this.extractLastName(client.name, client.contactType),
          companyName: client.contactType === 'company' ? client.name : '',
          email: client.email || '',
          phone: client.phone || '',
          address: client.address || '',
          city: client.city || '',
          state: client.state || '',
          zipCode: client.zipCode || '',
          company: client.company || '',
          taxId: client.taxId || '',
          notes: client.notes || ''
        };
        
        this.clientForm.patchValue(formData);
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading client:', error);
        this.snackBar.open('Error al cargar el cliente', 'Cerrar', {
          duration: 3000
        });
        this.loading = false;
        this.goBack();
      }
    });
  }

  private extractFirstName(name: string, contactType: string): string {
    if (contactType === 'company' || !name) return '';
    const parts = name.split(' ');
    return parts[0] || '';
  }

  private extractLastName(name: string, contactType: string): string {
    if (contactType === 'company' || !name) return '';
    const parts = name.split(' ');
    return parts.slice(1).join(' ') || '';
  }

  onSubmit(): void {
    if (this.clientForm.invalid) {
      this.markFormGroupTouched();
      return;
    }

    this.saving = true;
    const formValue = this.clientForm.value;

    // Transform form data to match backend expectations
    const transformedData = {
      name: this.buildClientName(formValue),
      email: formValue.email,
      phone: formValue.phone,
      address: formValue.address || '',
      city: formValue.city || '',
      state: formValue.state || '',
      zipCode: formValue.zipCode || '',
      contactType: formValue.contactType,
      notes: formValue.notes || '',
      company: formValue.contactType === 'company' ? (formValue.companyName || formValue.company) : undefined,
      taxId: formValue.contactType === 'company' ? formValue.taxId : undefined
    };

    if (this.isEditMode && this.clientId) {
      const updateRequest: UpdateClientRequest = transformedData;

      this.clientService.updateClient(this.clientId, updateRequest).subscribe({
        next: () => {
          this.snackBar.open('Cliente actualizado correctamente', 'Cerrar', {
            duration: 3000
          });
          this.goBack();
        },
        error: (error) => {
          console.error('Error updating client:', error);
          this.snackBar.open('Error al actualizar el cliente', 'Cerrar', {
            duration: 3000
          });
          this.saving = false;
        }
      });
    } else {
      const createRequest: CreateClientRequest = transformedData;

      this.clientService.createClient(createRequest).subscribe({
        next: () => {
          this.snackBar.open('Cliente creado correctamente', 'Cerrar', {
            duration: 3000
          });
          this.goBack();
        },
        error: (error) => {
          console.error('Error creating client:', error);
          this.snackBar.open('Error al crear el cliente', 'Cerrar', {
            duration: 3000
          });
          this.saving = false;
        }
      });
    }
  }

  private buildClientName(formValue: any): string {
    if (formValue.contactType === 'company') {
      return formValue.companyName || '';
    } else {
      const firstName = formValue.firstName || '';
      const lastName = formValue.lastName || '';
      return `${firstName} ${lastName}`.trim();
    }
  }

  markFormGroupTouched(): void {
    Object.keys(this.clientForm.controls).forEach(key => {
      const control = this.clientForm.get(key);
      control?.markAsTouched();
    });
  }

  onTypeChange(): void {
    const contactType = this.clientForm.get('contactType')?.value;
    if (contactType === 'company') {
      this.clientForm.get('companyName')?.setValidators([Validators.required]);
      this.clientForm.get('firstName')?.clearValidators();
      this.clientForm.get('lastName')?.clearValidators();
    } else {
      this.clientForm.get('firstName')?.setValidators([Validators.required]);
      this.clientForm.get('lastName')?.setValidators([Validators.required]);
      this.clientForm.get('companyName')?.clearValidators();
    }
    this.clientForm.get('companyName')?.updateValueAndValidity();
    this.clientForm.get('firstName')?.updateValueAndValidity();
    this.clientForm.get('lastName')?.updateValueAndValidity();
  }

  onCancel(): void {
    this.goBack();
  }

  goBack(): void {
    this.router.navigate(['/app/clients']);
  }
}