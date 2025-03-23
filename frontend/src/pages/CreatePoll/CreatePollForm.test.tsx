import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import CreatePollForm from './CreatePollForm';

describe('CreatePollForm', () => {
  beforeEach(() => (window.alert = vi.fn()));

  afterEach(() => vi.resetAllMocks());

  afterAll(() => vi.restoreAllMocks());

  it('should enforce unique choices', async () => {
    const mockAlert = window.alert as ReturnType<typeof vi.fn>;
    const mockSetQuestion = vi.fn();
    render(<CreatePollForm setQuestion={mockSetQuestion} />);
    const user = userEvent.setup();

    await user.type(
      screen.getByLabelText('Prompt:'),
      'What language should I use for my frontend?',
    );
    await user.type(screen.getByLabelText('Choice 1:'), 'TypeScript');
    await user.type(screen.getByLabelText('Choice 2:'), 'TypeScript');
    await user.click(screen.getByText('Submit'));

    expect(mockAlert).toHaveBeenCalledExactlyOnceWith('choices must be unique');
  });

  it('should correctly set the question on successful form submission', async () => {
    const mockSetQuestion = vi.fn();
    render(<CreatePollForm setQuestion={mockSetQuestion} />);
    const user = userEvent.setup();

    await user.type(
      screen.getByLabelText('Prompt:'),
      'What language should I use for my frontend?',
    );
    await user.type(screen.getByLabelText('Choice 1:'), 'TypeScript');
    await user.type(screen.getByLabelText('Choice 2:'), 'Rust');
    await user.click(screen.getByText('Submit'));

    expect(mockSetQuestion).toHaveBeenCalledExactlyOnceWith({
      prompt: 'What language should I use for my frontend?',
      choices: ['TypeScript', 'Rust'],
    });
  });
});
