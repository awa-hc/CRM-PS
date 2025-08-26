import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface DashboardStats {
  clients: {
    total: number;
    active: number;
  };
  projects: {
    total: number;
    active: number;
    completed: number;
  };
  quotes: {
    total: number;
    pending: number;
    accepted: number;
    accepted_value: number;
    total_value: number;
  };
  materials: {
    total: number;
    low_stock: number;
    inventory_value: number;
  };
  monthly: {
    clients: number;
    projects: number;
    quotes: number;
  };
}

export interface RecentActivity {
  id: string;
  type: 'client' | 'project' | 'quote';
  title: string;
  description: string;
  date: string;
  user: string;
}

export interface MonthlyData {
  month: string;
  revenue: number;
  projects: number;
  quotes: number;
}

export interface ProjectStatusData {
  status: string;
  count: number;
  percentage: number;
}

export interface QuoteStatusData {
  status: string;
  count: number;
  percentage: number;
}

@Injectable({
  providedIn: 'root'
})
export class DashboardService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.apiUrl}/dashboard`;

  /**
   * Obtiene las estadísticas generales del dashboard
   */
  getDashboardStats(): Observable<DashboardStats> {
    return this.http.get<DashboardStats>(`${this.apiUrl}/stats`);
  }

  /**
   * Obtiene la actividad reciente
   */
  getRecentActivity(limit: number = 10): Observable<RecentActivity[]> {
    return this.http.get<RecentActivity[]>(`${this.apiUrl}/recent-activity?limit=${limit}`);
  }

  /**
   * Obtiene datos mensuales para gráficos
   */
  getMonthlyData(months: number = 12): Observable<MonthlyData[]> {
    return this.http.get<MonthlyData[]>(`${this.apiUrl}/monthly-data?months=${months}`);
  }

  /**
   * Obtiene distribución de estados de proyectos
   */
  getProjectStatusDistribution(): Observable<ProjectStatusData[]> {
    return this.http.get<ProjectStatusData[]>(`${this.apiUrl}/project-status`);
  }

  /**
   * Obtiene distribución de estados de cotizaciones
   */
  getQuoteStatusDistribution(): Observable<QuoteStatusData[]> {
    return this.http.get<QuoteStatusData[]>(`${this.apiUrl}/quote-status`);
  }

  /**
   * Obtiene los próximos vencimientos
   */
  getUpcomingDeadlines(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/upcoming-deadlines`);
  }

  /**
   * Obtiene los clientes más activos
   */
  getTopClients(limit: number = 5): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/top-clients?limit=${limit}`);
  }

  /**
   * Obtiene los proyectos más rentables
   */
  getTopProjects(limit: number = 5): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/top-projects?limit=${limit}`);
  }
}