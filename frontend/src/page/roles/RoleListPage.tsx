import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import RoleList from 'page/roles/component/RoleList';

export default function RoleListPage() {
  return (
    <AuthenticatedLayout title='Manage your application roles'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Roles
      </Typography>

      <RoleList />
   </AuthenticatedLayout>
  );
}