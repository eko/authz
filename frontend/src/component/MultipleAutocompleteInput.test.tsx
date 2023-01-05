import { render, screen } from '@testing-library/react';
import { useState } from 'react';
import userEvent from '@testing-library/user-event';
import MultipleAutocompleteInput, { ItemType } from './MultipleAutocompleteInput';

test('MultipleAutocompleteInput: render default options', async () => {
    // Given
    const defaultValues: ItemType[] = [
        {id: '1', label: 'default-1'},
        {id: '2', label: 'default-2'},
    ];
  
    const errors: any = {};
  
    const fetcher = async (input: string): Promise<ItemType[]> => {
      return [];
    };
  
    const Wrapper = () => {
        const [values, setValues] = useState<ItemType[]>();
  
        return (
          <MultipleAutocompleteInput
              label='Rôles associés'
              placeholder='Rechercher un rôle ...'
              defaultValues={defaultValues}
              errorField={errors?.roles}
              fetcher={fetcher}
              setValue={(items: ItemType[]) => setValues(items)}
          />
        )
    };
  
    // Render
    render(
      <Wrapper />
    );
  
    // When
    await userEvent.click(screen.getByRole('multipleautocompleteinput-field'));
  
    // Then
    const listbox = screen.queryByRole('listbox');
    expect(listbox).toBeVisible();

    expect(listbox).toHaveTextContent('default-1');
    expect(listbox).toHaveTextContent('default-2');
  });

test('MultipleAutocompleteInput: render options using fetcher', async () => {
  // Given
  const defaultValues: ItemType[] = [];

  const errors: any = {};

  const fetcher = async (input: string): Promise<ItemType[]> => {
    return [
        { id: '1', label: 'role-1-' + input },
        { id: '2', label: 'role-2-' + input },
    ];
  };

  const Wrapper = () => {
      const [values, setValues] = useState<ItemType[]>();

      return (
        <MultipleAutocompleteInput
            label='Rôles associés'
            placeholder='Rechercher un rôle ...'
            defaultValues={defaultValues}
            errorField={errors?.roles}
            fetcher={fetcher}
            setValue={(items: ItemType[]) => setValues(items)}
        />
      )
  };

  // Render
  render(
    <Wrapper />
  );

  // When
  const inputTextField = screen.getByRole('multipleautocompleteinput-field');

  await userEvent.click(inputTextField);
  await userEvent.type(inputTextField, 'admin');

  // Then
  const listbox = screen.queryByRole('listbox');

  expect(listbox).toBeVisible();
  expect(listbox).toHaveTextContent('role-1-admin');
  expect(listbox).toHaveTextContent('role-2-admin');
});

test('MultipleAutocompleteInput: render list of selected items on select options', async () => {
    // Given
    const defaultValues: ItemType[] = [];

    const errors: any = {};

    const fetcher = async (input: string): Promise<ItemType[]> => {
        return [
            { id: '1', label: 'role-1-' + input },
            { id: '2', label: 'role-2-' + input },
        ];
    };

    const Wrapper = () => {
        const [values, setValues] = useState<ItemType[]>();

        return (
            <MultipleAutocompleteInput
                label='Rôles associés'
                placeholder='Rechercher un rôle ...'
                defaultValues={defaultValues}
                errorField={errors?.roles}
                fetcher={fetcher}
                setValue={(items: ItemType[]) => setValues(items)}
            />
        )
    };

    // Render
    render(
        <Wrapper />
    );

    // When
    const inputTextField = screen.getByRole('multipleautocompleteinput-field');

    await userEvent.click(inputTextField);
    await userEvent.type(inputTextField, 'admin');

    await userEvent.click(screen.getByText('role-2-admin'));

    // Then
    const items = screen.queryByRole('multipleautocompleteinput-item');

    expect(items).toHaveTextContent('role-2-admin');
});