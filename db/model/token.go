package model

// Token is a representation of a row in the tokens entity
type Token struct {
    ID           string `json:"auth_id" db:"id"`
    AccessToken  string `json:"access_token" db:"access_token"`
    TokenType    string `json:"token_type" db:"token_type"`
    Scope        string `json:"scope" db:"scope"`
    ExpiresIn    int    `json:"expires_in" db:"expires_in"`
    RefreshToken string `json:"refresh_token" db:"refresh_token"`
    CreatedAt    string `json:"created_at" db:"created_at"`
}
