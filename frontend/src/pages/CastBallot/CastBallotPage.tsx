import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import useGetRequest from '../../hooks/useGetRequest';
import EnterPollId from '../../components/EnterPollId/EnterPollId';
import CastBallotForm from './CastBallotForm';
import { Question } from '../../types';
import CastBallotSubmission from './CastBallotSubmission';

/**
 * Manages the display of `EnterPollId`, `CastBallotForm`, and `CastBallotSubmission` components.
 * Manages the timing of API requests to provide the voter with a ballot and to submit their
 * choices.
 */
const CastBallotPage = () => {
  const { pollId } = useParams<{ pollId?: string }>();
  const { data, error, loading, sendRequest } = useGetRequest<Question>('poll');
  const [rankOrder, setRankOrder] = useState<number[] | null>(null);

  useEffect(() => {
    if (pollId) sendRequest(pollId);
    else setRankOrder(null); // Avoid stale state when casting multiple ballots
  }, [sendRequest, pollId]);

  return (
    <>
      <h3>Cast a New Ballot</h3>
      <hr />
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
