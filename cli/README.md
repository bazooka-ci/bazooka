set BZK_URI to the endpoint of Bazooka
e.g. : BZK_URI=http://192.168.59.103:3000

# Create a Project

```
bzk project create --name testBuild --scm-uri git@bitbucket.org:bywan/bazooka-lang-example.git
```

# List all projects

```
bzk project list
```

# Start a Job

```
bzk project start-job --project-id 5463801acc60e30001000001 --scm-ref golang
```
