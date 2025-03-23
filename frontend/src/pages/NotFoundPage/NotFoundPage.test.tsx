import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import NotFoundPage from './NotFoundPage';

describe('NotFoundPage', () => {
  it('should link back to the home page', async () => {
    render(
      <MemoryRouter initialEntries={['/undefined-page']}>
        <Routes>
          <Route path='/' element={<p>Home page here!</p>} />
          <Route path='*' element={<NotFoundPage />} />
        </Routes>
      </MemoryRouter>,
    );
    const user = userEvent.setup();
    await user.click(screen.getByText('Back to home'));
    expect(screen.getByText('Home page here!')).toBeInTheDocument();
  });
});
