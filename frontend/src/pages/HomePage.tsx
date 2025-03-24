import { Link } from 'react-router-dom';

const HomePage = () => {
  return (
    <>
      <h2>Welcome to Verdict</h2>
      <p>
        <i>A place for ranked choice voting</i>
      </p>
      <h4>
        Read more about this project on{' '}
        <a
          href='https://github.com/noahkawaguchi/verdict'
          target='_blank'
          rel='noopener noreferrer'
        >
          GitHub
        </a>
      </h4>
      <p>
        Or <Link to='/create-poll'>create a new poll</Link>
      </p>
    </>
  );
};

export default HomePage;
