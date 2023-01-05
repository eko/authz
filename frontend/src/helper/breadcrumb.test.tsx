import { BreadcrumbItem, getBreadcrumbs } from './breadcrumb';

test('breadcrumb: compute single', async () => {
  // When
  const matches = [
    { id: 'dashboard', pathname: '/' },
  ];

  const breadcrumbs = getBreadcrumbs(matches, {});

  // Then
  const expected: BreadcrumbItem[] = [
    { label: 'Home', href: '/' },
  ];

  expect(breadcrumbs).toEqual(expected);
});


test('breadcrumb: compute nested', async () => {
  // When
  const matches = [
    { id: 'dashboard', pathname: '/' },
    { id: 'groups', pathname: '/groups' },
  ];

  const breadcrumbs = getBreadcrumbs(matches, {});

  // Then
  const expected: BreadcrumbItem[] = [
    { label: 'Home', href: '/' },
    { label: 'Groups', href: '/groups' },
  ];

  expect(breadcrumbs).toEqual(expected);
});