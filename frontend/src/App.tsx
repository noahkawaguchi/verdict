import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import HostPollPage from './pages/HostPollPage';
import './App.css';
import HomePage from './pages/HomePage';
import CastBallotPage from './pages/CastBallotPage';
import NotFoundPage from './pages/NotFoundPage';
import MainLayout from './layouts/MainLayout';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path='/' element={<HomePage />} />
          <Route path='/host-poll' element={<HostPollPage />} />
          <Route path='/cast-ballot/:pollID?' element={<CastBallotPage />} />
          <Route path='*' element={<NotFoundPage />} />
        </Route>
      </Routes>
    </Router>
  );
};

export default App;
