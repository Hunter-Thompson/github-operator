/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GlobalRepositorySpec defines the desired state of GlobalRepository
type GlobalRepositorySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// GitHub Org
	Organization string `json:"organization"`
	// A short description of the repository that will show up on GitHub
	Description *string `json:"description,omitempty"`
	// A URL with more information about the repository
	Homepage *string `json:"homepage,omitempty"`
	// A list of topics to set on the repository - can alternatively set like this: [github, probot, new-topic, another-topic, topic-12]
	Topics []string `json:"topics,omitempty"`
	// Either `true` to make the repository private, or `false` to make it public.
	Private *bool `json:"private,omitempty"`
	// Either `true` to enable issues for this repository, `false` to disable them.
	HasIssues *bool `json:"hasIssues,omitempty"`
	// Either `true` to enable pages for this repository, `false` to disable them.
	HasPages *bool `json:"hasPages,omitempty"`
	// Either `true` to enable projects for this repository, or `false` to disable them.
	HasProjects *bool `json:"hasProjects,omitempty"`
	// Either `true` to enable the wiki for this repository, `false` to disable it.
	HasWiki *bool `json:"hasWiki,omitempty"`
	// The default branch for this repository.
	DefaultBranch *string `json:"defaultBranch,omitempty"`
	// Either `true` to allow squash-merging pull requests, or `false` to prevent squash-merging.
	AllowSquashMerge *bool `json:"allowSquashMerge,omitempty"`
	// Either `true` to allow merging pull requests with a merge commit, or `false` to prevent merging pull requests with merge commits.
	AllowMergeCommit *bool `json:"allowMergeCommit,omitempty"`
	// Either `true` to allow merging pull requests with a merge commit, or `false` to prevent merging pull requests with merge commits.
	AllowRebaseMerge *bool `json:"allowRebaseMerge,omitempty"`
	// Either `true` to allow auto-merge on pull requests, or `false` to disallow auto-merge.
	AllowAutoMerge *bool `json:"allowAutoMerge,omitempty"`
	// Either `true` to allow automatically deleting head branches, when pull requests are merged, or `false` to prevent automatic deletion.
	DeleteBranchOnMerge *bool `json:"deleteBranchOnMerge,omitempty"`
	// Either `true` if its a template repo, false if its not
	IsTemplate *bool `json:"isTemplate,omitempty"`
	// Can be one of: PR_TITLE, COMMIT_OR_PR_TITLE
	SquashMergeCommitTitle *string `json:"squashMergeCommitTitle,omitempty"`
	// Can be one of: PR_BODY, COMMIT_MESSAGES, BLANK
	SquashMergeCommitMessage *string `json:"squashMergeCommitMessage,omitempty"`
	// Can be one of: PR_TITLE, MERGE_MESSAGE
	MergeCommitTitle *string `json:"mergeCommitTitle,omitempty"`
	// Can be one of: PR_BODY, PR_TITLE, BLANK
	MergeCommitMessage *string `json:"mergeCommitMessage,omitempty"`

	RepositoryCollaborators *RepositoryCollaborators `json:"repositoryCollaborators,omitempty"`
}

// GlobalRepositoryStatus defines the observed state of GlobalRepository
type GlobalRepositoryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GlobalRepository is the Schema for the globalrepositories API
type GlobalRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalRepositorySpec   `json:"spec,omitempty"`
	Status GlobalRepositoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlobalRepositoryList contains a list of GlobalRepository
type GlobalRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalRepository{}, &GlobalRepositoryList{})
}
