import React, { useEffect } from 'react';

const Welcome = () => {
  useEffect(() => {
    console.log('Welcome');
  }, []);
  return <div> Welcome </div>;
};

export default Welcome;
