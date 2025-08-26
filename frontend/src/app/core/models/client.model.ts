export interface Client {
  id: number;
  name: string;
  email: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  company?: string;
  taxId?: string;
  contactType: 'individual' | 'company';
  notes?: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateClientRequest {
  name: string;
  email: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  company?: string;
  taxId?: string;
  contactType: 'individual' | 'company';
  notes?: string;
}

export interface UpdateClientRequest {
  name?: string;
  email?: string;
  phone?: string;
  address?: string;
  city?: string;
  state?: string;
  zipCode?: string;
  company?: string;
  taxId?: string;
  contactType?: 'individual' | 'company';
  notes?: string;
  isActive?: boolean;
}

export interface ClientStats {
  total: number;
  active: number;
  inactive: number;
  individual: number;
  company: number;
  recent: number;
}