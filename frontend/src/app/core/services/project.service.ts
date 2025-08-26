import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Project, CreateProjectRequest, UpdateProjectRequest, ProjectStats } from '../models/project.model';
import { environment } from '../../../environments/environment';

export interface ProjectFilters {
  search?: string;
  status?: string;
  priority?: string;
  type?: string;
  client_id?: number;
  page?: number;
  limit?: number;
}

export interface ProjectsResponse {
  projects: Project[];
  total: number;
  page: number;
  totalPages: number;
}

@Injectable({
  providedIn: 'root'
})
export class ProjectService {
  private apiUrl = `${environment.apiUrl}/projects`;

  constructor(private http: HttpClient) {}

  getProjects(filters: ProjectFilters = {}): Observable<ProjectsResponse> {
    let params = new HttpParams();
    
    if (filters.search) {
      params = params.set('search', filters.search);
    }
    if (filters.status) {
      params = params.set('status', filters.status);
    }
    if (filters.priority) {
      params = params.set('priority', filters.priority);
    }
    if (filters.type) {
      params = params.set('type', filters.type);
    }
    if (filters.client_id) {
      params = params.set('client_id', filters.client_id.toString());
    }
    if (filters.page) {
      params = params.set('page', filters.page.toString());
    }
    if (filters.limit) {
      params = params.set('limit', filters.limit.toString());
    }

    return this.http.get<ProjectsResponse>(this.apiUrl, { params });
  }

  getProject(id: number): Observable<Project> {
    return this.http.get<Project>(`${this.apiUrl}/${id}`);
  }

  createProject(project: CreateProjectRequest): Observable<Project> {
    return this.http.post<Project>(this.apiUrl, project);
  }

  updateProject(id: number, project: UpdateProjectRequest): Observable<Project> {
    return this.http.put<Project>(`${this.apiUrl}/${id}`, project);
  }

  deleteProject(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  getProjectStats(): Observable<ProjectStats> {
    return this.http.get<ProjectStats>(`${this.apiUrl}/stats`);
  }

  updateProjectProgress(id: number, progress: number): Observable<Project> {
    return this.http.patch<Project>(`${this.apiUrl}/${id}/progress`, { progress });
  }

  getProjectsByClient(clientId: number): Observable<Project[]> {
    return this.http.get<Project[]>(`${this.apiUrl}/client/${clientId}`);
  }
}