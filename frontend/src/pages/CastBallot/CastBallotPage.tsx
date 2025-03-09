import { useParams } from 'react-router-dom';

const CastBallotPage = () => {
  const { pollID } = useParams<{ pollID?: string }>();
  return (
    <>
      {pollID ? (
        <p>This is the cast ballot page for poll {pollID}!</p>
      ) : (
        <p>This is the cast ballot page without a poll ID provided!</p>
      )}
    </>
  );
};

export default CastBallotPage;
