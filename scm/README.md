# Contract

## Input environment variables

* BZK_SCM_URL       : URL of the source code to fetch
* BZK_SCM_REFERENCE : Reference to fetch

## Input folders

### /bazooka-key (optional)

ssh private key that can be used to authenticate on the scm repo

## Output folders

### /bazooka

Source code of the project

### /meta

Should contain at least a file named `reference` which will contain the id of the last
commit on the branch that has just been fetched
