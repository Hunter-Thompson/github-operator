package v1beta1

func (repo *Repository) GetDescription() *string {
	if repo.Spec.Description == nil {
		return new(string)
	}
	return repo.Spec.Description
}

func (repo *Repository) GetHomepage() *string {
	if repo.Spec.Homepage == nil {
		return new(string)
	}
	return repo.Spec.Homepage
}

func (repo *Repository) GetTopics() []string {
	if repo.Spec.Topics == nil {
		return []string{}
	}
	return *repo.Spec.Topics
}

func (repo *Repository) GetPrivate() *bool {
	if repo.Spec.Private == nil {
		return boolPointer(true)
	}
	return repo.Spec.Private
}

func (repo *Repository) GetHasIssues() *bool {
	if repo.Spec.HasIssues == nil {
		return boolPointer(true)
	}
	return repo.Spec.HasIssues
}

func (repo *Repository) GetHasPages() *bool {
	if repo.Spec.HasPages == nil {
		return boolPointer(true)
	}
	return repo.Spec.HasPages
}

func (repo *Repository) GetHasProjects() *bool {
	if repo.Spec.HasProjects == nil {
		return boolPointer(true)
	}
	return repo.Spec.HasProjects
}

func (repo *Repository) GetHasWiki() *bool {
	if repo.Spec.HasWiki == nil {
		return boolPointer(true)
	}
	return repo.Spec.HasWiki
}

func (repo *Repository) GetAllowSquashMerge() *bool {
	if repo.Spec.AllowSquashMerge == nil {
		return boolPointer(true)
	}
	return repo.Spec.AllowSquashMerge
}

func (repo *Repository) GetAllowMergeCommit() *bool {
	if repo.Spec.AllowMergeCommit == nil {
		return boolPointer(true)
	}
	return repo.Spec.AllowMergeCommit
}

func (repo *Repository) GetAllowRebaseMerge() *bool {
	if repo.Spec.AllowRebaseMerge == nil {
		return boolPointer(true)
	}
	return repo.Spec.AllowRebaseMerge
}

func (repo *Repository) GetAllowAutoMerge() *bool {
	if repo.Spec.AllowAutoMerge == nil {
		return boolPointer(false)
	}
	return repo.Spec.AllowAutoMerge
}

func (repo *Repository) GetDeleteBranchOnMerge() *bool {
	if repo.Spec.DeleteBranchOnMerge == nil {
		return boolPointer(false)
	}
	return repo.Spec.DeleteBranchOnMerge
}

func (repo *Repository) GetIsTemplate() *bool {
	if repo.Spec.IsTemplate == nil {
		return boolPointer(false)
	}
	return repo.Spec.IsTemplate
}

func boolPointer(bool bool) *bool {
	return &bool
}
