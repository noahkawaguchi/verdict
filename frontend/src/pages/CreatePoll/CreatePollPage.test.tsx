import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import CreatePollPage from './CreatePollPage';

describe('CreatePollPage', () => {
  beforeEach(() => (globalThis.fetch = vi.fn()));

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should show CreatePollForm before a question is set, CreatePollSubmission after', async () => {
    const mockedFetch = globalThis.fetch as ReturnType<typeof vi.fn>;
    mockedFetch.mockResolvedValueOnce({ ok: true, json: async () => ({ pollId: 'poll-12345' }) });

    render(<CreatePollPage />);
    const user = userEvent.setup();
    const addChoiceBtn = screen.getByText('Add Choice');
    expect(addChoiceBtn).toBeInTheDocument();

    await user.type(screen.getByLabelText('Prompt:'), 'What is the best color?');
    await user.type(screen.getByLabelText('Choice 1:'), 'red');
    await user.type(screen.getByLabelText('Choice 2:'), 'blue');
    await user.click(screen.getByText('Submit'));

    expect(addChoiceBtn).not.toBeInTheDocument();
    expect(screen.getByText('Your poll has been created!')).toBeInTheDocument();
  });
});
