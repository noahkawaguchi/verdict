import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import useGetRequest from '../../hooks/useGetRequest';
import EnterPollId from '../../components/EnterPollId/EnterPollId';
import CastBallotForm from './CastBallotForm';
import { Question } from '../../types';

const CastBallotPage = () => {
  const { pollId } = useParams<{ pollId?: string }>();
  const { data, error, loading, sendRequest } = useGetRequest<Question>('poll');

  useEffect(() => {
    if (pollId) sendRequest(pollId);
  }, [sendRequest, pollId]);

  return (
    <>
      <h3>Cast a New Ballot</h3>
      {!pollId ? (
        <EnterPollId />
      ) : (
        <>
          {loading && <p>Loading...</p>}
          {error && <p>Failed to create ballot: {error.message}</p>}
          {data && <CastBallotForm question={data} />}
        </>
      )}
    </>
  );
};

export default CastBallotPage;
