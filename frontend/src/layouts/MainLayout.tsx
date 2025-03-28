import { Outlet } from 'react-router-dom';
import Header from '../components/Header/Header';
import styles from './MainLayout.module.css';

const MainLayout = () => {
  return (
    <div className={styles.fullScreenWrapper}>
      <div className={styles.mainLayoutDiv}>
        <Header />
        <main>
          <Outlet />
        </main>
        <hr />
      </div>
    </div>
  );
};

export default MainLayout;
