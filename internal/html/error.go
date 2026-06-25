package html

type HtmlError struct {
	Message string
}

func (error *HtmlError) Error() string {
	return error.Message
}
