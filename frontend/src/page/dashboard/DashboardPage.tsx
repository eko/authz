import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import Dashboard from 'page/dashboard/component/Dashboard';


export default function DashboardPage() {
  return (
    <AuthenticatedLayout title='What happens on your authorization backend'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Dashboard
      </Typography>

      <Dashboard />
   </AuthenticatedLayout>
  );
}