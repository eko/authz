import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import RoleCreateOrEdit from 'page/roles/component/RoleCreateOrEdit';

export default function RoleCreateOrEditPage() {
  return (
    <AuthenticatedLayout title='A role contains a set of policies'>
      <Breadcrumb />

      <RoleCreateOrEdit />
   </AuthenticatedLayout>
  );
}