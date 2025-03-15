import { useEffect } from 'react';
import useMutationRequest from '../../hooks/useMutationRequest';
import { Ballot } from '../../types';

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
