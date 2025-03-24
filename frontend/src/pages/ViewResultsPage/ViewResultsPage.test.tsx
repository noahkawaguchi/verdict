import { MemoryRouter, Route, Routes } from 'react-router-dom';
import ViewResultsPage from './ViewResultsPage';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { backendUrl } from '../../config';
import { Result } from '../../types';

describe('ViewResultsPage', () => {
  beforeEach(() => {
    globalThis.fetch = vi.fn();
    render(
      <MemoryRouter initialEntries={['/view-results']}>
        <Routes>
          <Route path='/view-results/:pollId?' element={<ViewResultsPage />} />
        </Routes>
      </MemoryRouter>,
    );
  });

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should show EnterPollId before a poll ID is set', () => {
    expect(screen.getByText('Paste the poll ID here:')).toBeInTheDocument();
  });

  it('should handle error responses', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: 'some error was encountered' }),
    });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText('Paste the poll ID here:'));
    await user.paste('poll-135');
    await user.click(screen.getByText('Submit'));
    expect(screen.getByText(/some error was encountered/)).toBeInTheDocument();
  });

  it('should invoke useGetRequest with the correct arguments', async () => {
    const result: Result = {
      prompt: 'What is the best apple color?',
      totalVotes: 16,
      winningVotes: 9,
      winningChoice: 'red',
      winningRound: 1,
    };
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => result });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText('Paste the poll ID here:'));
    await user.paste('poll-135');
    await user.click(screen.getByText('Submit'));
    expect(mockFetch).toHaveBeenCalledExactlyOnceWith(`${backendUrl}/result/poll-135`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    });
  });

  it('should display the results correctly', async () => {
    const result: Result = {
      prompt: 'What is the best apple color?',
      totalVotes: 16,
      winningVotes: 9,
      winningChoice: 'red',
      winningRound: 1,
    };
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => result });
    const user = userEvent.setup();
    await user.click(screen.getByLabelText('Paste the poll ID here:'));
    await user.paste('poll-135');
    await user.click(screen.getByText('Submit'));
    expect(
      screen.getByText(
        `In the poll "${result.prompt}," the choice "${result.winningChoice}" won with ` +
          `${result.winningVotes} out of ${result.totalVotes} votes in round ` +
          `${result.winningRound}.`,
      ),
    ).toBeInTheDocument();
  });
});
