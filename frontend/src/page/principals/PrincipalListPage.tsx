import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import PrincipalList from 'page/principals/component/PrincipalList';

export default function PrincipalListPage() {
  return (
    <AuthenticatedLayout title='Manage your principals'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Principals
      </Typography>

      <PrincipalList />
   </AuthenticatedLayout>
  );
}