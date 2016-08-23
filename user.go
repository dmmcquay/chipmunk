package chipmunk

import "fmt"

var authEmails []string = []string{"derekmcquay@gmail.com", "colleenmmcquay@gmail.com", "dmmllnl@gmail.com"}

type user struct {
	Info  userInfo `json:"info"`
	admin bool     `json:"admin"`
	txs   []tranx  `json:"Txs"`
}

type userInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (u *user) addTranx(t tranx) {
	u.txs = append(u.txs, t)
}

func authorizedEmail(e string) bool {
	b := false
	for _, i := range authEmails {
		if i == e {
			b = true
		}
	}
	return b
}

// getUser returns index of user with given email, otherwise it returns an
// error that it could not find that user
func getUser(e string) (int, error) {
	for i, u := range users {
		if e == u.Info.Email {
			return i, nil
		}
	}
	return 0, fmt.Errorf("could not find user")
}

func addUser(u userInfo) {
	users = append(
		users,
		user{
			Info:  u,
			admin: true,
		},
	)
}
