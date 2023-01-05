import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import PolicyCreateOrEdit from 'page/policies/component/PolicyCreateOrEdit';

export default function PolicyCreateOrEditPage() {
  return (
    <AuthenticatedLayout title='A policy contains a set of resources, actions and eventually attribute rules'>
      <Breadcrumb />

      <PolicyCreateOrEdit />
   </AuthenticatedLayout>
  );
}