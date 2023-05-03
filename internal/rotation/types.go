package rotation

type RotationEvent struct {
	Token *string `json:"ClientRequestToken"`
	Arn   *string `json:"SecretId"`
	Step  *string `json:"Step"`
}
