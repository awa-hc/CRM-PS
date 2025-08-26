import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatSliderModule } from '@angular/material/slider';
import { ProjectService } from '../../../core/services/project.service';
import { ClientService } from '../../../core/services/client.service';
import { Project, CreateProjectRequest, UpdateProjectRequest } from '../../../core/models/project.model';
import { Client } from '../../../core/models/client.model';

@Component({
  selector: 'app-project-form',
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
    MatDatepickerModule,
    MatNativeDateModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatSliderModule
  ],
  templateUrl: './project-form.component.html',
  styleUrls: ['./project-form.component.css']
})
export class ProjectFormComponent implements OnInit {
  projectForm: FormGroup;
  clients: Client[] = [];
  loading = false;
  isEditMode = false;
  projectId?: number;

  constructor(
    private fb: FormBuilder,
    private projectService: ProjectService,
    private clientService: ClientService,
    private router: Router,
    private route: ActivatedRoute,
    private snackBar: MatSnackBar
  ) {
    this.projectForm = this.createForm();
  }

  ngOnInit(): void {
    this.loadClients();
    this.checkEditMode();
  }

  createForm(): FormGroup {
    return this.fb.group({
      code: [{ value: '', disabled: true }], // Auto-generated, read-only
      name: ['', [Validators.required]],
      description: [''],
      client_id: ['', [Validators.required]],
      status: ['planning', [Validators.required]],
      priority: ['medium', [Validators.required]],
      type: ['construction', [Validators.required]],
      address: ['', [Validators.required]],
      city: ['', [Validators.required]],
      state: ['', [Validators.required]],
      zipCode: ['', [Validators.required]],
      startDate: [''],
      endDate: [''],
      budget: ['', [Validators.required, Validators.min(0.01)]],
      estimatedCost: ['', [Validators.required, Validators.min(0.01)]],
      actualCost: [''],
      progress: [0, [Validators.min(0), Validators.max(100)]],
      notes: ['']
    });
  }

  loadClients(): void {
    this.clientService.getClients({ limit: 1000 }).subscribe({
      next: (response) => {
        this.clients = response.clients;
      },
      error: (error) => {
        console.error('Error loading clients:', error);
        this.snackBar.open('Error al cargar clientes', 'Cerrar', { duration: 3000 });
      }
    });
  }

  checkEditMode(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id && id !== 'new') {
      this.isEditMode = true;
      this.projectId = parseInt(id, 10);
      this.loadProject();
    }
  }

  loadProject(): void {
    if (!this.projectId) return;

    this.loading = true;
    this.projectService.getProject(this.projectId).subscribe({
      next: (project) => {
        console.log("project", project)
        this.populateForm(project);
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading project:', error);
        this.snackBar.open('Error al cargar proyecto', 'Cerrar', { duration: 3000 });
        this.loading = false;
        this.router.navigate(['/app/projects']);
      }
    });
  }

  populateForm(project: any): void {
    console.log('Datos del proyecto recibidos:', project);
    
    // Parse dates first
    const startDate = project.startDate ? this.parseDate(project.startDate) : null;
    const endDate = project.endDate ? this.parseDate(project.endDate) : null;

    console.log('Fechas parseadas:', {
      start_date_original: project.startDate,
      end_date_original: project.endDate,
      startDate_parsed: startDate,
      endDate_parsed: endDate
    });
    
    console.log('Campos específicos:', {
      zipCode: project.zipCode,
      estimatedCost: project.estimatedCost,
      actualCost: project.actualCost
    });
    const formData = {
      code: project.code || '', // Add code field mapping
      name: project.name || '',
      description: project.description || '',
      client_id: project.clientId, // Backend uses client_id
      status: project.status || 'planning',
      priority: project.priority || 'medium',
      type: project.project_type || 'construction', // Backend uses project_type
      address: project.address || '',
      city: project.city || '',
      state: project.state || '',
      zipCode: project.zipCode || '', // Backend uses zip_code
      startDate: startDate,
      endDate: endDate,
      budget: project.budget || 0,
      estimatedCost: project.estimatedCost || 0, // Backend uses estimated_cost
      actualCost: project.actualCost || 0, // Backend uses actual_cost
      progress: project.progress || 0,
      notes: project.notes || ''
    };
    
    console.log('Datos que se van a asignar al formulario:', formData);
    
    this.projectForm.patchValue(formData);
    
    // Verificar valores después del patchValue
    console.log('Valores del formulario después del patchValue:', {
      zipCode: this.projectForm.get('zipCode')?.value,
      startDate: this.projectForm.get('startDate')?.value,
      endDate: this.projectForm.get('endDate')?.value,
      estimatedCost: this.projectForm.get('estimatedCost')?.value
    });
  }

  private parseDate(dateString: string | null | undefined): Date | null {
    if (!dateString) return null;
    try {
      // Handle timezone-aware dates like "2025-08-08T00:00:00-03:00"
      const date = new Date(dateString);
      if (isNaN(date.getTime())) {
        console.warn('Invalid date string:', dateString);
        return null;
      }
      
      // For date inputs, we need to create a date without timezone offset
      // to avoid timezone conversion issues in the date picker
      const localDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());
      console.log('parseDate:', {
        original: dateString,
        parsed: date,
        localDate: localDate
      });
      
      return localDate;
    } catch (error) {
      console.error('Error parsing date:', dateString, error);
      return null;
    }
  }

  private formatDate(date: Date): string {
    if (!date) return '';
    return date.toISOString().split('T')[0];
  }

  private formatDateForBackend(date: Date): string {
    if (!date) return '';
    // Format as ISO string with time and timezone for backend
    return date.toISOString();
  }

  onSubmit(): void {
    if (this.projectForm.invalid) {
      this.markFormGroupTouched();
      return;
    }

    this.loading = true;
    const formValue = this.projectForm.value;

    // Format dates and map to CreateProjectRequest interface format (camelCase)
    const projectData = {
      name: this.projectForm.get('name')?.value,
      description: this.projectForm.get('description')?.value,
      client_id: this.projectForm.get('client_id')?.value,
      status: this.projectForm.get('status')?.value,
      priority: this.projectForm.get('priority')?.value,
      type: this.projectForm.get('type')?.value,
      address: this.projectForm.get('address')?.value,
      city: this.projectForm.get('city')?.value,
      state: this.projectForm.get('state')?.value,
      zipCode: this.projectForm.get('zipCode')?.value,
      startDate: this.formatDateForBackend(this.projectForm.get('startDate')?.value),
      endDate: this.formatDateForBackend(this.projectForm.get('endDate')?.value),
      budget: this.projectForm.get('budget')?.value,
      estimatedCost: this.projectForm.get('estimatedCost')?.value,
      notes: this.projectForm.get('notes')?.value
    };

    if (this.isEditMode && this.projectId) {
      this.updateProject(projectData);
    } else {
      this.createProject(projectData);
    }
  }

  createProject(projectData: CreateProjectRequest): void {
    this.projectService.createProject(projectData).subscribe({
      next: (project) => {
        this.snackBar.open('Proyecto creado exitosamente', 'Cerrar', { duration: 3000 });
        this.router.navigate(['/app/projects']);
      },
      error: (error) => {
        console.error('Error creating project:', error);
        this.snackBar.open('Error al crear proyecto', 'Cerrar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  updateProject(projectData: any): void {
    if (!this.projectId) return;

    // For updates, we need to map to snake_case for backend
    const updateData = {
      name: projectData.name,
      description: projectData.description,
      client_id: projectData.client_id,
      status: projectData.status,
      priority: projectData.priority,
      project_type: projectData.type,
      address: projectData.address,
      city: projectData.city,
      state: projectData.state,
      zip_code: projectData.zipCode,
      start_date: projectData.startDate,
      end_date: projectData.endDate,
      budget: projectData.budget,
      estimated_cost: projectData.estimatedCost,
      actual_cost: this.projectForm.get('actualCost')?.value,
      progress: this.projectForm.get('progress')?.value,
      notes: projectData.notes
    };

    this.projectService.updateProject(this.projectId, updateData).subscribe({
      next: (project) => {
        this.snackBar.open('Proyecto actualizado exitosamente', 'Cerrar', { duration: 3000 });
        this.router.navigate(['/app/projects']);
      },
      error: (error) => {
        console.error('Error updating project:', error);
        this.snackBar.open('Error al actualizar proyecto', 'Cerrar', { duration: 3000 });
        this.loading = false;
      }
    });
  }

  markFormGroupTouched(): void {
    Object.keys(this.projectForm.controls).forEach(key => {
      const control = this.projectForm.get(key);
      control?.markAsTouched();
    });
  }

  goBack(): void {
    this.router.navigate(['/app/projects']);
  }
}