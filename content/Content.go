package content

type Content = *string

var (
	TheContent Content
)

func init() {
	var theContent string
	TheContent = &theContent
}
