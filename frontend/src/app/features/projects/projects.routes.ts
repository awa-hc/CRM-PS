import { Routes } from '@angular/router';
import { ProjectsListComponent } from './projects-list/projects-list.component';
import { ProjectFormComponent } from './project-form/project-form.component';
import { ProjectDetailComponent } from './project-detail/project-detail.component';

export const projectsRoutes: Routes = [
  {
    path: '',
    component: ProjectsListComponent
  },
  {
    path: 'new',
    component: ProjectFormComponent
  },
  {
    path: ':id',
    component: ProjectDetailComponent
  },
  {
    path: ':id/edit',
    component: ProjectFormComponent
  }
];