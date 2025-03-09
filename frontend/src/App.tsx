import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import MainLayout from './layouts/MainLayout';
import HomePage from './pages/HomePage';
import CreatePollPage from './pages/CreatePoll/CreatePollPage';
import CastBallotPage from './pages/CastBallot/CastBallotPage';
import ViewResultsPage from './pages/ViewResults/ViewResultsPage';
import NotFoundPage from './pages/NotFoundPage';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path='/' element={<HomePage />} />
          <Route path='/create-poll' element={<CreatePollPage />} />
          <Route path='/cast-ballot/:pollID?' element={<CastBallotPage />} />
          <Route path='/view-results/:pollID?' element={<ViewResultsPage />} />
          <Route path='*' element={<NotFoundPage />} />
        </Route>
      </Routes>
    </Router>
  );
};

export default App;
