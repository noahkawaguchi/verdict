import React, { useEffect, useState } from 'react';
import { Question } from '../../types';
import styles from './CastBallot.module.css';

type CastBallotFormProps = {
  question: Question;
  setRankOrder: (rankOrder: number[] | null) => void;
};

/**
 * A form for casting a ballot. Allows the voter to rank their choices.
 *
 * @param question - The question to pose to the voter.
 * @param setRankOrder - A function to update the rank order in the parent component.
 */
const CastBallotForm: React.FC<CastBallotFormProps> = ({ question, setRankOrder }) => {
  const [ranks, setRanks] = useState(question.choices);

  // Clear old ranks and rank order when the question changes in case this component was just 
  // hidden without fully remounting
  useEffect(() => {
    setRanks(question.choices);
    setRankOrder(null);
  }, [question, setRankOrder]);

  const moveUp = (rankIdx: number) => {
    const updatedRanks = [...ranks];
    if (rankIdx > 0) {
      [updatedRanks[rankIdx - 1], updatedRanks[rankIdx]] = [
        updatedRanks[rankIdx],
        updatedRanks[rankIdx - 1],
      ];
    }
    setRanks(updatedRanks);
  };

  const moveDown = (rankIdx: number) => {
    const updatedRanks = [...ranks];
    if (rankIdx < updatedRanks.length - 1) {
      [updatedRanks[rankIdx + 1], updatedRanks[rankIdx]] = [
        updatedRanks[rankIdx],
        updatedRanks[rankIdx + 1],
      ];
    }
    setRanks(updatedRanks);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Create an array of indices in the format expected in the backend
    const rankOrder = ranks.map((rank) => question.choices.indexOf(rank));
    setRankOrder(rankOrder);
  };

  return (
    <form onSubmit={handleSubmit}>
      <p>
        <i>Prompt:</i> {question.prompt}
      </p>
      <hr />
      {ranks.map((choice, idx) => (
        <div key={idx}>
          <div className={styles.rank}>
            <p>
              <i>Rank {idx + 1}:</i> {choice}
            </p>
            <button
              type='button'
              className={styles.moveUpBtn}
              onClick={() => moveUp(idx)}
              disabled={idx === 0}
            >
              Move up
            </button>
            <button
              type='button'
              onClick={() => moveDown(idx)}
              disabled={idx === ranks.length - 1}
            >
              Move down
            </button>
          </div>
          <hr style={{ border: '1px solid black', width: '95%' }} />
        </div>
      ))}
      <button type='submit'>Submit</button>
    </form>
  );
};

export default CastBallotForm;
