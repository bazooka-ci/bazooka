package main

type BitbucketPayload struct {
	Push BitbucketPush `json:"push"`
}

type BitbucketPush struct {
	Changes []BitbucketChanges `json:"changes"`
}

type BitbucketChanges struct {
	New BitbucketNew `json:"new"`
}

type BitbucketNew struct {
	Type   string          `json:"type"`
	Name   string          `json:"name"`
	Target BitbucketTarget `json:"hash"`
}

type BitbucketTarget struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
}
