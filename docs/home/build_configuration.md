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

### Services

Services allow you to have the ability to use external services within your build environment, such as databases...

```yaml
services:
  - mongodb
```

More details on the [services page](../home/services)
