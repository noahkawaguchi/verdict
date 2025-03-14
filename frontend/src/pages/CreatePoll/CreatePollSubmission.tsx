import { useEffect } from 'react';
import useMutationRequest from '../../hooks/useMutationRequest';
import { Question } from '../../types';

type CreatePollSubmissionProps = {
  question: Question;
};

const CreatePollSubmission: React.FC<CreatePollSubmissionProps> = ({ question }) => {
  const { data, error, loading, sendRequest } = useMutationRequest<Question, { pollID: string }>(
    'poll',
    'POST',
  );

  useEffect(() => {
    sendRequest(question);
  }, [sendRequest, question]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Failed to create poll: {error.message}</p>;

  return <p>Poll ID: {data?.pollID}</p>;
};

export default CreatePollSubmission;
