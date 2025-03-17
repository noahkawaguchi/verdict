import { useEffect } from 'react';
import useMutationRequest from '../../hooks/useMutationRequest';
import { Question } from '../../types';

type CreatePollSubmissionProps = {
  question: Question;
};

/**
 * Makes a POST request to create the new poll. 
 * Displays loading, error, or links to vote in the poll and view its results.
 * 
 * @param question - The completed question to be used in the new poll.
 */
const CreatePollSubmission: React.FC<CreatePollSubmissionProps> = ({ question }) => {
  // Get the base URL to dynamically combine with the poll ID later
  const baseUrl = `${window.location.protocol}//${window.location.host}`;
  const { data, error, loading, sendRequest } = useMutationRequest<Question, { pollId: string }>(
    'poll',
    'POST',
  );

  useEffect(() => {
    sendRequest(question);
  }, [sendRequest, question]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Failed to create poll: {error.message}</p>;

  return (
    data && (
      <>
        <h4>
          <em>Your poll has been created!</em>
        </h4>
        <p>Voters can cast their ballots at this link:</p>
        <p>
          <a
            href={`${baseUrl}/cast-ballot/${data.pollId}`}
            target='_blank'
            rel='noopener noreferrer'
          >{`${baseUrl}/cast-ballot/${data.pollId}`}</a>
        </p>
        <button
          type='button'
          onClick={() => navigator.clipboard.writeText(`${baseUrl}/cast-ballot/${data.pollId}`)}
        >
          Copy voting link to clipboard
        </button>
        <p>The results will be available at this link:</p>
        <p>
          <a
            href={`${baseUrl}/view-results/${data.pollId}`}
            target='_blank'
            rel='noopener noreferrer'
          >{`${baseUrl}/view-results/${data.pollId}`}</a>
        </p>
        <button
          type='button'
          onClick={() => navigator.clipboard.writeText(`${baseUrl}/view-results/${data.pollId}`)}
        >
          Copy results link to clipboard
        </button>
      </>
    )
  );
};

export default CreatePollSubmission;
