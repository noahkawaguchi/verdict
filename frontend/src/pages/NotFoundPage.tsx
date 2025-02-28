import { Link } from 'react-router-dom';
import styles from './pages.module.css';

const NotFoundPage = () => {
  return (
    <>
      <p>404 - Page Not Found</p>
      <Link to='/' className={styles.routerLink}>
        Go to homepage
      </Link>
    </>
  );
};

export default NotFoundPage;
