import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import Header from './Header';

describe('Header', () => {
  beforeEach(() =>
    render(
      <MemoryRouter initialEntries={['/']}>
        <Header />
        <Routes>
          <Route path='/' element={<p>Home page here</p>} />
          <Route path='/create-poll' element={<p>Create poll page here</p>} />
          <Route path='/cast-ballot/:pollId?' element={<p>Cast ballot page here</p>} />
          <Route path='/view-results/:pollId?' element={<p>View results page here</p>} />
        </Routes>
      </MemoryRouter>,
    ),
  );

  it('should render', () => expect(screen.getByText('Verdict')).toBeInTheDocument());

  it('should redirect to the correct page', async () => {
    const user = userEvent.setup();

    await user.click(screen.getByText('Create Poll'));
    expect(screen.getByText('Create poll page here')).toBeInTheDocument();

    await user.click(screen.getByText('Cast Ballot'));
    expect(screen.getByText('Cast ballot page here')).toBeInTheDocument();

    await user.click(screen.getByText('View Results'));
    expect(screen.getByText('View results page here')).toBeInTheDocument();

    await user.click(screen.getByText('Verdict'));
    expect(screen.getByText('Home page here')).toBeInTheDocument();
  });
});
