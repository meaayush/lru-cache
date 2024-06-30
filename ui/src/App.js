// src/App.js
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import CacheDisplay from './CacheDisplay';
import AddKeyValue from './AddKeyValue';
import GetKeyValue from './GetKeyValue';

function App() {
  const [cache, setCache] = useState({});
  const [singleValue, setSingleValue] = useState("");

  const fetchCache = async () => {
    try {
      const response = await axios.get('http://localhost:3001/mycache/all');
      setCache(response.data);
    } catch (error) {
      console.error('Error fetching cache:', error);
    }
  };

  const fetchSingleKey = async (key) => {
    try {
      const response = await axios.get(`http://localhost:3001/mycache?key=${key}`);
      console.log(response.data)
      setSingleValue(response.data);
    } catch (error) {
      console.error('Error fetching key:', error);
      setSingleValue('Key not found');
    }
  };

  useEffect(() => {
    fetchCache();
  }, []);

  const addKeyValue = async (key, value, ttl) => {
    try {
      await axios.post('http://localhost:3001/mycache', { key, value, ttl});
      fetchCache();
    } catch (error) {
      console.error('Error adding key-value pair:', error);
    }
  };

  return (
    <div className="App">
      <h1>LRU Cache</h1>
      <AddKeyValue onAdd={addKeyValue} />
      <GetKeyValue onGet={fetchSingleKey} singleValue={singleValue} />
      <CacheDisplay cache={cache} />
    </div>
  );
}

export default App;
