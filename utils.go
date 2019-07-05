package vos

func home(user string) string {
	if user == "root" {
		return "/root"
	}
	return "/home/" + user
}
