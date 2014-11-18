set BZK_URI to the endpoint of Bazooka
e.g. : BZK_URI=http://192.168.59.103:3000

# Create a Project

```
> bzk project create --name testBuild --scm-uri git@bitbucket.org:bywan/bazooka-lang-example.git
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
bzk job start --project-id 5463801acc60e30001000001 --scm-ref golang
```

or, set BZK_PROJECT_ID to the ID of the project and run

```
bzk job start --scm-ref golang
```
