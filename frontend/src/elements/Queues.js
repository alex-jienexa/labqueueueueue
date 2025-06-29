import { useEffect, useState } from "react";
import { BREAKPOINTS,  BASE_API_URL} from "../api/api";
import axios from "axios";


function Queues() {

    const [currencies, setCurriencies] = useState();

    const getQueues = () => {
        const token = localStorage.getItem('jwt');
        axios.get(BASE_API_URL + BREAKPOINTS.QUEUES, {
            headers: {Authorization: `${token}`}
        }).then(r => {
            const currenciesResponce = r.data;
            console.log(r.data);
        })
    }

    useEffect( () => {
        // getQueues();
    })

    return (
       <div>
       </div>
     )
};

export default Queues;