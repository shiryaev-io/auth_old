package responses

// Содержит статус код и тело ответа
type Common struct {
	Status int
	Body   []byte
}
