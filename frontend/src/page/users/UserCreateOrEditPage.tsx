import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import UserCreateOrEdit from 'page/users/component/UserCreateOrEdit';

export default function UserCreateOrEditPage() {
  return (
    <AuthenticatedLayout title='A user is used to authenticate on Authz'>
      <Breadcrumb />

      <UserCreateOrEdit />
   </AuthenticatedLayout>
  );
}