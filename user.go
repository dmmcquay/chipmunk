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

func getUser(e string) (*user, error) {
	for _, i := range users {
		if e == i.Info.Email {
			return &i, nil
		}
	}
	return &user{}, fmt.Errorf("could not find user")
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
