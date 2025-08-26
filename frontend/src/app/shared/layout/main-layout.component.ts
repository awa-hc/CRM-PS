import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterOutlet, RouterLink } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatDividerModule } from '@angular/material/divider';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { AuthService } from '../../core/services/auth.service';
import { User } from '../../core/models/user.model';

interface NavigationItem {
  label: string;
  icon: string;
  route: string;
  roles?: string[];
}

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    RouterLink,
    MatToolbarModule,
    MatSidenavModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatMenuModule,
    MatDividerModule
  ],
  templateUrl: './main-layout.component.html',
  styleUrls: ['./main-layout.component.css']
})
export class MainLayoutComponent {
  private readonly breakpointObserver = inject(BreakpointObserver);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  public readonly isHandset = signal(false);
  public readonly currentUser = signal<User | null>(null);

  private readonly navigationItems: NavigationItem[] = [
    {
      label: 'Dashboard',
      icon: 'dashboard',
      route: '/app/dashboard'
    },
    {
      label: 'Clientes',
      icon: 'people',
      route: '/app/clients'
    },
    {
      label: 'Proyectos',
      icon: 'construction',
      route: '/app/projects'
    },
    {
      label: 'Cotizaciones',
      icon: 'description',
      route: '/app/quotes'
    },
    {
      label: 'Reportes',
      icon: 'analytics',
      route: '/app/reports',
      roles: ['admin', 'manager']
    },
    {
      label: 'Configuración',
      icon: 'settings',
      route: '/app/settings',
      roles: ['admin']
    }
  ];

  constructor() {
    this.breakpointObserver.observe(Breakpoints.Handset)
      .subscribe(result => {
        this.isHandset.set(result.matches);
      });

    this.currentUser.set(this.authService.currentUser());
  }

  getVisibleNavigationItems(): NavigationItem[] {
    const user = this.currentUser();
    if (!user) return [];

    return this.navigationItems.filter(item => {
      if (!item.roles) return true;
      return item.roles.includes(user.role);
    });
  }

  viewProfile(): void {
    this.router.navigate(['/app/profile']);
  }

  changePassword(): void {
    this.router.navigate(['/app/change-password']);
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/auth/login']);
  }
}