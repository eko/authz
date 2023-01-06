import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import ResourceCreateOrEdit from 'page/resources/component/ResourceCreateOrEdit';

export default function ResourceCreateOrEditPage() {
  return (
    <AuthenticatedLayout title='A resource contains a set of policies'>
      <Breadcrumb />

      <ResourceCreateOrEdit />
   </AuthenticatedLayout>
  );
}