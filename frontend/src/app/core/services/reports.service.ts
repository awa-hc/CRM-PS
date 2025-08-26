import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface ReportFilters {
  startDate?: string;
  endDate?: string;
  clientId?: number;
  projectId?: number;
  status?: string;
  clientType?: string;
}

export interface FinancialReport {
  totalRevenue: number;
  totalExpenses: number;
  profit: number;
  profitMargin: number;
  quotesCount: number;
  acceptedQuotes: number;
  rejectedQuotes: number;
  pendingQuotes: number;
  averageQuoteValue: number;
  monthlyRevenue: { month: string; revenue: number }[];
  expensesByCategory: { category: string; amount: number }[];
}

export interface ProjectReport {
  totalProjects: number;
  activeProjects: number;
  completedProjects: number;
  onHoldProjects: number;
  cancelledProjects: number;
  averageProjectDuration: number;
  averageCompletionRate: number;
  totalBudget: number;
  totalSpent: number;
  projectsByStatus: { status: string; count: number }[];
  projectsByPriority: { priority: string; count: number }[];
  budgetUtilization: { projectId: number; projectName: string; budgetUsed: number; totalBudget: number }[];
  delayedProjects: { projectId: number; projectName: string; daysDelayed: number }[];
  topClients: { clientId: number; clientName: string; projectsCount: number; totalValue: number }[];
  recentProjects: { projectId: number; projectName: string; clientName: string; status: string; startDate: string }[];
}

export interface ClientReport {
  totalClients: number;
  activeClients: number;
  inactiveClients: number;
  newClientsThisMonth: number;
  totalRevenue: number;
  averageProjectValue: number;
  clientsByType: { type: string; count: number; totalRevenue: number }[];
  topClientsByRevenue: { clientId: number; clientName: string; totalRevenue: number }[];
  topClients: { clientId: number; clientName: string; totalRevenue: number; projectsCount: number }[];
  clientActivity: { clientId: number; clientName: string; lastActivity: string; projectsCount: number; quotesCount: number }[];
  clientPerformance: { name: string; projectCount: number; totalValue: number; averageProjectValue: number; completedProjects: number; completionRate: number; lastProject: string }[];
  recentActivity: { type: string; clientId: number; clientName: string; description: string; date: string }[];
  geographicDistribution: { location: string; count: number; percentage: number }[];
}

export interface DashboardStats {
  totalClients: number;
  totalProjects: number;
  totalQuotes: number;
  monthlyRevenue: number;
  recentActivity: {
    type: 'client' | 'project' | 'quote';
    id: number;
    title: string;
    date: string;
    status?: string;
  }[];
  upcomingDeadlines: {
    type: 'project' | 'quote';
    id: number;
    title: string;
    deadline: string;
    daysLeft: number;
  }[];
}

@Injectable({
  providedIn: 'root'
})
export class ReportsService {
  private apiUrl = `${environment.apiUrl}/reports`;

  constructor(private http: HttpClient) {}

  // Financial Reports
  getFinancialReport(filters?: ReportFilters): Observable<FinancialReport> {
    let params = new HttpParams();
    if (filters) {
      if (filters.startDate) params = params.set('startDate', filters.startDate);
      if (filters.endDate) params = params.set('endDate', filters.endDate);
      if (filters.clientId) params = params.set('clientId', filters.clientId.toString());
      if (filters.status) params = params.set('status', filters.status);
    }
    return this.http.get<FinancialReport>(`${this.apiUrl}/financial`, { params });
  }

  // Project Reports
  getProjectReport(filters?: ReportFilters): Observable<ProjectReport> {
    let params = new HttpParams();
    if (filters) {
      if (filters.startDate) params = params.set('startDate', filters.startDate);
      if (filters.endDate) params = params.set('endDate', filters.endDate);
      if (filters.clientId) params = params.set('clientId', filters.clientId.toString());
      if (filters.status) params = params.set('status', filters.status);
    }
    return this.http.get<ProjectReport>(`${this.apiUrl}/projects`, { params });
  }

  // Client Reports
  getClientReport(filters?: ReportFilters): Observable<ClientReport> {
    let params = new HttpParams();
    if (filters) {
      if (filters.startDate) params = params.set('startDate', filters.startDate);
      if (filters.endDate) params = params.set('endDate', filters.endDate);
    }
    return this.http.get<ClientReport>(`${this.apiUrl}/clients`, { params });
  }

  // Dashboard Stats
  getDashboardStats(): Observable<DashboardStats> {
    return this.http.get<DashboardStats>(`${this.apiUrl}/dashboard`);
  }

  // Export Reports
  exportFinancialReport(filters?: ReportFilters, format: 'pdf' | 'excel' = 'pdf'): Observable<Blob> {
    let params = new HttpParams().set('format', format);
    if (filters) {
      if (filters.startDate) params = params.set('startDate', filters.startDate);
      if (filters.endDate) params = params.set('endDate', filters.endDate);
      if (filters.clientId) params = params.set('clientId', filters.clientId.toString());
      if (filters.status) params = params.set('status', filters.status);
    }
    return this.http.get(`${this.apiUrl}/financial/export`, {
      params,
      responseType: 'blob'
    });
  }

  exportProjectReport(filters?: ReportFilters, format: 'pdf' | 'excel' = 'pdf'): Observable<Blob> {
    let params = new HttpParams().set('format', format);
    if (filters) {
      if (filters.startDate) params = params.set('startDate', filters.startDate);
      if (filters.endDate) params = params.set('endDate', filters.endDate);
      if (filters.clientId) params = params.set('clientId', filters.clientId.toString());
      if (filters.status) params = params.set('status', filters.status);
    }
    return this.http.get(`${this.apiUrl}/projects/export`, {
      params,
      responseType: 'blob'
    });
  }

  exportClientReport(filters?: ReportFilters, format: 'pdf' | 'excel' = 'pdf'): Observable<Blob> {
    let params = new HttpParams().set('format', format);
    if (filters) {
      if (filters.startDate) params = params.set('startDate', filters.startDate);
      if (filters.endDate) params = params.set('endDate', filters.endDate);
    }
    return this.http.get(`${this.apiUrl}/clients/export`, {
      params,
      responseType: 'blob'
    });
  }
}