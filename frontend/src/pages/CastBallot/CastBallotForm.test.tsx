import { render, screen } from '@testing-library/react';
import { Question } from '../../types';
import CastBallotForm from './CastBallotForm';
import userEvent from '@testing-library/user-event';

describe('CastBallotForm', () => {
  const question: Question = {
    prompt: 'How would you describe VS Code?',
    choices: ['text editor', 'IDE', "don't know/just got here"],
  };

  it('should display the prompt and choices correctly', () => {
    render(<CastBallotForm question={question} setRankOrder={vi.fn()} />);
    expect(screen.getByText('How would you describe VS Code?')).toBeInTheDocument();
    expect(screen.getByText('text editor')).toBeInTheDocument();
    expect(screen.getByText('IDE')).toBeInTheDocument();
    expect(screen.getByText("don't know/just got here")).toBeInTheDocument();
  });

  it("should convert and set the user's choices as a properly formatted rank order", async () => {
    const mockSetRankOrder = vi.fn();
    render(<CastBallotForm question={question} setRankOrder={mockSetRankOrder} />);
    const user = userEvent.setup();
    await user.click(screen.getAllByText('Move up')[2]); // text editor, don't know, IDE
    await user.click(screen.getAllByText('Move up')[1]); // don't know, text editor, IDE
    await user.click(screen.getAllByText('Move down')[1]); // don't know, IDE, text editor
    await user.click(screen.getByText('Submit'));
    expect(mockSetRankOrder).toHaveBeenCalledExactlyOnceWith([2, 1, 0]);
  });
});
