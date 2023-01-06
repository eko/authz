import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import PrincipalCreateOrEdit from 'page/principals/component/PrincipalCreateOrEdit';

export default function PrincipalCreateOrEditPage() {
  return (
    <AuthenticatedLayout title='A principal contains a set of roles and attributes'>
      <Breadcrumb />

      <PrincipalCreateOrEdit />
   </AuthenticatedLayout>
  );
}