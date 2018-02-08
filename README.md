# Kubernetes Sample Cluster Applications

A sample Micro-Service cluster for Kubernetes.

## Micro-Service One

A simple service sending something to a message queue.

### Distributed

Needs a MongoDB to share information between instances.

## Micro-Service Two

A simple service listening for work to do on a queue. It is possible that multiple instances pick up the same work to do.
This needs to be addressed.

## Messaging

`nsqlookupd`.

## Deploying

Using Kubernetes to deploy sample application into various clusters.
