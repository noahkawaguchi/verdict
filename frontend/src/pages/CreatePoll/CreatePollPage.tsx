import { useState } from 'react';
import CreatePollForm from './CreatePollForm';
import { Question } from '../../types';
import CreatePollSubmission from './CreatePollSubmission';

const CreatePollPage = () => {
  const [question, setQuestion] = useState<Question | null>(null);

  return (
    <>
      <p>This is the create poll page!</p>
      {question ? (
        <CreatePollSubmission question={question} />
      ) : (
        <CreatePollForm setQuestion={setQuestion} />
      )}
    </>
  );
};

export default CreatePollPage;
