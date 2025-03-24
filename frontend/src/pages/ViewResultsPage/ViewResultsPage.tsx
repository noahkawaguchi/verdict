import { useParams } from 'react-router-dom';
import EnterPollId from '../../components/EnterPollId/EnterPollId';
import useGetRequest from '../../hooks/useGetRequest';
import { Result } from '../../types';
import { useEffect } from 'react';

/** Displays `EnterPollId` if one is not provided, then displays the results of the poll. */
const ViewResultsPage = () => {
  const { pollId } = useParams<{ pollId?: string }>();
  const { data, error, loading, sendRequest } = useGetRequest<Result>('result');

  useEffect(() => {
    if (pollId) sendRequest(pollId);
  }, [sendRequest, pollId]);

  return (
    <>
      <h3>View a Poll's Results</h3>
      <hr />
      {!pollId ? (
        <EnterPollId />
      ) : (
        <>
          {loading && <p>Loading...</p>}
          {error && <p>Failed to get result: {error.message}</p>}
          {data && (
            <p>
              In the poll "{data.prompt}," the choice "{data.winningChoice}" won with{' '}
              {data.winningVotes} out of {data.totalVotes} votes in round {data.winningRound}.
            </p>
          )}
        </>
      )}
    </>
  );
};

export default ViewResultsPage;
