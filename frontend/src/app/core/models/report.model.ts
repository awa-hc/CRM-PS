export interface DashboardStats {
  totalClients: number;
  totalProjects: number;
  totalQuotes: number;
  totalRevenue: number;
  activeProjects: number;
  completedProjects: number;
  pendingQuotes: number;
  monthlyRevenue: number;
}

export interface RecentActivity {
  id: number;
  type: 'client' | 'project' | 'quote';
  title: string;
  description: string;
  date: string;
  status?: string;
}

export interface ClientReport {
  totalClients: number;
  newClientsThisMonth: number;
  activeClients: number;
  averageProjectValue: number;
  topClients: TopClient[];
  clientPerformance: ClientPerformance[];
  clientGrowth: MonthlyData[];
  clientsByType: ClientTypeData[];
  recentActivity: ClientActivity[];
  geographicDistribution: GeographicData[];
}

export interface ClientPerformance {
  id: number;
  name: string;
  totalProjects: number;
  projectCount: number;
  totalValue: number;
  averageProjectValue: number;
  completedProjects: number;
  completionRate: number;
  lastProjectDate: string;
  status: string;
}

export interface TopClient {
  id: number;
  clientName: string;
  projectsCount: number;
  totalRevenue: number;
  completionRate: number;
}

export interface ClientTypeData {
  type: string;
  count: number;
  totalRevenue: number;
  percentage: number;
}

export interface ClientActivity {
  id: number;
  type: string;
  clientName: string;
  description: string;
  date: string;
  value?: number;
}

export interface GeographicData {
  location: string;
  count: number;
  percentage: number;
  totalRevenue: number;
}

export interface ProjectReport {
  totalProjects: number;
  activeProjects: number;
  completedProjects: number;
  averageProjectDuration: number;
  projectsByStatus: StatusData[];
  projectsByType: TypeData[];
  monthlyProjects: MonthlyData[];
  budgetVariance: BudgetVariance[];
}

export interface QuoteReport {
  totalQuotes: number;
  pendingQuotes: number;
  acceptedQuotes: number;
  rejectedQuotes: number;
  conversionRate: number;
  averageQuoteValue: number;
  quotesByStatus: StatusData[];
  monthlyQuotes: MonthlyData[];
}

export interface FinancialReport {
  totalRevenue: number;
  monthlyRevenue: number;
  yearlyRevenue: number;
  profitMargin: number;
  averageProjectValue: number;
  revenueByMonth: MonthlyData[];
  revenueByClient: ClientRevenue[];
  expensesByCategory: ExpenseData[];
}

export interface MonthlyData {
  month: string;
  value: number;
  label?: string;
}

export interface StatusData {
  status: string;
  count: number;
  percentage: number;
}

export interface TypeData {
  type: string;
  count: number;
  percentage: number;
}

export interface BudgetVariance {
  projectId: number;
  projectName: string;
  budget: number;
  actualCost: number;
  variance: number;
  variancePercentage: number;
}

export interface ClientRevenue {
  clientId: number;
  clientName: string;
  revenue: number;
  projectCount: number;
}

export interface ExpenseData {
  category: string;
  amount: number;
  percentage: number;
}

export interface ReportFilters {
  startDate?: string;
  endDate?: string;
  clientId?: number;
  projectType?: string;
  status?: string;
}