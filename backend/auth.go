package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TokenInfo struct {
	Issuer        string `json:"iss"`
	Subject       string `json:"sub"`
	Audience      string `json:"aud"`
	Expiry        string `json:"exp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
}

func verifyGoogleToken(token string) (*TokenInfo, error) {
	googleOAuthURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token)
	resp, err := http.Get(googleOAuthURL)
	if err != nil {
		log.Printf("Error making request to Google: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

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
	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tokenInfo, err := verifyGoogleToken(req.Token)
	if err != nil {
		log.Printf("Token verification error: %v", err)
		http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("Google User ID: %s, Email: %s\n", tokenInfo.Subject, tokenInfo.Email)

	response := map[string]string{
		"userID":  tokenInfo.Subject,
		"email":   tokenInfo.Email,
		"message": "User authenticated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
