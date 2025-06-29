import "../styles/home.css";
import Hat from "../elements/Hat.js"
import "../styles/Hat.css";
import { Link } from 'react-router-dom';

const Create = () => {
    return (
        <div>
            <Hat/>
            <div id="HomePage"> 
                <div id="homeNavigation">
                    <ul>
                        <li>
                            <Link to="/home">Просмотреть очереди</Link>
                        </li>
                        <li>
                            <Link to="/home/create">Создать Очередь</Link>
                        </li>
                    </ul>
                </div>
                <div id="queueList">
                    <div className="queue">
                        ФОрма для создания очереди
                    </div>
                </div>
            </div>

        </div>
    );
};

export default Create;