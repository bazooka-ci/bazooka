package main

type ConfigGolang struct {
	Language      string   "language"
	Setup         []string "setup,omitempty"
	BeforeInstall []string "before_install,omitempty"
	Install       []string "install,omitempty"
	BeforeScript  []string "before_script,omitempty"
	Script        []string "script,omitempty"
	AfterScript   []string "after_script,omitempty"
	AfterSuccess  []string "after_success,omitempty"
	AfterFailure  []string "after_failure,omitempty"
	Services      []string "services,omitempty"
	Env           []string "env,omitempty"
	GoVersions    []string "go,omitempty"
	FromImage     string   "from"
}
