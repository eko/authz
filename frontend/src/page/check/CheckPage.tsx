import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import Check from 'page/check/component/Check';

export default function CheckPage() {
  return (
    <AuthenticatedLayout title='Want to check an access? You can do it here'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Check access
      </Typography>

      <Check />
   </AuthenticatedLayout>
  );
}