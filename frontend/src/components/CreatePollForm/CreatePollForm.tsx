import { useState } from 'react';

const CreatePollForm = () => {
  const [prompt, setPrompt] = useState('');
  const [choices, setChoices] = useState(['', '']);

  const handleChoiceChange = (index: number, value: string) => {
    const updatedChoices = [...choices];
    updatedChoices[index] = value;
    setChoices(updatedChoices);
  };

  const addChoice = () => {
    setChoices([...choices, '']);
  };

  const removeChoice = (index: number) => {
    setChoices((prev) => [...prev.slice(0, index), ...prev.slice(index + 1)]);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('submitted!'); // This will use the hook
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
          <button onClick={() => removeChoice(index)} disabled={choices.length <= 2}>
            Remove
          </button>
        </div>
      ))}
      <button onClick={addChoice}>Add Choice</button>
      <button type='submit'>Submit</button>
    </form>
  );
};

export default CreatePollForm;
