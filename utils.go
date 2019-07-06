package vos

func home(user string) string {
	if user == "root" {
		return "/root"
	}
	return "/home/" + user
}

// BitIf x86_64/x86
func BitIf() string {
	if 32<<(^uint(0)>>63) == 64 {
		return "x86_64"
	}
	return "x86"
}
