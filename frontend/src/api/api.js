import axios from 'axios';

const BASE_API_URL = "http://localhost:8080"; 

const BREAKPOINTS = {
    PING: "/ping",
}

const API_METHODS = 
{
    GET: "GET",
    POST: "POST",
    PUT: "PUT",
    DELETE: "DELETE",
}


export const fetchCurriencies = async () => {
    await axios.get(BASE_API_URL + BREAKPOINTS.PING).then(r => {
        console.log(r);
    })
}

export const postLogin = () => {

}