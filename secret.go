package marathon

type SecretSource struct {
	Source string `json:"source"`
}

type EnvironmentSecret struct {
	Secret string `json:"secret"`
}
