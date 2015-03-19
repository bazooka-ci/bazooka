# Building a Java project

This guide walks you through the build environment and configuration topics specific to Java projects. Please make sure to read our [Getting Started](../home/getting_started) and [general build configuration](../home/build_configuration) guides first.

## Choose the Java language

First of all, you need to specify in your build configuration that you project main language will be Java

```yaml
language: java
```

## Java versions

Bazooka built-in Java support comes with the following jvm implementations

* openjdk6
* openjdk7
* openjdk8
* oraclejdk6
* oraclejdk8

You can choose multiple versions on which your project will be tested

```yaml
language: java
jdk:
  - openjdk8
  - oraclejdk7
  - oraclejdk8
```

This will generate [permutations](../home/permutations) on which your project will be tested

## Discovering your build tool

Bazooka will try to find which build tool you are using in your project and define reasonable defaults for your build phases. Of course, these values can be overridden if you don't want to use them

To determine which build tool your project is using :

* If a file named `build.gradle` exists in your project:
    * If a file named `gradlew` exists in your project, the build tool is **gradlew**
    * Otherwise, the build tool to **gradle**
* If a file name *pom.xml* exists in your project, the build tool is **maven**
* Otherwise, the build tool is **ant** by default

## Project using Maven

the defaults generated for a *maven* project are equivalent to the following configuration

```yaml
install: mvn install -DskipTests=true
script: mvn test
```

## Project using Gradlew

the defaults generated for a *gradlew* project are equivalent to the following configuration

```yaml
install: ./gradlew assemble
script: ./gradlew check
```

## Project using Gradle

the defaults generated for a *gradle* project are equivalent to the following configuration

```yaml
install: gradle assemble
script: gradle check
```
