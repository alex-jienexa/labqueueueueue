import {BrowserRouter, Navigate, Route, Routes, Router, useLocation, useNavigate} from "react-router-dom";
import './App.css';
import Login from "./pages/Login"
import Home from "./pages/Home"
import Queue from "./pages/Queue"
import Create from "./pages/Create"

export const ROUTES = {
    LOGIN: "/login",
    HOME: "/home",
    QINFO: "/queue",
    CREATE: "/qcreate"
}


function App() {

  return (
    <div className="App">
      <BrowserRouter>
        <div>
          <Routes>
            <Route exec path={ROUTES.LOGIN} element={<Login />} />
            <Route path={ROUTES.HOME} element={<Home />} />
            <Route path={ROUTES.QINFO} element={<Queue />} />
            <Route path={ROUTES.CREATE} element={<Create />} />
          </Routes>
        </div>
      </BrowserRouter>
    </div>
  );
}

export default App;
