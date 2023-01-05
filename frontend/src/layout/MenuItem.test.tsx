import { List, ThemeProvider } from '@mui/material';
import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import {
    createBrowserRouter,
    RouterProvider,
  } from 'react-router-dom';
import type { Router } from "@remix-run/router";
import MenuItem from 'layout/MenuItem';
import theme from 'layout/theme';

const initializeRouter = (nestedOpened: boolean, setNestedOpened: Function) : Router => createBrowserRouter([
  {
    path: '/',
    element: (
      <ThemeProvider theme={theme}>
        <MenuItem
          label='My simple button'
          path='/my/simple/button'
          menuOpened={true}
        />
        <MenuItem
          label='My parent button'
          handleClick={() => setNestedOpened(!nestedOpened)}
          menuOpened={true}
          nestedOpened={nestedOpened}
          nested={
            <List component="div">
              <MenuItem
                label='My nested button'
                path='/my/nested/button'
                menuOpened={true}
              />
            </List>
          }
        />
      </ThemeProvider>
    ),
  },
  {
    path: '/my/simple/button',
    element: <div role="expected">My simple button page</div>,
  },
  {
    path: '/my/nested/button',
    element: <div role="expected">My nested button page</div>,
  },
]);

test('renders a menu item link', async () => {
  // Given
  let router: Router | null = null;

  const Wrapper = () => {
    const [nestedOpened, setNestedOpened] = useState(false);
    router = initializeRouter(nestedOpened, setNestedOpened);

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

  await userEvent.click(screen.getByText('My simple button'));

  // Then
  const element = expect(screen.getByRole('expected'));

  element.toHaveTextContent('My simple button page');
  element.toBeVisible();
});

test('renders a menu nested when button not clicked: nested items hidden', async () => {
  // Given
  let router: Router | null = null;

  const Wrapper = () => {
    const [nestedOpened, setNestedOpened] = useState(false);
    router = initializeRouter(nestedOpened, setNestedOpened);

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
  expect(screen.queryByText('My nested button')).toBeNull();
});

test('renders a menu nested item when parent button clicked: nested items shown', async () => {
  // Given
  let router: Router | null = null;

  const Wrapper = () => {
    const [nestedOpened, setNestedOpened] = useState(false);
    router = initializeRouter(nestedOpened, setNestedOpened);

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

  await userEvent.click(screen.getByText('My parent button'));

  // Then
  expect(screen.queryByText('My nested button')).toBeVisible();
});
