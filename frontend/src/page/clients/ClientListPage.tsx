import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import ClientList from 'page/clients/component/ClientList';

export default function ClientListPage() {
  return (
    <AuthenticatedLayout title='Manage your clients'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Service accounts
      </Typography>

      <ClientList />
   </AuthenticatedLayout>
  );
}