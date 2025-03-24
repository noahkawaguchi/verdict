import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import CastBallotPage from './CastBallotPage';
import { MemoryRouter, Routes, Route } from 'react-router-dom';

describe('CastBallotPage', () => {
  beforeEach(() => {
    globalThis.fetch = vi.fn();
    render(
      <MemoryRouter initialEntries={['/cast-ballot']}>
        <Routes>
          <Route path='/cast-ballot/:pollId?' element={<CastBallotPage />} />
        </Routes>
      </MemoryRouter>,
    );
  });

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should show EnterPollId before a poll ID is set', () => {
    expect(screen.getByText('Paste the poll ID here:')).toBeInTheDocument();
  });

  it('should show CastBallotForm after a poll ID is set but before a rank order is set', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        prompt: 'What is the best apple color?',
        choices: ['red', 'green', 'yellow'],
      }),
    });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText('Paste the poll ID here:'));
    await user.paste('poll-789');
    await user.click(screen.getByText('Submit'));
    expect(screen.getAllByText('Move up')[0]).toBeInTheDocument();
  });

  it('should show CastBallotSubmission after a rank order is set', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          prompt: 'What is the best apple color?',
          choices: ['red', 'green', 'yellow'],
        }),
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({ message: 'Successfully cast ballot' }),
      });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText('Paste the poll ID here:'));
    await user.paste('poll-789');
    await user.click(screen.getByText('Submit')); // EnterPollId submit button
    await user.click(screen.getByText('Submit')); // CastBallotForm submit button
    expect(screen.getByText('Successfully cast ballot')).toBeInTheDocument();
  });
});
