import { Breadcrumbs, Link, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import { useMatches, useParams } from 'react-router';
import { BreadcrumbItem, getBreadcrumbs } from 'helper/breadcrumb';

export default function Breadcrumb() {
  const [breadcrumbs, setBreadcrumbs] = useState<BreadcrumbItem[]>([]);
  const params = useParams();
  const matches = useMatches();

  useEffect(() => {
    const items = getBreadcrumbs(matches, params);
    setBreadcrumbs(items);
  }, [matches, params]);

  return (
    <Breadcrumbs aria-label='breadcrumb'>
      {breadcrumbs.map((breadcrumb, index) => {
        return breadcrumbs.length - 1 === index || breadcrumb.href === '' ? (
          <Typography key={index} color="text.primary">{breadcrumb.label}</Typography>
        ) : (
          <Link key={index} underline='hover' color='inherit' href={breadcrumb.href}>
            {breadcrumb.label}
          </Link>
        )
      })}
    </Breadcrumbs>
  );
}