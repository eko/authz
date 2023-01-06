import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import ResourceList from 'page/resources/component/ResourceList';

export default function ResourceListPage() {
  return (
    <AuthenticatedLayout title='Manage your application resources'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Resources
      </Typography>

      <ResourceList />
   </AuthenticatedLayout>
  );
}