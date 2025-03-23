import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import MainLayout from './layouts/MainLayout';
import HomePage from './pages/HomePage/HomePage';
import CreatePollPage from './pages/CreatePoll/CreatePollPage';
import CastBallotPage from './pages/CastBallot/CastBallotPage';
import ViewResultsPage from './pages/ViewResults/ViewResultsPage';
import NotFoundPage from './pages/NotFoundPage/NotFoundPage';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path='/' element={<HomePage />} />
          <Route path='/create-poll' element={<CreatePollPage />} />
          <Route path='/cast-ballot/:pollId?' element={<CastBallotPage />} />
          <Route path='/view-results/:pollId?' element={<ViewResultsPage />} />
          <Route path='*' element={<NotFoundPage />} />
        </Route>
      </Routes>
    </Router>
  );
};

export default App;
