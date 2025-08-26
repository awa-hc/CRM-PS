import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
  Quote,
  CreateQuoteRequest,
  UpdateQuoteRequest,
  QuoteStats,
  QuoteFilters
} from '../models/quote.model';

@Injectable({
  providedIn: 'root'
})
export class QuoteService {
  private apiUrl = `${environment.apiUrl}/quotes/`;

  constructor(private http: HttpClient) {}

  getQuotes(filters?: QuoteFilters, page: number = 1, limit: number = 10): Observable<{
    quotes: Quote[];
    total: number;
    page: number;
    totalPages: number;
  }> {
    let params = new HttpParams()
      .set('page', page.toString())
      .set('limit', limit.toString());

    if (filters) {
      if (filters.search) {
        params = params.set('search', filters.search);
      }
      if (filters.clientId) {
        params = params.set('client_id', filters.clientId.toString());
      }
      if (filters.projectId) {
        params = params.set('project_id', filters.projectId.toString());
      }
      if (filters.status) {
        params = params.set('status', filters.status);
      }
      if (filters.dateFrom) {
        params = params.set('date_from', filters.dateFrom);
      }
      if (filters.dateTo) {
        params = params.set('date_to', filters.dateTo);
      }
    }

    return this.http.get<{
      quotes: Quote[];
      total: number;
      page: number;
      totalPages: number;
    }>(this.apiUrl, { params });
  }

  getQuote(id: number): Observable<Quote> {
    return this.http.get<Quote>(`${this.apiUrl}/${id}`);
  }

  createQuote(quote: CreateQuoteRequest): Observable<Quote> {
    return this.http.post<Quote>(this.apiUrl, quote);
  }

  updateQuote(id: number, quote: UpdateQuoteRequest): Observable<Quote> {
    return this.http.put<Quote>(`${this.apiUrl}/${id}`, quote);
  }

  deleteQuote(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }

  getQuoteStats(): Observable<QuoteStats> {
    return this.http.get<QuoteStats>(`${this.apiUrl}/stats`);
  }

  sendQuote(id: number): Observable<Quote> {
    return this.http.post<Quote>(`${this.apiUrl}/${id}/send`, {});
  }

  acceptQuote(id: number): Observable<Quote> {
    return this.http.post<Quote>(`${this.apiUrl}/${id}/accept`, {});
  }

  rejectQuote(id: number): Observable<Quote> {
    return this.http.post<Quote>(`${this.apiUrl}/${id}/reject`, {});
  }

  duplicateQuote(id: number): Observable<Quote> {
    return this.http.post<Quote>(`${this.apiUrl}/${id}/duplicate`, {});
  }

  generatePDF(id: number): Observable<Blob> {
    return this.http.get(`${this.apiUrl}/${id}/pdf`, {
      responseType: 'blob'
    });
  }

  getQuotesByClient(clientId: number): Observable<Quote[]> {
    return this.http.get<Quote[]>(`${this.apiUrl}/client/${clientId}`);
  }

  getQuotesByProject(projectId: number): Observable<Quote[]> {
    return this.http.get<Quote[]>(`${this.apiUrl}/project/${projectId}`);
  }
}