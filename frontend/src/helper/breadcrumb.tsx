import { Params } from "react-router";

type LabelGetter = (params: Readonly<Params<string>>) => string;

export const labels = new Map<string, string | LabelGetter>([
  ['check', 'Check access'],
  ['dashboard', 'Home'],
  ['signin', 'Sign In'],
  ['groups', 'Groups'],
  ['groups-create', 'Create'],
  ['groups-edit', 'Edit'],
  ['roles', 'Roles'],
  ['roles-create', 'Create'],
  ['roles-edit', 'Edit'],
  ['policies', 'Policies'],
  ['policies-create', 'Create'],
  ['policies-edit', 'Edit'],
  ['principals', 'Principals'],
  ['principals-create', 'Create'],
  ['principals-edit', 'Edit'],
  ['resources', 'Resources'],
  ['resources-create', 'Create'],
  ['resources-edit', 'Edit'],
  ['users', 'Users'],
  ['users-create', 'Create'],
  ['clients', 'Service accounts'],
  ['clients-create', 'Create'],
  ['clients-edit', 'Edit'],
]);


export type BreadcrumbItem = {
  label: string
  href: string
}

export const getBreadcrumbs = (
  matches: { id: string, pathname: string }[] | undefined,
  params: Readonly<Params<string>>
  ): BreadcrumbItem[] => {
    const items: BreadcrumbItem[] = [];
    
    matches?.forEach(match => {
      let label = labels.get(match.id);

      if (match.id.endsWith('-index')) {
        return;
      }

      if (typeof label === 'function') {
        label = label(params);
      }
      
      items.push({
        label: label || `Unlabelized: ${match.id}`,
        href: match.id !== 'dashboard' && match.pathname === '/' ? '' : match.pathname,
      });
    });
    
    return items;
  }
;