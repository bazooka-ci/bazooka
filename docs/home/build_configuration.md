# Build Configuration

The entire configuration of your build is contained in a `.bazooka.yml` file in the root of your repository.
As the file extension suggests, a Bazooka build descriptor must be a valid YAML file.

This page details the configuration options available.

## Language

```yaml
language: go
```

You can choose any available programming language available in your Bazooka platform. Built-in languages are `go`, `java`, `nodejs` and `python`.
Other languages can be used if you install the necessary plugins.

Every language comes with its own specific fields, like the language version(s) for example...

Languages specifics are described in their own pages.

## Environment variables

You can inject any number of environment variables into your build using the `env` key:

```yaml
env:
  - X=42
  - Y=true
```

In the example above, during the build, two environment variables `X` and `Y` will be available and set to `42` and `true` respectively.

### Environment variables permutations

You can specify multiple values for a single environment variable by repeating the variable assignement.
Bazooka will the automatically generate all the possible [permutations](../home/permutations).

```yaml
env:
  - X=42
  - X=6
  - Y=true
```

For instance, with the previous configuration, bazooka will build your project twice. One time with the environment variables set to:

* X=42
* Y=true

And a second time with the environment variables set to:

* X=6
* Y=true

This allow you to make sure your project works with different configurations.

### Matrix

The `matrix` allows you to have a finer control over the generated [permutations](../home/permutations).
At this time, only `exclude` parameter is supported and lets you exclude some specific possible permutations.

```yaml
matrix:
  exclude:
    - go: 1.2.2
      env:
        - B=testb1
```


More details on the [permutation page](../home/permutations)

### Archiving

Archiving lets you catpure and store the build generated artifacts, like a `jar` file for example for a java project.

```yaml
archive:
  - target/*.jar
  - target/failsafe-reports/*.xml
```

The `archive` key takes one or multiple shell-compatible globs.
After the build, and for every build variant, any artifacts matching the specified globs will be captured and stored for later retrieval.

The captured artifacts can then be downloaded using the following API endpoint: `server:port/variant/{variant_id}/artifacts/{artifact_path}`, where:

* `{variant_id}` is the variant identifier
* `{artifact_path}` is the file path of the artifact

For example: `https://bzk.intranet:3000/variant/42c9d47fg/artifacts/unicorn.jar`.

With `archive`, the artifacts get captured whether the build succeeds or fails (but not when it is errored).

For finer control on when to capture build artifacts, bazooka also supports the following keys:

* `archive_success`: like with `archive`, takes one or multiple globs, but the maching artifacts are only captured if the build succeeds. Could be used to capture the deployable artifact which is only produced when the build succeeds
* `archive_failure`: like with `archive`, takes one or multiple globs, but the maching artifacts are only captured if the build fails. Could be used to capture the build system log and reports files when the build fails, e.g. the failsafe or surefire reports for java projects using maven as the build system

Here's an example:

```yaml
archive_success: target/*.jar
archive_failure: target/failsafe-reports/*.xml
```

### Services

Services allow you to have the ability to use external services within your build environment, such as databases...

```yaml
services:
  - mongodb
```

More details on the [services page](../home/services)
