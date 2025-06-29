import "../styles/Hat.css";
import Image from "./icon.png"
import { Link } from 'react-router-dom';

const Hat = () =>{
    return (
    <nav>
        <ul>
            <li className="icon">
                <Link to={"https://www.youtube.com/watch?v=dQw4w9WgXcQ"}>
                    <img src={Image}></img>
                </Link>
            </li>
        </ul>
        <div className="exit">
            <Link to={"/login"}>exit</Link>
        </div>
        
    </nav>
)};

export default Hat;