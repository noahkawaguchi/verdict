import { useParams } from 'react-router-dom';
import EnterPollId from '../../components/EnterPollId/EnterPollId';

const ViewResultsPage = () => {
  const { pollId } = useParams<{ pollId?: string }>();
  return (
    <>
      <h3>View a Poll's Results</h3>
      {pollId ? <p>This is the view results page for {pollId}!</p> : <EnterPollId />}
    </>
  );
};

export default ViewResultsPage;
