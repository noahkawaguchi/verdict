import React, { useState } from 'react';
import { Question } from '../../types';
import styles from './CastBallot.module.css';

const CastBallotForm: React.FC<{ question: Question }> = ({ question }) => {
  const [ranks, setRanks] = useState(question.choices);

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
    console.log(rankOrder);
  }

  return (
    <form onSubmit={handleSubmit}>
      <p>Prompt: {question.prompt}</p>
      {ranks.map((choice, idx) => (
        <div key={idx} className={styles.rank}>
          <p>
            Rank {idx + 1}: {choice}
          </p>
          <button type='button' onClick={() => moveUp(idx)}>
            Move up
          </button>
          <button type='button' onClick={() => moveDown(idx)}>
            Move down
          </button>
        </div>
      ))}
      <button type='submit'>Submit</button>
    </form>
  );
};

export default CastBallotForm;
