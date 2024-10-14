import React, {useState} from "react";
import {GoogleOAuthProvider, GoogleLogin} from '@react-oauth/google';

function Login() {
    const [player, setPlayer] = useState(null);

    const handleLoginSuccess = (response) => {
        const token = response.credential;

        fetch('http://localhost:8080/api/auth/google', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json', 
            },
            body: JSON.stringify({token}),   
        })
        .then(res => res.json())
        .then(data => {
            setPlayer(data);
        })
        .catch(error => console.error("Login failed: ", error))
    };

    return (

        <GoogleOAuthProvider clientId="1054033557847-44dp5pgltkb24hhn39514b5ubh97qjbs.apps.googleusercontent.com">
            <div>
                {player ? (
                    <div>
                        <h2>
                            Welcome {player.email}!
                        </h2>
                    </div>
                ) : (
                    <GoogleLogin 
                    onSuccess={handleLoginSuccess} 
                    onError={() => console.log("login failed :(")}
                    />
                )}
            </div>
        </GoogleOAuthProvider>
    );
};



export default Login;