import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import useGetRequest from '../../hooks/useGetRequest';
import EnterPollId from '../../components/EnterPollId/EnterPollId';
import CastBallotForm from './CastBallotForm';
import { Question } from '../../types';
import CastBallotSubmission from './CastBallotSubmission';

const CastBallotPage = () => {
  const { pollId } = useParams<{ pollId?: string }>();
  const { data, error, loading, sendRequest } = useGetRequest<Question>('poll');
  const [rankOrder, setRankOrder] = useState<number[] | null>(null);

  useEffect(() => {
    if (pollId) sendRequest(pollId);
  }, [sendRequest, pollId]);

  return (
    <>
      <h3>Cast a New Ballot</h3>
      {!pollId ? (
        <EnterPollId />
      ) : !rankOrder ? (
        <>
          {loading && <p>Loading...</p>}
          {error && <p>Failed to create ballot: {error.message}</p>}
          {data && <CastBallotForm question={data} setRankOrder={setRankOrder} />}
        </>
      ) : (
        <CastBallotSubmission ballot={{ pollId, rankOrder }} />
      )}
    </>
  );
};

export default CastBallotPage;
