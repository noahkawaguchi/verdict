import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import MainLayout from './layouts/MainLayout';
import HomePage from './pages/HomePage';
import CreatePollPage from './pages/CreatePoll/CreatePollPage';
import CastBallotPage from './pages/CastBallot/CastBallotPage';
import ViewResultsPage from './pages/ViewResultsPage/ViewResultsPage';
import NotFoundPage from './pages/NotFoundPage/NotFoundPage';

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path='/' element={<HomePage />} />
          <Route path='/create-poll' element={<CreatePollPage />} />
          <Route path='/cast-ballot/:pollId?' element={<CastBallotPage />} />
          <Route path='/view-results/:pollId?' element={<ViewResultsPage />} />
          <Route path='*' element={<NotFoundPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
