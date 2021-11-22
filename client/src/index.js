import React from 'react';
import ReactDOM from 'react-dom';

import "./library/axios.js";
import "./library/toaster.js";
import App from './App';

ReactDOM.render(
  <React.StrictMode>
    <App/>
  </React.StrictMode>,
  document.getElementById('root')
);