import { Question } from '../../types';

const CastBallotForm: React.FC<{ question: Question }> = ({ question }) => {
  return (
    <form>
      <p>
        This is the cast ballot form for poll {question.prompt} with choices {question.choices}!
      </p>
    </form>
  );
};

export default CastBallotForm;
