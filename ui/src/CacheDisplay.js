import React from 'react';

function CacheDisplay({ cache }) {
  return (
    <div>
      <h2>Cache State</h2>
      <ul>
        {Object.entries(cache).map(([key, details]) => (
          <li key={key}>
            <strong>Key: </strong>{key} <br></br>
              <strong>Value:</strong> {details.value} <br></br>
              <strong>Expiration:</strong> {new Date(details.expiration).toLocaleString()}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default CacheDisplay;
