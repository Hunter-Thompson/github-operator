package v1beta1

func (team *Team) GetDescription() *string {
	if team.Spec.Description == nil {
		return new(string)
	}
	return team.Spec.Description
}

func (team *Team) GetMaintainers() []string {
	if team.Spec.Maintainers == nil {
		return []string{}
	}
	return *team.Spec.Maintainers
}

func (team *Team) GetRepoNames() []string {
	if team.Spec.RepoNames == nil {
		return []string{}
	}
	return *team.Spec.RepoNames
}
