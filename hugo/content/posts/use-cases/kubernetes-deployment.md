---
title: "Kubernetes Deployment"
date: 2018-10-20T19:46:18+13:00
draft: true
---

Kubernetes resources can be described declaratively, by yaml files. A typical microservice might have 3 or 4 resources. For example, a Deployment, a Service, a ConfigMap and an Ingress.

When you have 20+ microservices (or more), then you can end up with a lot of lengthy, similar yaml files.

In the past we've considered 'templating' yaml files. You could use one of many commandline templating tools out there, and we have tried. There's one problem - templated yaml is no longer yaml. It generates yaml, but your IDE and static analysis tools can no longer help you.

Let's take a file with 3 resources …


