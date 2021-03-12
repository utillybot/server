package helpers

type ContextKey string
func (c ContextKey) String() string {
	return "utilly-server context key " + string(c)
}
