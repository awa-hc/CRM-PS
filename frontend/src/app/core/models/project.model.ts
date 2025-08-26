export interface Project {
  id: number;
  code: string;
  name: string;
  description?: string;
  client_id: number;
  clientName?: string;
  status: 'planning' | 'active' | 'completed' | 'cancelled' | 'on_hold';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  type: 'construction' | 'renovation' | 'maintenance';
  address: string;
  city: string;
  state: string;
  zipCode: string;
  startDate?: string;
  endDate?: string;
  estimatedEndDate?: string;
  budget: number;
  estimatedCost: number;
  actualCost: number;
  progress: number;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
  client_id: number;
  status: 'planning' | 'active' | 'completed' | 'cancelled' | 'on_hold';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  type: 'construction' | 'renovation' | 'maintenance';
  address: string;
  city: string;
  state: string;
  zipCode: string;
  startDate?: string;
  endDate?: string;
  estimatedEndDate?: string;
  budget: number;
  estimatedCost: number;
  notes?: string;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  client_id?: number;
  status?: 'planning' | 'active' | 'completed' | 'cancelled' | 'on_hold';
  priority?: 'low' | 'medium' | 'high' | 'urgent';
  type?: 'construction' | 'renovation' | 'maintenance';
  address?: string;
  city?: string;
  state?: string;
  zipCode?: string;
  startDate?: string;
  endDate?: string;
  estimatedEndDate?: string;
  budget?: number;
  estimatedCost?: number;
  actualCost?: number;
  progress?: number;
  notes?: string;
}

export interface ProjectStats {
  total: number;
  active: number;
  completed: number;
  planning: number;
  lowPriority: number;
  mediumPriority: number;
  highPriority: number;
  urgentPriority: number;
  construction: number;
  renovation: number;
  maintenance: number;
  totalBudget: number;
  totalActualCost: number;
  averageProgress: number;
}