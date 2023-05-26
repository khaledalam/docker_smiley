import {applyMiddleware, compose, createStore} from 'redux';
import {persistStore} from 'redux-persist';
import ReduxThunk from 'redux-thunk';
import RootReducer from './reducers';

const store = createStore(RootReducer, compose(
    applyMiddleware(ReduxThunk), (f) => f));

persistStore(store);

export default store;


// Dev:

// import { createStore, applyMiddleware } from 'redux';
// import { composeWithDevTools } from "redux-devtools-extension";
// import thunk from 'redux-thunk';
// import rootReducer from './reducers';
//
// const initialState = {};
//
// const middleware = [thunk];
//
// const store = createStore(
//     rootReducer,
//     initialState,
//     composeWithDevTools(applyMiddleware(...middleware))
// );
//
// export default store;
