import "../styles/login.css"
import { useState } from "react";

const Login = () => {

    const [selectedOption, setSelectedOption] = useState('option1');

    const handleOptionChange = (changeEvent) => {
        setSelectedOption(changeEvent.target.value);
    };

    return (
        <div>
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
                <label >
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
                
                <form method="post" action="http://localhost:8080/login">
                    <div>
                        <div>Логин</div>
                        <input type="text" name="username" pattern="^[A-Za-zА-Яа-яЁё\-\s]+$"></input>
                    </div>
                    <div>
                        <div>Пароль</div>
                        <input type="password" name="password" pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^\w\s]).{6,}"></input>
                        <div className="erol">"Пароль от 6 символов. Содержит: цифру, спец. символ, строчный и прописной символ латиницы"</div>
                    </div>
                    <button> Войти </button>
                </form>
            </div>
            <div className="LogDialog" id="Authorize">
                <h2>Авторизоваться</h2>
                <form method="post" action="http://localhost:8080/register">
                    <div>
                        <div>Имя</div>
                        <input type="text" name="name" pattern="^[A-Za-zА-Яа-яЁё\-\s]+$"></input>
                    </div>
                    <div>
                        <div>Фамилия</div>
                        <input type="text" name="surname" pattern="^[A-Za-zА-Яа-яЁё\-\s]+$"></input>
                    </div>
                    <div>
                        <div>Логин</div>
                        <input type="text" name="username"></input>
                    </div>
                    <div>
                        <div>Пароль</div>
                        <input type="password" name="password" pattern="^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^\w\s]).{6,}"></input>
                        <div className="erol">"Пароль от 6 символов. Содержит: цифру, спец. символ, строчный и прописной символ латиницы"</div>
                    </div>
                    <div>
                        <div>Староста?</div>
                        <input type="checkbox" name="is_admin" value={true}></input>
                    </div>
                    <button> Войти </button>
                </form>
            </div>
            
        </div>
        
    );
};

export default Login;