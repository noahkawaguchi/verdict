import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import CreatePollPage from './pages/CreatePollPage';
import './App.css';
import HomePage from './pages/HomePage';
import CastBallotPage from './pages/CastBallotPage';
import NotFoundPage from './pages/NotFoundPage';
import MainLayout from './layouts/MainLayout';
import ViewResultsPage from './pages/ViewResultsPage';

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
