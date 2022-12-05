package authorization

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func hasPermission(p Permission, roles ...string) bool {
	if len(roles) == 0 {
		return false
	}

	for _, r := range roles {
		if globalAuthorizer.IsGranted(r, p, nil) {
			return true
		}
	}

	return false
}
