import { Link } from "react-router-dom";

const HomePage = () => {
  return (
    <>
      <h3>Welcome to Verdict</h3>
      <p>
        <i>A place for ranked choice voting</i>
      </p>
      <h4>
        Check this project out on{' '}
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
