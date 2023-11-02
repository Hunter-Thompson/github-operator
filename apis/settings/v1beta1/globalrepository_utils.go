package v1beta1

func (gr *GlobalRepository) GetDescription() *string {
	return gr.Spec.Description
}

func (gr *GlobalRepository) GetHomepage() *string {
	return gr.Spec.Homepage
}

func (gr *GlobalRepository) GetTopics() []string {
	return gr.Spec.Topics
}

func (gr *GlobalRepository) GetPrivate() *bool {
	return gr.Spec.Private
}

func (gr *GlobalRepository) GetHasIssues() *bool {
	return gr.Spec.HasIssues
}

func (gr *GlobalRepository) GetHasPages() *bool {
	return gr.Spec.HasPages
}

func (gr *GlobalRepository) GetHasProjects() *bool {
	return gr.Spec.HasProjects
}

func (gr *GlobalRepository) GetHasWiki() *bool {
	return gr.Spec.HasWiki
}

func (gr *GlobalRepository) GetAllowSquashMerge() *bool {
	return gr.Spec.AllowSquashMerge
}

func (gr *GlobalRepository) GetAllowMergeCommit() *bool {
	return gr.Spec.AllowMergeCommit
}

func (gr *GlobalRepository) GetAllowRebaseMerge() *bool {
	return gr.Spec.AllowRebaseMerge
}

func (gr *GlobalRepository) GetAllowAutoMerge() *bool {
	return gr.Spec.AllowAutoMerge
}

func (gr *GlobalRepository) GetDeleteBranchOnMerge() *bool {
	return gr.Spec.DeleteBranchOnMerge
}

func (gr *GlobalRepository) GetIsTemplate() *bool {
	return gr.Spec.IsTemplate
}
