import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import PolicyList from 'page/policies/component/PolicyList';

export default function PolicyListPage() {
  return (
    <AuthenticatedLayout title='Manage your application policies'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Policies
      </Typography>

      <PolicyList />
   </AuthenticatedLayout>
  );
}