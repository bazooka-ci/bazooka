# Clone a git repo with Docker

## Without Authentication

```bash
sudo docker run -v /home/user/dest:/bazooka \
                bazooka/scm-git $GIT_URL $GIT_BRANCH
```

## With SSH key Authentication

```bash
sudo docker run -v /home/user/dest:/bazooka \
                -v /home/user/.ssh/id_rsa:/bazooka-key \
                bazooka/scm-git $GIT_URL $GIT_BRANCH
```
