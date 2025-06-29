import axios from 'axios';
import CryptoJS from 'crypto-js';

export const BASE_API_URL = "http://localhost:8080"; 

export const BREAKPOINTS = {
    PING: "/ping",
    JPING: "/rping",
    LOGIN: "/login",
    REGISTER: "/register",
    ACTIVE_QUEUES: "/queues/active",
    QUEUES: "/queues"
}

export const fetchCurriencies = async () => {
    await axios.get(BASE_API_URL + BREAKPOINTS.PING).then(r => {
        console.log(r);
    })
}

async function jwtPing (token) {
    await axios.post(BASE_API_URL + BREAKPOINTS.JPING, {}, {
        headers: {Authorization: `${token}`}
    }).then(r => {
        console.log(r);
    })
}

export async function postLogin(data) {
    let result;
    await axios.post(BASE_API_URL + BREAKPOINTS.LOGIN, data, {headers: {
        'Content-Type': 'application/json'
    }}).then(response => {
        const athorizeResponce = response.data;
        const token = athorizeResponce.token;

        console.log(token);

        localStorage.setItem('jwt', token);
;
        result = response;
    });
    return result;
}

export async function postAuthorize(data) {
    await axios.post(BASE_API_URL + BREAKPOINTS.REGISTER, data, {headers: {
        'Content-Type': 'application/json'
    }}).then(response => {
        const cResponce = response.data;

        console.log(cResponce);

    });
}