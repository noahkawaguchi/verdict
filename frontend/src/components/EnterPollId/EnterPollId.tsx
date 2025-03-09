import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';

const EnterPollId = () => {
  const [idInput, setIdInput] = useState('');
  const navigate = useNavigate();
  const location = useLocation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    navigate(`${location.pathname}/${idInput}`);
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Paste the poll ID here:
        <br />
        <input
          value={idInput}
          onChange={(e) => setIdInput(e.target.value)}
          placeholder='Ex: 8471a9ab...'
          required
          autoFocus
        />
      </label>
      <button type='submit'>Submit</button>
    </form>
  );
};

export default EnterPollId;
