import { useEffect } from 'react';
import useMutationRequest from '../../hooks/useMutationRequest';
import { Ballot } from '../../types';

/**
 * Makes a POST request to cast the voter's ballot and displays loading, error, or a confirmation 
 * message.
 * 
 * @param ballot - The voter's completed ballot.
 */
const CastBallotSubmission: React.FC<{ ballot: Ballot }> = ({ ballot }) => {
  const { data, error, loading, sendRequest } = useMutationRequest<Ballot, { message: string }>(
    'ballot',
    'POST',
  );

  useEffect(() => {
    sendRequest(ballot);
  }, [sendRequest, ballot]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Failed to cast ballot: {error.message}</p>;

  return <p>{data?.message}</p>;
};

export default CastBallotSubmission;
