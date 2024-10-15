package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// may use localstorage for high scores and information so

type TokenInfo struct {
	Issuer        string `json:"iss"`
	Subject       string `json:"sub"`
	Audience      string `json:"aud"`
	Expiry        string `json:"exp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

// NEED TO PERSIST A SESSION VIA SESSION COOKIE

func verifyGoogleToken(token string) (*TokenInfo, error) {
	// make url to authorize token with google
	googleOAuthURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token)
	// send a get request to the url with the token
	resp, err := http.Get(googleOAuthURL)
	if err != nil {
		log.Printf("Error making request to Google: %v", err)
		return nil, err
	}
	defer resp.Body.Close() // close the response body when the function ends

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Google token verification failed with status: %s, body: %s", resp.Status, body)
		return nil, fmt.Errorf("failed to verify token: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenInfo TokenInfo
	err = json.Unmarshal(body, &tokenInfo) // unmarshal token data into the tokenInfo struct
	if err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	// decode the http request as json into the struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// authorize the token
	tokenInfo, err := verifyGoogleToken(req.Token)
	if err != nil {
		log.Printf("Token verification error: %v", err)
		http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("Google User ID: %s, Email: %s\n", tokenInfo.Subject, tokenInfo.Email)

	// capture the index of @ to give player a username
	atIndex := strings.Index(tokenInfo.Email, "@")

	// set cookies to persist user session so they dont need to login every refresh
	/*
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    tokenInfo.Subject,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour), // Adjust session duration
			HttpOnly: true,
			Secure:   false, // Set to true in production with HTTPS
		})
	*/

	// create a response to be sent to the client with their info
	response := map[string]string{
		"userID":   tokenInfo.Subject,
		"username": tokenInfo.Email[0:atIndex],
		"email":    tokenInfo.Email,
		"message":  "User authenticated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
