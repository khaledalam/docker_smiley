import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import 'bootstrap/dist/css/bootstrap.css';
import { ToastContainer } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  
  <div>
    <ToastContainer />
    <App />
  </div>

  // For dev
  // <React.StrictMode>
  //   <ToastContainer />
  //   <App />
  // </React.StrictMode>
);
