package v1beta1

func (repo *Repository) GetDescription() *string {
	return repo.Spec.Description
}

func (repo *Repository) GetHomepage() *string {
	return repo.Spec.Homepage
}

func (repo *Repository) GetTopics() []string {
	return repo.Spec.Topics
}

func (repo *Repository) GetPrivate() *bool {
	return repo.Spec.Private
}

func (repo *Repository) GetHasIssues() *bool {
	return repo.Spec.HasIssues
}

func (repo *Repository) GetHasPages() *bool {
	return repo.Spec.HasPages
}

func (repo *Repository) GetHasProjects() *bool {
	return repo.Spec.HasProjects
}

func (repo *Repository) GetHasWiki() *bool {
	return repo.Spec.HasWiki
}

func (repo *Repository) GetAllowSquashMerge() *bool {
	return repo.Spec.AllowSquashMerge
}

func (repo *Repository) GetAllowMergeCommit() *bool {
	return repo.Spec.AllowMergeCommit
}

func (repo *Repository) GetAllowRebaseMerge() *bool {
	return repo.Spec.AllowRebaseMerge
}

func (repo *Repository) GetAllowAutoMerge() *bool {
	return repo.Spec.AllowAutoMerge
}

func (repo *Repository) GetDeleteBranchOnMerge() *bool {
	return repo.Spec.DeleteBranchOnMerge
}

func (repo *Repository) GetIsTemplate() *bool {
	return repo.Spec.IsTemplate
}
