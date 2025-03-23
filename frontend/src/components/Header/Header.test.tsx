import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import HomePage from '../../pages/HomePage/HomePage';
import CreatePollPage from '../../pages/CreatePoll/CreatePollPage';
import CastBallotPage from '../../pages/CastBallot/CastBallotPage';
import ViewResultsPage from '../../pages/ViewResults/ViewResultsPage';
import NotFoundPage from '../../pages/NotFoundPage';
import Header from './Header';

describe('Header', () => {
  beforeEach(() =>
    render(
      <MemoryRouter initialEntries={['/']}>
        <Header />
        <Routes>
          <Route path='/' element={<HomePage />} />
          <Route path='/create-poll' element={<CreatePollPage />} />
          <Route path='/cast-ballot/:pollId?' element={<CastBallotPage />} />
          <Route path='/view-results/:pollId?' element={<ViewResultsPage />} />
          <Route path='*' element={<NotFoundPage />} />
        </Routes>
      </MemoryRouter>,
    ),
  );

  it('should render', () => expect(screen.getByText('Verdict')).toBeInTheDocument());

  it('should redirect to the correct page', async () => {
    const user = userEvent.setup();

    await user.click(screen.getByText('Create Poll'));
    expect(screen.getByText('Create a New Poll')).toBeInTheDocument();

    await user.click(screen.getByText('Cast Ballot'));
    expect(screen.getByText('Cast a New Ballot')).toBeInTheDocument();

    await user.click(screen.getByText('View Results'));
    expect(screen.getByText("View a Poll's Results")).toBeInTheDocument();

    await user.click(screen.getByText('Verdict'));
    expect(screen.getByText('Welcome to Verdict')).toBeInTheDocument();
  });
});
