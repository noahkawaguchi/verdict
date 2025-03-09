import { useParams } from 'react-router-dom';

const ViewResultsPage = () => {
  const { pollID } = useParams<{ pollID?: string }>();
  return (
    <>
      {pollID ? (
        <p>This is the view results page for {pollID}!</p>
      ) : (
        <p>This is the view results page without a poll ID provided!</p>
      )}
    </>
  );
};

export default ViewResultsPage;
