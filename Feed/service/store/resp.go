package store

type registerResponse struct {
	id string
}

func (r registerResponse) GetId() string {
	return r.id
}

type loginResponse struct {
	token    string
	userInfo UserInfo
}

func (l loginResponse) GetToken() string {
	return l.token
}
func (l loginResponse) GetUserInfo() UserInfo {
	return l.userInfo
}

type userInfo struct {
	lastname  string
	firstname string
	username  string
	email     string
	id        string
}

func (u userInfo) GetLastName() string {
	return u.lastname
}

func (u userInfo) GetFirstName() string {
	return u.firstname
}

func (u userInfo) GetUserName() string {
	return u.username
}

func (u userInfo) GetEmail() string {
	return u.email
}

func (u userInfo) GetId() string {
	return u.id
}

type feedIder struct {
	FeedInfo
	id string
}

func (f feedIder) GetId() string {
	return f.id
}
