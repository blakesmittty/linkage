import React from "react";
import Login from "./Login";
import "./home.css"
import { useNavigate } from 'react-router-dom';


function Home() {
    const navigate = useNavigate();

    const goToGameSolo = () => {
        navigate("/game")
    };

    return (
        
        <div className="home">
            <h1>Linkage!</h1>
            <div className="login">
                <Login />
            </div>
            <div className="single-player">
                <button onClick={goToGameSolo}>Solo</button>
            </div>
            <div>
                <button>Versus</button>
            </div>
        </div>
       
    );
};

export default Home;