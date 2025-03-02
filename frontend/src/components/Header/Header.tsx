import { Link } from 'react-router-dom';
import styles from './Header.module.css';

const Header = () => {
  return (
    <header className={styles.headerBox}>
      <Link to='/' className={styles.routerLink}>
        <h2 className={styles.headerTitle}>Verdict</h2>
      </Link>
      <nav className={styles.navbar}>
        <Link to='/create-poll' className={styles.routerLink}>
          Create Poll
        </Link>
        <Link to='/cast-ballot' className={styles.routerLink}>
          Cast Ballot
        </Link>
        <Link to='/view-results' className={styles.routerLink}>
          View Results
        </Link>
      </nav>
    </header>
  );
};

export default Header;
