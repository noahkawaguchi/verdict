import { render, screen } from '@testing-library/react';
import CastBallotSubmission from './CastBallotSubmission';
import { Ballot } from '../../types';
import { act } from 'react';

describe('CastBallotSubmission', () => {
  beforeEach(() => (globalThis.fetch = vi.fn()));

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should fetch and display submission results correctly on render', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ message: 'Hello from mock fetch!' }),
    });
    const ballot: Ballot = { pollId: 'poll-987', rankOrder: [3, 0, 1, 2] };
    await act(async () => render(<CastBallotSubmission ballot={ballot} />));
    expect(screen.getByText('Hello from mock fetch!')).toBeInTheDocument();
  });

  it('should handle error responses', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: 'something went wrong!' }),
    });
    const ballot: Ballot = { pollId: 'poll-987', rankOrder: [3, 0, 1, 2] };
    await act(async () => render(<CastBallotSubmission ballot={ballot} />));
    expect(screen.getByText(/something went wrong!/)).toBeInTheDocument();
  });
});
