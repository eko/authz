import { Typography } from '@mui/material';
import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import UserList from 'page/users/component/UserList';

export default function UserListPage() {
  return (
    <AuthenticatedLayout title='Manage your users'>
      <Breadcrumb />

      <Typography variant="h3" gutterBottom marginTop={1} marginBottom={2}>
          Users
      </Typography>

      <UserList />
   </AuthenticatedLayout>
  );
}