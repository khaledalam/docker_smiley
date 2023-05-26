import * as ActionTypes from './types';

export function authLogin(payload) {
    return {
        type: ActionTypes.AUTH_LOGIN,
        payload,
    };
}

export function authLogout() {
    return {
        type: ActionTypes.AUTH_LOGOUT,
    };
}

export function authCheck() {
    return {
        type: ActionTypes.AUTH_CHECK,
    };
}


export function appGlobalMessage(payload) {
    return {
        type: ActionTypes.APP_GLOBAL_MESSAGE,
        payload
    }
}

export function appUpdateDarkMode(payload) {
    return {
        type: ActionTypes.UPDATE_DARK_MODE,
        payload
    }
}