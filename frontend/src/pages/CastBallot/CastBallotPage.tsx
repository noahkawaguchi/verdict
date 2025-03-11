import { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import useGetRequest from '../../hooks/useGetRequest';
import EnterPollId from '../../components/EnterPollId/EnterPollId';
import CastBallotForm from './CastBallotForm';
import { Question } from '../../types';

const CastBallotPage = () => {
  const { pollID } = useParams<{ pollID?: string }>();
  const { data, error, loading, sendRequest } = useGetRequest<Question>('/ballot/new');

  useEffect(() => {
    if (pollID) sendRequest(pollID);
  }, [sendRequest, pollID]);

  return (
    <>
      <h3>Cast a New Ballot</h3>
      {!pollID ? (
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
