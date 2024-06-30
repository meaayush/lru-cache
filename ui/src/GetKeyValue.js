import React, { useState } from 'react';

function GetKeyValue({ onGet, singleValue }) {
  const [key, setKey] = useState('');

  const handleGet = (e) => {
    e.preventDefault();
    if (key) {
      onGet(key);
    }
  };

  return (
    <div>
      <h2>Get Value by Key</h2>
      <form onSubmit={handleGet}>
        <div>
          <label>Key:</label>
          <input
            type="text"
            value={key}
            onChange={(e) => setKey(e.target.value)}
          />
        </div>
        <button type="submit">Get Value</button>
      </form>
      {singleValue !== null && (
        <div>
          <h3>Value for "{key}": {singleValue}</h3>
        </div>
      )}
    </div>
  );
}

export default GetKeyValue;
