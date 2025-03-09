import { useParams } from 'react-router-dom';
import EnterPollId from '../../components/EnterPollId/EnterPollId';

const CastBallotPage = () => {
  const { pollID } = useParams<{ pollID?: string }>();
  return (
    <>
      <h3>Cast a New Ballot</h3>
      {pollID ? <p>This is the cast ballot page for poll {pollID}!</p> : <EnterPollId />}
    </>
  );
};

export default CastBallotPage;
