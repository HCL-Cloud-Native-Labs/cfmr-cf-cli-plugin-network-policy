# GitHub Workflows

This directory contains the definitions for [GitHub Actions workflows].

[GitHub Actions workflows]: https://docs.github.com/en/free-pro-team@latest/actions/reference/workflow-syntax-for-github-actions

<!-- omit in toc -->
## Table of Contents

* [Building & Pushing](#building-and-pushing)


## Building and Pushing

The pull request & continuous integration workflow is:

* On pull requests, [`main.yml`] does the following:
  * Pulls the code branch using [`actions/checkout`]
  * Login to the IBM docker entitled registry using [`docker/login-action`]
  * Builds the eirinix annotate docker image.
  * Push the image to staging IBM entitled registry using docker [`docker/build-push-action`]
* On pull requests approvals, same thing as above except done on the merged code with main branch

[`main.yml`]: main.yml
[`actions/checkout`]: https://github.com/actions/checkout
[`docker/login-action`]: https://github.com/docker/login-action
[`docker/build-push-action`]: https://github.com/docker/build-push-action

