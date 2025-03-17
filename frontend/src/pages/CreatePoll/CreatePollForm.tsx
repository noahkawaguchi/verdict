import { useState } from 'react';
import { Question } from '../../types';

type CreatePollFormProps = {
  setQuestion: (question: Question) => void;
};

/**
 * A form for entering the prompt and choices for a new poll.
 * 
 * @param setQuestion - Function to update the question in the parent component.
 */
const CreatePollForm: React.FC<CreatePollFormProps> = ({ setQuestion }) => {
  const [prompt, setPrompt] = useState('');
  const [choices, setChoices] = useState(['', '']);

  const handleChoiceChange = (index: number, value: string) => {
    const updatedChoices = [...choices];
    updatedChoices[index] = value;
    setChoices(updatedChoices);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Enforce unique choices
    if (new Set(choices).size !== choices.length) {
      alert('choices must be unique');
      return;
    }
    setQuestion({
      prompt: prompt,
      choices: choices,
    });
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Prompt:{' '}
        <input
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          placeholder='What is the best fruit?'
          required
          autoFocus
        />
      </label>
      {choices.map((choice, idx) => (
        <div key={idx}>
          <label>
            Choice {idx + 1}:{' '}
            <input
              value={choice}
              onChange={(e) => handleChoiceChange(idx, e.target.value)}
              placeholder='Banana'
              required
            />
          </label>
          {/* Specify the type as button so it doesn't try to submit the form */}
          <button
            type='button'
            onClick={() => setChoices((prev) => prev.filter((_, i) => i !== idx))}
            disabled={choices.length <= 2}
          >
            Remove
          </button>
        </div>
      ))}
      <button type='button' onClick={() => setChoices([...choices, ''])}>
        Add Choice
      </button>
      <button type='submit'>Submit</button>
    </form>
  );
};

export default CreatePollForm;
