import { Link } from 'react-router-dom';
import styles from './Header.module.css';

const Header = () => {
  return (
    <header>
      <div>
        <h1>Verdict</h1>
        <nav className={styles.navbar}>
          <Link to='/' className={styles.routerLink}>
            Home
          </Link>
          <Link to='/host-poll' className={styles.routerLink}>
            Host Poll
          </Link>
          <Link to='/cast-ballot' className={styles.routerLink}>
            Cast Ballot
          </Link>
        </nav>
      </div>
      <hr />
    </header>
  );
};

export default Header;
