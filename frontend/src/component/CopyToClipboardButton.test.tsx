import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CopyToClipboardButton from './CopyToClipboardButton';

test('CopyToClipboardButton: copy text to clipboard', async () => {
    // Given
    const user = userEvent.setup();

    // Render
    render(
        <CopyToClipboardButton text='This is my test message' />
    );

    // When clicking on the icon
    // it should be copied to clipboard
    await user.click(screen.getByRole('copy-icon'));

    // Then check confirmation message to be present
    const confirmation = screen.getByRole('copy-confirmation');

    expect(confirmation).toBeVisible();
    expect(confirmation).toHaveTextContent('Copié');

    // And ensure right content has been copied to clipboard
    const value = await navigator.clipboard.readText();
    expect(value).toEqual('This is my test message');

    // Finally, confirmation should disappear after 1 second
    await act(async () => {
        await new Promise(r => setTimeout(r, 1500));
    });

    expect(screen.queryByRole('copy-confirmation')).toBeNull();
});