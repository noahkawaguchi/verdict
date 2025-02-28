import { Link } from 'react-router-dom';
import styles from './pages.module.css';

const HomePage = () => {
  return (
    <>
      <Link to='/host-poll' className={styles.routerLink}>Host Poll</Link>
      <br />
      <Link to='/cast-ballot' className={styles.routerLink}>Cast Ballot</Link>
    </>
  );
};

export default HomePage;
