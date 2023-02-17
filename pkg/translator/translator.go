package translator

type ITranslator interface {
	SendSuccessCallback() error
	Translate()
}
