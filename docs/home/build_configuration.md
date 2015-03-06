# Build Configuration

The entire configuration of your build is contained in a `.bazooka.yml` file at the root of your repository. This page details the configuration options available.

## Language

```yaml
language: go
```

You can choose any available programming language available in your Bazooka platform. Built-in languages are `go`, `java`, `nodejs` and `python`. If any other plugin has been added to your Bazooka platform, you can use it your configuration file. Each language comes with its own specific fields, such as languages versions... Languages specifics are described on their own pages.

## Environment variables

```yaml
env:
  - X=42
  - Y=true
```

You can declare any environment variables to be available for your build scripts. In this case, the environment variables `X` and `Y` are defined

### Environment variables permutations

```yaml
env:
  - X=42
  - X=6
  - Y=true
```

You can specify multiple values for a single environment variable. Bazooka will automatically generate [permutations](../home/permutations). For instance, with the previous configuration, bazooka will build your project twice. One time with the environment variables set to:

* X=42
* Y=true

Another time with the environment variables set to:

* X=6
* Y=true

This allow you to make sure your project works with different configurations

### Matrix

```yaml
matrix:
  exclude:
    - go: 1.2.2
      env:
        - B=testb1
```

Matrix allows you to manage your [permutations](../home/permutations) easily. Currently only `exclude` parameter is supported

More details on the [permutation page](../home/permutations)

### Services

```yaml
services:
  - mongodb
```

Services allow you to have the ability to use external services within your build environment, such as databases...

More details on the [services page](../home/services)
