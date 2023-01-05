import { Params } from "react-router";

type LabelGetter = (params: Readonly<Params<string>>) => string;

export const labels = new Map<string, string | LabelGetter>([
  ['dashboard', 'Home'],
  ['signin', 'Sign In'],
  ['groups', 'Groups'],
  ['groups-create', 'Create'],
  ['groups-edit', 'Edit'],
  ['roles', 'Roles'],
  ['roles-create', 'Create'],
  ['roles-edit', 'Edit'],
  ['policies', 'Policies'],
  ['policies-edit', 'Edit'],
  ['users', 'Users'],
  ['users-edit', 'Edit'],
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