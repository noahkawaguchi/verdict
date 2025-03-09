import { useParams } from 'react-router-dom';
import EnterPollId from '../../components/EnterPollId/EnterPollId';

const ViewResultsPage = () => {
  const { pollID } = useParams<{ pollID?: string }>();
  return (
    <>
      <h3>View a Poll's Results</h3>
      {pollID ? <p>This is the view results page for {pollID}!</p> : <EnterPollId />}
    </>
  );
};

export default ViewResultsPage;
