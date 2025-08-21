package contact

import "fmt"

const (
	ContactMessage = `Admin bilan chatga o‘tish uchun "Bog‘lanish" tugmasini bosing`
	ContactButton  = "Bog‘lanish"
)

const UserUrlTemplate = "https://t.me/%s"

var CreateUserLinkByUsername = func(userName string) string {
	return fmt.Sprintf(UserUrlTemplate, userName)
}
