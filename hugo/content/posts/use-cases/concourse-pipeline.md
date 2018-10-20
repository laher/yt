---
title: "Concourse Pipeline"
date: 2018-10-20T19:46:25+13:00
draft: true
---

[Concourse CI](https://concourse-ci.org/) pipelines are typically described in YAML. It seems to be a common characteristic for CI tools at the moment. _For other CI platforms using yaml configurations, see Circle CI, Travis, droid.io, Gitlab CI, …_

For an introduction to Concourse:

> _If you're just getting started with Concourse, we recommend that you start with the [Concourse Tutorials](https://concoursetutorial.com/) developed and maintained by our friends at Stark & Wayne._

## The Use Case

 * Team X wants to set up pipelines for 20+ services, each pretty similar but with small differences.
 * The pipelines
   * Each service will have a 'main' pipeline:
     * Track master branch
     * Run unit tests & build
     * Deploy to 'staging'
     * Run integrations tests
     * Deploy to production
     * Run smoke tests
   * Each service will need additional pipelines for 'feature testing' environments:
     * Track a branch
     * Run unit tests & build
     * Deploy to feature-testing environment
     * Run integration tests

So, each pipeline will have some repetition, and we'll need a lot of them. Concourse pipelines can get quite big, and maintainability becomes a concern.

Ideally we'd like to:
 * Break the pipeline into building blocks.
 * Concisely define the differences between the pipelines.

## The anatomy of a pipeline

At the top level, our pipelines have three properties:

```
resource_types: {}
resources: {}
jobs: {}
```

 * Resource types and resources are usually very similar for each service - git, S3, Slack, Docker registries, …
 * Some jobs are very similar for each pipeline (flick-version, unit-test, build, docker-push).
 * Some jobs are repeated with small differences. (deploy-dev vs deploy-prod. integration-test-dev vs smoke-test-prod)

## Breaking the files down

…
