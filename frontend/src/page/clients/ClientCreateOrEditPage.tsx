import AuthenticatedLayout from 'layout/AuthenticatedLayout';
import Breadcrumb from 'layout/Breadcrumb';
import ClientCreateOrEdit from 'page/clients/component/ClientCreateOrEdit';

export default function ClientCreateOrEditPage() {
  return (
    <AuthenticatedLayout title='A service account can be used to authenticate your applications'>
      <Breadcrumb />

      <ClientCreateOrEdit />
   </AuthenticatedLayout>
  );
}