package cmd

type Token struct {
  Token string `json:"cspAuthToken"`
}

type vRACToken struct {
  vRACToken string `json:"token"`
}

type Catalog struct {
  Catalog []Item `json:"content"`
}

type Item struct {
  Name string `json:"name"`
	ID   string `json:"id"`
}

type Projects struct {
  Project []Project `json:"content"`
}

type Project struct {
  Name string `json:"name"`
	ID   string `json:"id"`
}

type Deployments struct {
  Deployment []Deployment `json:"content"`
}

type Deployment struct {
  Name string `json:"name"`
	ID   string `json:"id"`
}

type Action struct {
  Name string `json:"name"`
	ID   string `json:"id"`
}

type AuthenticationRequestCloud struct {
	RefreshToken string `json:"refreshToken"`
}

// AuthenticationResponseCloud - Authentication response structure for Cloud
type AuthenticationResponseCloud struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

type AuthenticationError struct {
	Timestamp     int64  `json:"timestamp"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Error         string `json:"error"`
	ServerMessage string `json:"serverMessage"`
}
