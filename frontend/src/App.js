import {BrowserRouter, Navigate, Route, Routes, Router, useLocation, useNavigate} from "react-router-dom";
import './App.css';

import Login from "./pages/Login"
import Home from "./pages/Home"
import Queue from "./pages/Queue"
import Create from "./pages/Create"
import { useEffect } from "react";

export const ROUTES = {
    LOGIN: "/login",
    HOME: "/home",
    QINFO: "/queue",
    CREATE: "/qcreate"
}

function checkObj(obj, check_val){
  return Object.values(obj).includes(check_val)
}


function AppContext() {

  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (!checkObj(ROUTES, location.pathname)) {
      navigate(-1);
    }
  });


  return (
    <div className="App">
      <div>
        <Routes>
          <Route exec path={ROUTES.LOGIN} element={<Login />} />
          <Route path={ROUTES.HOME} element={<Home />} />
          <Route path={ROUTES.QINFO} element={<Queue />} />
          <Route path={ROUTES.CREATE} element={<Create />} />
        </Routes>
      </div>
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AppContext />
    </BrowserRouter>
  );
}

export default App;
