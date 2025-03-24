import { render, screen } from '@testing-library/react';
import { act } from 'react';
import CreatePollSubmission from './CreatePollSubmission';

describe('CreatePollSubmission', () => {
  beforeEach(() => (globalThis.fetch = vi.fn()));

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should correctly invoke useMutationRequest when it renders', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({ message: 'success!' }) });
    const question = {
      prompt: 'What language should I use for my frontend?',
      choices: ['TypeScript', 'Rust'],
    };
    await act(async () => render(<CreatePollSubmission question={question} />));
    expect(mockFetch).toHaveBeenCalledWith('undefined/poll', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: `{"prompt":"What language should I use for my frontend?","choices":["TypeScript","Rust"]}`,
    });
  });

  it('should display the received data correctly', async () => {
    const mockFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => ({ pollId: 'poll-678' }) });
    const question = {
      prompt: 'What language should I use for my frontend?',
      choices: ['TypeScript', 'Rust'],
    };
    await act(async () => render(<CreatePollSubmission question={question} />));
    expect(screen.getAllByText(/poll-678/)[0]).toBeInTheDocument();
  });
});
