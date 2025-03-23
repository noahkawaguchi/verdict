import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import EnterPollId from './EnterPollId';

describe('EnterPollId', () => {
  it('redirects to the correct URL on submit', async () => {
    render(
      <MemoryRouter initialEntries={['/dummy-path']}>
        <Routes>
          <Route path='/dummy-path' element={<EnterPollId />} />
          <Route path='/dummy-path/poll-id-123' element={<p>Successfully redirected!</p>} />
        </Routes>
      </MemoryRouter>,
    );
    const user = userEvent.setup();
    await user.type(screen.getByLabelText('Paste the poll ID here:'), 'poll-id-123');
    await user.click(screen.getByText('Submit'));
    expect(screen.getByText('Successfully redirected!')).toBeInTheDocument();
  });
});
