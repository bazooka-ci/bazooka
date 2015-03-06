# Build Lifecycle

## Lifecycle phases

Bazooka provides default build phases for each programming languages. For instance, Java runners will try to use Maven is a `pom.xml` is available in your repository. But any step is customizable.

The build phases are :

* `before_install`
* `install`
* `before_script`
* `script`
* `after_success`
* `after_failure`
* `after_script`

### `before_install`

In a before_install step, you can install additional dependencies required by your project, Ubuntu packages for instance, or custom services, downloaded and installed from the internet.

### `install`

In this step, you will install any dependencies required. For instance, you will run `bundle install` for a ruby application

### `before_script`

Run any command necessary before running your build script.

### `script`

Run a script that effectively runs the build or the tests. For instance, in a ruby application, you would run `bundle exec rake`

### `after_success`

`after_success` commands are executed when the `script` commands are successful. A common task for after_success is to generate documentation, or to upload a build artifact to S3 for later use. You can also use this step to deploy your code to your staging or production servers.

### `after_failure`

`after_failure` commands are executed when any of the `script` commands failed. `after_failure` can be used in similar ways to `after_success`, for example to upload any log files that could help debugging a failure to S3.

### `after_script`

`after_script` commands are executed whatever the result of the previous commands was.

## Build Status: Success, Error, Failure

The build status of a bazooka build consists of 3 states:

### Success

### Error

When any of the steps in the `before_install`, `install` or `before_script` stages fails with a non-zero exit code, the build will be marked as errored.

### Failed

When any of the steps in the `script` stage fails with a non-zero exit code, the build will be marked as failed.
