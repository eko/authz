import { render, screen } from '@testing-library/react';
import { useState } from 'react';
import userEvent from '@testing-library/user-event';
import { ItemType } from './MultipleAutocompleteInput';
import SingleAutocompleteInput from './SingleAutocompleteInput';

test('SingleAutocompleteInput: render default option', async () => {
    // Given
    const defaultValue: ItemType = {
        id: '1',
        label: 'default-1',
    };
  
    const errors: any = {};
  
    const fetcher = async (input: string): Promise<ItemType[]> => {
      return [];
    };
  
    const Wrapper = () => {
        const [value, setValue] = useState<ItemType>();
  
        return (
            <SingleAutocompleteInput
                label='Item'
                placeholder='Search for an item...'
                defaultValue={defaultValue}
                errorField={errors?.item}
                fetcher={fetcher}
                setValue={(value: ItemType) => setValue(value)}
            />
        )
    };
  
    // Render
    render(
      <Wrapper />
    );
  
    // When
    await userEvent.click(screen.getByRole('singleautocompleteinput-field'));
  
    // Then
    const listbox = screen.queryByRole('listbox');

    expect(listbox).toBeVisible();
    expect(listbox).toHaveTextContent('default-1');
});

test('SingleAutocompleteInput: render options using fetcher', async () => {
    // Given
    const errors: any = {};
  
    const fetcher = async (input: string): Promise<ItemType[]> => {
      return [
          {id: '1', label: 'item-1-' + input},
          {id: '2', label: 'item-2-' + input},
      ];
    };
  
    const Wrapper = () => {
        const [value, setValue] = useState<ItemType>();
  
        return (
            <SingleAutocompleteInput
                label='Item'
                placeholder='Search for an item...'
                errorField={errors?.item}
                fetcher={fetcher}
                setValue={(value: ItemType) => setValue(value)}
            />
        )
    };
  
    // Render
    render(
      <Wrapper />
    );
  
    // When
    const inputTextField = screen.getByRole('singleautocompleteinput-field');

    await userEvent.click(inputTextField);
    await userEvent.type(inputTextField, 'something');
    
    // Then
    const listbox = screen.queryByRole('listbox');

    expect(listbox).toBeVisible();
    expect(listbox).toHaveTextContent('item-1-something');
    expect(listbox).toHaveTextContent('item-2-something');
});