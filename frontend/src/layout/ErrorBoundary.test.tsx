import { render, screen } from '@testing-library/react';
import { useEffect } from 'react';
import ErrorBoundary from 'layout/ErrorBoundary';

// Hide all console errors during tests output by mocking console.error.
beforeEach(() => {
    jest.spyOn(console, 'error').mockImplementation(() => {});
});

test('ErrorBoundary: render component when no error', async () => {
  // Given
  const Children = () => {
    return (
      <div role='expected'>My normal case component</div>
    )
  }

  // Render
  render(
    <ErrorBoundary>
        <Children />
    </ErrorBoundary>
  );

  // When - Then
  const element = expect(screen.getByRole('expected'));
  element.toHaveTextContent('My normal case component');
  element.toBeVisible();

  expect(screen.queryByText('Une erreur est survenue.')).toBeNull();
});

test('ErrorBoundary: render component when error is thrown', async () => {
  // Given
  const Children = () => {
      useEffect(() => {
          throw new Error('This is an unexpected error...');
      }, []);

      return (
          <div role='expected'>My normal case component</div>
      )
  }

  // Render
  render(
      <ErrorBoundary>
          <Children />
      </ErrorBoundary>
    );
  
    // When - Then
    expect(screen.queryByText('An error occurred.')).toBeVisible();
    expect(screen.queryByText('This is an unexpected error...')).toBeVisible();
});