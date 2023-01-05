import { RouteObject } from "react-router";

import DashboardPage from "page/dashboard/DashboardPage";
import SigninPage from "page/SigninPage";
import NotFoundPage from "page/NotFoundPage";
import RoleListPage from "page/roles/RoleListPage";
import RoleCreateOrEditPage from "page/roles/RoleCreateOrEditPage";
import PolicyListPage from "page/policies/PolicyListPage";
import PolicyCreateOrEditPage from "page/policies/PolicyCreateOrEditPage";

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
            element: <NotFoundPage />,
          },
          {
            id: 'resources-create',
            path: '/resources/create',
            element: <NotFoundPage />,
          },
          {
            id: 'resources-edit',
            path: '/resources/edit/:id',
            element: <NotFoundPage />,
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
            element: <NotFoundPage />,
          },
          {
            id: 'principals-create',
            path: '/principals/create',
            element: <NotFoundPage />,
          },
          {
            id: 'principals-edit',
            path: '/principals/edit/:id',
            element: <NotFoundPage />,
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
            element: <NotFoundPage />,
          },
          {
            id: 'clients-create',
            path: '/clients/create',
            element: <NotFoundPage />,
          },
          {
            id: 'clients-edit',
            path: '/clients/edit/:id',
            element: <NotFoundPage />,
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
            element: <NotFoundPage />,
          },
          {
            id: 'users-edit',
            path: '/users/edit/:id',
            element: <NotFoundPage />,
          },
        ],
      },
    ],
  },
];
