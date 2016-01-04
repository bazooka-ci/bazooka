set BZK_API_URI to the endpoint of Bazooka
e.g. : BZK_API_URI=http://192.168.59.103:3000

# Create a Project

```
> bzk project create <name> <scm> <scm_url>
PROJECT ID                 NAME           SCM TYPE       SCM URI
546b1f2be0973d0001000001   testBuild2     git            git@bitbucket.org:bywan/bazooka-lang-example.git
```

# List all projects

```
bzk project list
PROJECT ID                 NAME           SCM TYPE       SCM URI
546b1f2be0973d0001000001   testBuild2     git            git@bitbucket.org:bywan/bazooka-lang-example.git
55693233eade456001000001   testVagrant    git            git@github.com:ggiamarchi/vagrant-openstack-provider.git
```

# Start a Job

```
bzk job start <project_id> <scm_ref>
```

# List variants for a job

```
bzk variant list <project_id> <job_id>
```
