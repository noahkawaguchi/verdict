import { useState } from 'react';
import CreatePollForm from '../components/CreatePollForm/CreatePollForm';
import { Question } from '../types';
import CreatePollSubmission from '../components/CreatePollSubmission/CreatePollSubmission';

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
