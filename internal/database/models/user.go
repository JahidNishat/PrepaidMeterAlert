package models

type Platform string

const (
	PlatformTelegram Platform = "telegram"
)

type User struct {
	BaseTimeStampedModel

	Platform   Platform `bun:"platform,notnull"`
	PlatformID string   `bun:"platform_id,notnull"`
}
