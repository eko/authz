import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Cookies from 'universal-cookie';
import {
    createBrowserRouter,
    RouterProvider,
  } from 'react-router-dom';
import type { Router } from "@remix-run/router";
import { AuthContext, User } from 'context/auth';
import UserMenu from './UserMenu';

const initializeRouter = (
  user?: User,
  logout?: Function
) : Router => {
  if (logout === undefined) {
    logout = () => {}
  }

  return createBrowserRouter([
    {
      path: '/',
      element: (
        <AuthContext.Provider value={{
            user,
            logout,
        }}>
            <UserMenu />
        </AuthContext.Provider>
      ),
    }
  ]);
};

test('user menu: should not be visible when user not authenticated', async () => {
  // Given
  let router: Router | null = null;

  const Wrapper = () => {
    router = initializeRouter();

    return (
      <RouterProvider router={router} />
    )
  }

  // Render
  render(
      <Wrapper />
  );

  // When
  act(() => {
    router!.navigate('/');
  });

  // Then
  expect(screen.queryByText('Déconnexion')).toBeNull();
});

test('user menu: should be visible when user authenticated', async () => {
  // Given
  let router: Router | null = null;

  const logout = jest.fn();

  const Wrapper = () => {
    router = initializeRouter({
      username: 'John',
      token: 'xyz',
    }, logout);

    return (
      <RouterProvider router={router} />
    )
  }

  // Render
  render(
      <Wrapper />
  );

  // When
  act(() => {
    router!.navigate('/');
  });

  await userEvent.click(screen.getByRole('user-menu'));

  // Then
  expect(screen.queryByText('Déconnexion')).toBeVisible();
  expect(logout).toBeCalledTimes(0);
});

test('user menu: click on logout button', async () => {
  // Given
  let router: Router | null = null;

  const logout = jest.fn();

  const Wrapper = () => {
    router = initializeRouter({
      username: 'John',
      token: 'xyz',
    }, logout);

    return (
      <RouterProvider router={router} />
    )
  }

  // Render
  render(
      <Wrapper />
  );

  // When
  act(() => {
    router!.navigate('/');
  });

  await userEvent.click(screen.getByRole('user-menu'));
  await userEvent.click(screen.getByText('Déconnexion'));

  // Then
  expect(logout).toBeCalledTimes(1);
});