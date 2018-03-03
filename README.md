# Kubernetes Sample Cluster Applications

A sample Micro-Service cluster for Kubernetes.

## Receiver

An REST API which receives images to process. The image path is saved into a db and a message is placed in the queue for the image processor.

## Image Processor

The image processor picks up the image ID from the queue and gets the image path from the database.

Identifies the person in question and updates the image record in the db with an affiliated person id.

From here on anything can view the db and see them. Images which don't have a person yet are shown as `pending`, and images who have a person assigned, can link to a profile.

### Distributed

[BoltDB](https://github.com/coreos/bbolt) will be used to share image information.

## Micro-Service Two

A simple service listening for work to do on a queue. It is possible that multiple instances pick up the same work to do.
This needs to be addressed.

## Messaging

`nsqlookupd`.

## Deploying

Using Kubernetes to deploy sample application into various clusters.
