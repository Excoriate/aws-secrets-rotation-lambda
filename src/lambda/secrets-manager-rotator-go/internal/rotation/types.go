package rotation

type Event struct {
	Token *string `json:"ClientRequestToken"`
	Arn   *string `json:"SecretId"`
	Step  *string `json:"Step"`
}

type Input struct {
	SecretName        string
	SecretARN         string
	IsRotationEnabled bool
}

type Rotation struct {
	SecretId     string
	VersionId    string
	VersionStage string
}

type DiscoveredSecrets struct {
	SecretName        string
	SecretARN         string
	SecretDescription string
	IsRotationEnabled bool
}

type SecretType struct {
	Db     string // TODO: To implement later.
	Static string
}
