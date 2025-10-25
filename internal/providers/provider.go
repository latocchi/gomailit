package providers

type Provider interface {
	SendEmail() error
}
