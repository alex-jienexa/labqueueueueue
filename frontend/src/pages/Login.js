import "../styles/login.css"
import { useState } from "react";
import {clickLogin, clickAuthorize} from "../scripts/click";
import { useNavigate } from 'react-router-dom'

const Login = () => {

    const [selectedOption, setSelectedOption] = useState('option1');

    const navigate = useNavigate();

    const handleOptionChange = (changeEvent) => {
        setSelectedOption(changeEvent.target.value);
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        
        try{
            const response = await clickLogin();
            navigate('/home');
        }
        catch (error){
            throw new Error(error);
        };
    };

    return (
        <div className="dialog">
            <div className="RadioHandle">
                <label>
                    <input
                    id="log_radio"
                    type="radio"
                    value="option1"
                    checked={selectedOption === 'option1'}
                    onChange={handleOptionChange}
                    />
                    Вход
                </label>
                <label>
                    <input
                    id="auth_radio"
                    type="radio"
                    value="option2"
                    checked={selectedOption === 'option2'}
                    onChange={handleOptionChange}
                    />
                    Регистрация
                </label>
            </div>
            <div className="LogDialog" id="LogIn">
                <h2>Войти в систему</h2>
                
                <div className="form">
                    <div className="handler">
                        <div>Логин</div>
                        <input type="text" name="username"></input>
                    </div>
                    <div className="handler">
                        <div>Пароль</div>
                        <input type="password" name="password" pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^\w\s]).{6,}"></input>
                        <div className="erol">"Пароль от 6 символов. Содержит: цифру, спец. символ, строчный и прописной символ латиницы"</div>
                    </div>
                    <button onClick={handleSubmit}> Войти </button>
                </div>
            </div>
            <div className="LogDialog" id="Authorize">
                <h2>Авторизоваться</h2>
                <div className="form">
                    <div className="handler">
                        <div>Имя</div>
                        <input type="text" name="name" pattern="^[A-Za-zА-Яа-яЁё\-\s]+$"></input>
                    </div>
                    <div className="handler">
                        <div>Фамилия</div>
                        <input type="text" name="surname" pattern="^[A-Za-zА-Яа-яЁё\-\s]+$"></input>
                    </div>
                    <div className="handler">
                        <div>Логин</div>
                        <input type="text" name="username"></input>
                    </div>
                    <div className="handler">
                        <div>Пароль</div>
                        <input type="password" name="password" pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^\w\s]).{6,}"></input>
                        <div className="erol">"Пароль от 6 символов. Содержит: цифру, спец. символ, строчный и прописной символ латиницы"</div>
                    </div>
                    <div className="handler">
                        <div>Староста?</div>
                        <input type="checkbox" name="is_admin" value={true}></input>
                    </div>
                    <button onClick={clickAuthorize}> Войти </button>
                </div>
                <div className="message" id="LogMessage"></div>
            </div>
            
        </div>
        
    );
};

export default Login;