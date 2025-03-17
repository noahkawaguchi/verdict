import { useState } from 'react';
import CreatePollForm from './CreatePollForm';
import { Question } from '../../types';
import CreatePollSubmission from './CreatePollSubmission';

const CreatePollPage = () => {
  const [question, setQuestion] = useState<Question | null>(null);

  return (
    <>
      <h3>Create a New Poll</h3>
      {!question ? (
        <CreatePollForm setQuestion={setQuestion} />
      ) : (
        <CreatePollSubmission question={question} />
      )}
    </>
  );
};

export default CreatePollPage;
