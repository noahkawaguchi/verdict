import { useEffect } from 'react';
import useApiRequest from '../../hooks/useApiRequest';
import { Question } from '../../types';

type CreatePollSubmissionProps = {
  question: Question;
};

const CreatePollSubmission: React.FC<CreatePollSubmissionProps> = ({ question }) => {
  const { data, error, loading, sendRequest } = useApiRequest<Question, { pollID: string }>(
    '/poll/create',
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
