import React, { useState } from 'react';

function AddKeyValue({ onAdd }) {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [ttl, setTtl] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (key && value && ttl) {
      onAdd(key, value, ttl);
      setKey('');
      setValue('');
      setTtl('');
    }
  };

  return (
    <div>
      <h2>Add Key-Value Pair</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Key:</label>
          <input
            type="text"
            value={key}
            onChange={(e) => setKey(e.target.value)}
          />
        </div>
        <div>
          <label>Value:</label>
          <input
            type="text"
            value={value}
            onChange={(e) => setValue(e.target.value)}
          />
        </div>
        <div>
          <label>TTL (seconds):</label>
          <input
            type="number"
            value={ttl}
            onChange={(e) => setTtl(e.target.value)}
          />
        </div>
        <button type="submit">Add</button>
      </form>
    </div>
  );
}

export default AddKeyValue;
