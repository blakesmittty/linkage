import React from "react";
import Login from "./Login";
import "./home.css"

function Home() {



    return (
        
        <div>
            <div className="login">
                <Login />
            </div>
            <div className="single-player">
                <button>Solo</button>
            </div>
            <div>
                <button>Versus</button>
            </div>
        </div>
       
    );
};

export default Home;