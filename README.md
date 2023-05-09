# github-operator
An operator to manage policy-as-code and apply repository settings to repositories across an organization.

## Guide

### Deployment

1. Create GitHub Token with access to `repo`, `write:org`, and `read:org`
2. Deploy on Kubernetes cluster:

```sh
git clone github.com/Hunter-Thompson/github-operator
cd github-operator
kubectl create secret generic github-operator-secret --from-literal GITHUB_TOKEN=ghp_sjdnfksdfnjksdf 
kubectl apply -f docs/install.yaml
```

### Team

```yaml
# https://docs.github.com/en/rest/teams/teams?apiVersion=2022-11-28#create-a-team
# https://docs.github.com/en/rest/teams/teams?apiVersion=2022-11-28#update-a-team

apiVersion: settings.github.com/v1beta1
kind: Team
metadata:
  name: team-sample
spec:
  organization: ocktokit
  description: 'asd'
  maintainers:
  - Hunter-Thompson
  repoNames:
  - repo1
  - anotherrepo
  privacy: secret
  parentTeamId: 1829
```

### Repository

```yaml
# https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#create-an-organization-repository
# https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#update-a-repository

apiVersion: settings.github.com/v1beta1
kind: Repository
metadata:
  name: repo-sample
spec:
  organization: ocktokit
  description: "asd"
  homepage: "https://home.page"
  topics:
  - topic1
  - topic2
  private: false
  hasIssues: false
  hasPages: false
  hasWiki: false
  defaultBranch: master
  allowSquashMerge: true
  allowMergeCommit: true
  allowRebaseMerge: true
  allowAutoMerge: true
  deleteBranchOnMerge: true
  isTemplate: false
  squashMergeCommitTitle: "PR_TITLE"
  squashMergeCommitMessage: "COMMIT_MESSAGES"
  mergeCommitTitle: "PR_TITLE"
  mergeCommitMessage: "PR_BODY"
```

## License

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

