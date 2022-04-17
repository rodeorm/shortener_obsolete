package repo

type URL struct {
	Key string `json:"url,omitempty"`
}

type ShortenURL struct {
	Key string `json:"result,omitempty"`
}

type URLPair struct {
	Origin string `json:"origin,omitempty"`
	Short  string `json:"short,omitempty"`
}
