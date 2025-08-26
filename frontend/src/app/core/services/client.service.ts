import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Client, CreateClientRequest, UpdateClientRequest, ClientStats } from '../models/client.model';

// Re-export Client for use in other components
export type { Client } from '../models/client.model';
import { environment } from '../../../environments/environment';

export interface ClientsResponse {
  clients: Client[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ClientFilters {
  search?: string;
  contactType?: 'individual' | 'company';
  isActive?: boolean;
  page?: number;
  limit?: number;
}

@Injectable({
  providedIn: 'root'
})
export class ClientService {
  private readonly apiUrl = `${environment.apiUrl}/clients`;

  constructor(private http: HttpClient) {}

  getClients(filters?: ClientFilters): Observable<ClientsResponse> {
    let params = new HttpParams();
    
    if (filters) {
      if (filters.search) {
        params = params.set('search', filters.search);
      }
      if (filters.contactType) {
        params = params.set('contactType', filters.contactType);
      }
      if (filters.isActive !== undefined) {
        params = params.set('isActive', filters.isActive.toString());
      }
      if (filters.page) {
        params = params.set('page', filters.page.toString());
      }
      if (filters.limit) {
        params = params.set('limit', filters.limit.toString());
      }
    }

    return this.http.get<ClientsResponse>(this.apiUrl, { params });
  }

  getClient(id: number): Observable<Client> {
    return this.http.get<Client>(`${this.apiUrl}/${id}`);
  }

  createClient(client: CreateClientRequest): Observable<Client> {
    return this.http.post<Client>(this.apiUrl, client);
  }

  updateClient(id: number, client: UpdateClientRequest): Observable<Client> {
    return this.http.put<Client>(`${this.apiUrl}/${id}`, client);
  }

  deleteClient(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  getClientStats(): Observable<ClientStats> {
    return this.http.get<ClientStats>(`${this.apiUrl}/stats`);
  }
}