import axios from 'axios';

const BASE_API_URL = "http://localhost:8080"; 

const BREAKPOINTS = {
    PING: "/ping",
    LOGIN: "/login",
    REGISTER: "/register",
}


export const fetchCurriencies = async () => {
    await axios.get(BASE_API_URL + BREAKPOINTS.PING).then(r => {
        console.log(r);
    })
}

export async function postLogin(data) {
    await axios.post(BASE_API_URL + BREAKPOINTS.LOGIN, data, {headers: {
        'Content-Type': 'application/json'
    }}).then(response => {console.log(response)});
}

export async function postAuthorize(data) {
    await axios.post(BASE_API_URL + BREAKPOINTS.REGISTER, data, {headers: {
        'Content-Type': 'application/json'
    }}).then(response => {console.log(response)});
}