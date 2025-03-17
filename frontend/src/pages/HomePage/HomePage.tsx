import { Link } from 'react-router-dom';
import styles from './HomePage.module.css';

const HomePage = () => {
  return (
    <>
      <h2>Welcome to Verdict</h2>
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
      <div className={`${styles.explanationParent} ${styles.explanationBlock}`}>
        <h3>What is ranked choice voting?</h3>
        <p>
          Verdict implements a voting algorithm known as <strong>ranked choice voting</strong> or{' '}
          <strong>instant runoff voting</strong>. There are several variations, but this is the one
          used here:
        </p>
        <div className={styles.explanationBlock}>
          <h4>Voting process</h4>
          <p>
            Instead of selecting a single choice, voters <strong>rank each choice</strong> as their
            first choice, second choice, third choice, etc., all the way to last choice.
          </p>
        </div>
        <div className={styles.explanationBlock}>
          <h4>Determining a winner</h4>
          <p>
            Instead of immediately selecting a choice that has only a <strong>plurality</strong> of
            votes, the algorithm first checks if any choice has a <strong>strict majority</strong>{' '}
            of votes. If no choice does,{' '}
            <strong>the choice with the fewest votes is eliminated</strong>, and its votes are{' '}
            <strong>redistributed</strong> to the voters' next highest choices. This process of
            elimination <strong>continues until a single choice has a strict majority</strong> of
            votes.
          </p>
        </div>
        <div className={styles.explanationBlock}>
          <h4>What about ties for last?</h4>
          <p>
            In each round, the choice with the fewest votes is eliminated, but what if{' '}
            <strong>multiple choices are tied</strong> for last place? In this case, a{' '}
            <strong>sub-poll is simulated</strong> between only the tied choices. This is possible
            because voters provide a rank for every choice, allowing the algorithm to{' '}
            <strong>determine their preferences amongst any subset</strong> of choices. If there is
            another tie for last place, the tie-breaking algorithm{' '}
            <strong>continues recursively</strong>.
          </p>
          <p>
            While unlikely unless the numbers of choices and voters are very small, it is possible
            that multiple choices are tied for last place and received the{' '}
            <strong>exact same rankings</strong>. In this case, a choice is eliminated by
            pseudorandom number generation.
          </p>
        </div>
      </div>
    </>
  );
};

export default HomePage;
