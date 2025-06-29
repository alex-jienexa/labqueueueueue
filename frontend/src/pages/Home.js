import "../styles/home.css";
import Hat from "../elements/Hat.js"
import "../styles/Hat.css"

const Home = () => {
    return (
        <div>
            <Hat/>
            <div id="HomePage"> 
                <div id="homeNavigation">
                    <div>Элемент 1</div>
                    <div>Элемент 2</div>
                    <div>Элемент 3</div>
                </div>
                <div id="queueList">
                    <div className="queue">
                        Очередь
                    </div>
                </div>
            </div>

        </div>
        
    );
};

export default Home;