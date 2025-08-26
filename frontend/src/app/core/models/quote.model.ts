export interface Quote {
  id: number;
  quoteNumber: string;
  clientId: number;
  client?: {
    id: number;
    name: string;
    email: string;
  };
  projectId?: number;
  project?: {
    id: number;
    name: string;
    code: string;
  };
  title: string;
  description: string;
  status: 'draft' | 'sent' | 'accepted' | 'rejected' | 'expired';
  validUntil?: string;
  subtotal: number;
  taxRate: number;
  taxAmount: number;
  discount: number;
  total: number;
  notes: string;
  terms: string;
  createdAt: string;
  updatedAt: string;
  items?: QuoteItem[];
}

export interface QuoteItem {
  id: number;
  quoteId: number;
  description: string;
  quantity: number;
  unit: string;
  unitPrice: number;
  total: number;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateQuoteRequest {
  clientId: number;
  projectId?: number;
  title: string;
  description: string;
  validUntil?: string;
  taxRate: number;
  discount: number;
  notes: string;
  terms: string;
  items: CreateQuoteItemRequest[];
}

export interface CreateQuoteItemRequest {
  description: string;
  quantity: number;
  unit: string;
  unitPrice: number;
  notes: string;
}

export interface UpdateQuoteRequest {
  title?: string;
  description?: string;
  status?: 'draft' | 'sent' | 'accepted' | 'rejected' | 'expired';
  validUntil?: string;
  taxRate?: number;
  discount?: number;
  notes?: string;
  terms?: string;
}

export interface QuoteStats {
  totalQuotes: number;
  draftQuotes: number;
  sentQuotes: number;
  acceptedQuotes: number;
  rejectedQuotes: number;
  totalValue: number;
  acceptedValue: number;
}

export interface QuoteFilters {
  search?: string;
  clientId?: number;
  projectId?: number;
  status?: string;
  dateFrom?: string;
  dateTo?: string;
}