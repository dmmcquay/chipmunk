package chipmunk

import "fmt"

var authEmails []string = []string{"derekmcquay@gmail.com", "colleenmmcquay@gmail.com"}

type user struct {
	Info  userInfo `json:"info"`
	admin bool     `json:"admin"`
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

//authorizedEmail checks whether the email coming in is in the preapproved list
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

//addUser adds user to slice of users
func addUser(u userInfo) {
	_, err := getUser(u.Email)
	if err != nil {
		users = append(
			users,
			user{
				Info:  u,
				admin: true,
			},
		)
	}
}
