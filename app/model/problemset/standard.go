package problemset

type standardProblemsetUserRoleInfo struct {
	IsGuest  bool `json:"is_guest"`
	IsMember bool `json:"is_member"`
	IsAdmin  bool `json:"is_admin"`
}
type standardProblemsetInfo struct {
	IsPublic bool `json:"is_public"`
}
type standardProblemsetProvider struct{}

func (standardProblemsetProvider) GetDefaultProblemsetInfo() ProblemsetInfo {
	return new(standardProblemsetInfo)
}

func (*standardProblemsetInfo) GetCreatorRole() ProblemsetUserRole {
	return &standardProblemsetUserRoleInfo{IsAdmin: true, IsGuest: false, IsMember: true}
}

func (*standardProblemsetInfo) GetGuestRole() ProblemsetUserRole {
	return &standardProblemsetUserRoleInfo{IsAdmin: false, IsGuest: false, IsMember: false}
}

func (*standardProblemsetInfo) GetRegisteredUserRole() ProblemsetUserRole {
	return &standardProblemsetUserRoleInfo{IsAdmin: false, IsGuest: true, IsMember: false}
}

func (*standardProblemsetInfo) GetDefaultRole() ProblemsetUserRole {
	return &standardProblemsetUserRoleInfo{IsAdmin: false, IsGuest: false, IsMember: false}
}

func (ps *standardProblemsetInfo) CheckPrivilege(u_ ProblemsetUserRole, p ProblemsetPrivilege) error {
	u := u_.(*standardProblemsetUserRoleInfo)
	if u.IsAdmin {
		return nil
	}
	switch p {
	case ProblemsetCreateProblemPrivilege:
		return ProblemsetPermissionDeniedError
	case ProblemsetViewProblemPrivilege:
		if ps.IsPublic {
			return nil
		} else {
			if u.IsMember {
				return nil
			}
		}
		return ProblemsetPermissionDeniedError
	case ProblemsetSubmitProblemPrivilege:
		if u.IsGuest {
			return ProblemsetPermissionDeniedError
		}
		if ps.IsPublic {
			return nil
		} else {
			if u.IsMember {
				return nil
			}
		}
		return nil
	}
	panic(ProblemsetPermissionInvalidError)
}
