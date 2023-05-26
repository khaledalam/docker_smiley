import React, {useState} from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import 'bootstrap/dist/css/bootstrap.css';
import { ToastContainer } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';
import {useSelector} from 'react-redux';

import {Provider} from 'react-redux';
import store from './store';
import * as action from './store/actions';

const root = ReactDOM.createRoot(document.getElementById('root'));



store.dispatch(action.authCheck());

root.render(

  <Provider store={store}>

      <App />
  </Provider>

  // For dev
  // <React.StrictMode>
  //   <ToastContainer />
  //   <App />
  // </React.StrictMode>
);
