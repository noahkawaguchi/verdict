import { useState } from 'react';
import { Question } from '../../types';

type CreatePollFormProps = {
  setQuestion: (question: Question) => void;
};

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
    setQuestion({
      prompt: prompt,
      choices: choices,
    });
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        placeholder='Prompt'
        required
      />
      {choices.map((choice, index) => (
        <div key={index}>
          <input
            value={choice}
            onChange={(e) => handleChoiceChange(index, e.target.value)}
            placeholder='Choice'
            required
          />
          {/* Specify the type as button so it doesn't try to submit the form */}
          <button
            type='button'
            onClick={() => setChoices((prev) => prev.filter((_, i) => i !== index))}
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
