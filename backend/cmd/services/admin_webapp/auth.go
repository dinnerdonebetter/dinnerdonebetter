package main

type userSessionDetails struct {
	Token     string `json:"token"`
	UserID    string `json:"userID"`
	AccountID string `json:"accountID"`
}
