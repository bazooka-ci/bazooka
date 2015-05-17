package main

type BitbucketPayload struct {
	Commits []BitbucketCommit `json:"commits"`
}

type BitbucketCommit struct {
	RawNode string `json:"raw_node"`
}
