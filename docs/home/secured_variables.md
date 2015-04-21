# Secured environment variables

Some configuration of your build need to be secured, and not accessible from other developers, or even for people who have read access to your repository.

For instance, login or password, to push artifacts to a repository manager, need to be accessible within your build, but also need to stay private.

Bazooka allows you to encrypt some of your environment variables, put the encrypted data in your `.bazooka.yml` file, and they will decrypted by the server.

Each encrypted data is specific to a project, and can only be decrypted by this specific project.

## Encrypt

The cli comes with a new verb to encrypt data

```bash
> bzk encrypt --help

Usage: bzk encrypt PROJECT_ID DATA

Encrypt some data

Arguments:
  PROJECT_ID=""   Project id
  DATA=""         Data to Encrypt
```

For instance, to encrypt the environment variable `PASSWORD=secret` for the project `my_project`

```bash
>  bzk encrypt my_project PASSWORD=secret
Encrypted data: (to add to your .bazooka.yml file)
d68e754fdcd7defdaaf02c585b4615e905659272dd1d862b48237ee55a601d9ad78a0ea3
```

The result of this command needs to be added to the `.bazooka.yml` of your project

## .bazooka.yml configuration

To add an encrypted environment variable to your build, add it in your `.bazooka.yml` like this

```yaml
env:
  - X=42
  - secure: d68e754fdcd7defdaaf02c585b4615e905659272dd1d862b48237ee55a601d9ad78a0ea3
```

During your build, the data will be decrypted by the server, and the env var (for instance PASSWORD), will be accessible as any other variable.
