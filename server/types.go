package main

type BitbucketPayload struct {
	Commits []BitbucketCommit `json:"commits"`
}

type BitbucketCommit struct {
	RawNode string `json:"raw_node"`
}

type GithubPayload struct {
	HeadCommit GithubCommit `json:"head_commit"`
}

type GithubCommit struct {
	ID string `json:"id"`
}
