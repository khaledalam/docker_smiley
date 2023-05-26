import * as ActionTypes from '../actions/types';
import {UPDATE_DARK_MODE} from "../actions/types";

const darkMode = localStorage.getItem('darkMode') === '1';

const initialState = {
    globalInfo: false,
    globalError: false,
    globalSuccess: false,
    darkMode: darkMode
};

const updateApp = (state, payload) => {
    return {
        ...state,
        ...payload
    };
};


const App = (state = initialState, {type, payload = null}) => {

    switch (type) {
        case ActionTypes.APP_GLOBAL_MESSAGE:
        case ActionTypes.UPDATE_DARK_MODE:
            return updateApp(state, payload);
        default:
            return state;
    }
};

export default App;
