import { RouteObject } from "react-router";

import DashboardPage from "page/dashboard/DashboardPage";
import SigninPage from "page/SigninPage";
import NotFoundPage from "page/NotFoundPage";
import RoleListPage from "page/roles/RoleListPage";
import RoleCreateOrEditPage from "page/roles/RoleCreateOrEditPage";
import PolicyListPage from "page/policies/PolicyListPage";
import PolicyCreateOrEditPage from "page/policies/PolicyCreateOrEditPage";
import ResourceListPage from "page/resources/ResourceListPage";
import ResourceCreateOrEditPage from "page/resources/ResourceCreateOrEditPage";
import PrincipalListPage from "page/principals/PrincipalListPage";
import PrincipalCreateOrEditPage from "page/principals/PrincipalCreateOrEditPage";
import UserListPage from "page/users/UserListPage";
import UserCreateOrEditPage from "page/users/UserCreateOrEditPage";
import ClientCreateOrEditPage from "page/clients/ClientCreateOrEditPage";
import ClientListPage from "page/clients/ClientListPage";
import CheckPage from "page/check/CheckPage";

export const routes: RouteObject[] = [
  {
    id: 'dashboard',
    errorElement: <NotFoundPage />,
    children: [
      {
        id: 'dashboard-index',
        path: '/',
        index: true,
        element: <DashboardPage />,
      },
      {
        id: 'signin',
        path: "/signin",
        element: <SigninPage />,
      },
      {
        id: 'check',
        path: '/check',
        children: [
          {
            id: 'check-index',
            path: '/check',
            index: true,
            element: <CheckPage />,
          },
        ],
      },
      {
        id: 'groups',
        path: '/groups',
        children: [
          {
            id: 'groups-index',
            path: '/groups',
            index: true,
            element: <NotFoundPage />,
          },
          {
            id: 'groups-create',
            path: '/groups/create',
            element: <NotFoundPage />,
          },
          {
            id: 'groups-edit',
            path: '/groups/edit/:id',
            element: <NotFoundPage />,
          },
        ],
      },
      {
        id: 'roles',
        path: '/roles',
        children: [
          {
            id: 'roles-index',
            path: '/roles',
            index: true,
            element: <RoleListPage />,
          },
          {
            id: 'roles-create',
            path: '/roles/create',
            element: <RoleCreateOrEditPage />,
          },
          {
            id: 'roles-edit',
            path: '/roles/edit/:id',
            element: <RoleCreateOrEditPage />,
          },
        ],
      },
      {
        id: 'policies',
        path: '/policies',
        children: [
          {
            id: 'policies-index',
            path: '/policies',
            index: true,
            element: <PolicyListPage />,
          },
          {
            id: 'policies-create',
            path: '/policies/create',
            element: <PolicyCreateOrEditPage />,
          },
          {
            id: 'policies-edit',
            path: '/policies/edit/:id',
            element: <PolicyCreateOrEditPage />,
          },
        ],
      },
      {
        id: 'resources',
        path: '/resources',
        children: [
          {
            id: 'resources-index',
            path: '/resources',
            index: true,
            element: <ResourceListPage />,
          },
          {
            id: 'resources-create',
            path: '/resources/create',
            element: <ResourceCreateOrEditPage />,
          },
          {
            id: 'resources-edit',
            path: '/resources/edit/:id',
            element: <ResourceCreateOrEditPage />,
          },
        ],
      },
      {
        id: 'principals',
        path: '/principals',
        children: [
          {
            id: 'principals-index',
            path: '/principals',
            index: true,
            element: <PrincipalListPage />,
          },
          {
            id: 'principals-create',
            path: '/principals/create',
            element: <PrincipalCreateOrEditPage />,
          },
          {
            id: 'principals-edit',
            path: '/principals/edit/:id',
            element: <PrincipalCreateOrEditPage />,
          },
        ],
      },
      {
        id: 'clients',
        path: '/clients',
        children: [
          {
            id: 'clients-index',
            path: '/clients',
            index: true,
            element: <ClientListPage />,
          },
          {
            id: 'clients-create',
            path: '/clients/create',
            element: <ClientCreateOrEditPage />,
          },
          {
            id: 'clients-edit',
            path: '/clients/edit/:id',
            element: <ClientCreateOrEditPage />,
          },
        ],
      },
      {
        id: 'users',
        path: '/users',
        children: [
          {
            id: 'users-index',
            path: '/users',
            index: true,
            element: <UserListPage />,
          },
          {
            id: 'users-create',
            path: '/users/create',
            element: <UserCreateOrEditPage />,
          },
        ],
      },
    ],
  },
];
