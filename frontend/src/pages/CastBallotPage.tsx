import { useParams } from 'react-router-dom';

const CastBallotPage = () => {
  const { pollID } = useParams<{ pollID?: string }>();
  if (pollID) {
    return <p>This is the cast ballot page for poll {pollID}!</p>;
  } else {
    return <p>This is the cast ballot page without a poll ID provided!</p>;
  }
};

export default CastBallotPage;
