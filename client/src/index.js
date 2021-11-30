import React from 'react';
import ReactDOM from 'react-dom';

import "./library/axios.js";
import "./library/toaster.js";
import App from './App';

const renderApp = (data = false) => {
  ReactDOM.render(
    <React.StrictMode>
      <App auth={data}/>
    </React.StrictMode>,
    document.getElementById('root')
  );
}

if(localStorage.getItem("user-token")){
  window.$axios.get("/me")
  .then(res => {    
    renderApp(res.data)
  })
  .catch(err => {
    localStorage.removeItem("user-token")
    renderApp(false)
  })
}else{
  renderApp(false)
}