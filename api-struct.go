package cmd

type Token struct {
  Token string `json:"cspAuthToken"`
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
