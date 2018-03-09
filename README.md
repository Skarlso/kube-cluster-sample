# Kubernetes Sample Cluster Applications

A sample Micro-Service cluster for Kubernetes.

## Db

Database is MySQL which can run anywhere but needs to be accessible by all services.

## Receiver

An REST API which receives images to process. The image path is saved into a db and a message is placed in the queue for the image processor.

## Image Processor

The image processor picks up the image ID from the queue and gets the image path from the database.

Identifies the person in question and updates the image record in the db with an affiliated person id.

From here on anything can view the db and see them. Images which don't have a person yet are shown as `pending`, and images who have a person assigned, can link to a profile.

## Micro-Service Two

A simple service listening for work to do on a queue. It is possible that multiple instances pick up the same work to do.
This needs to be addressed.

## Messaging

`nsqlookupd`
`nsqd --lookupd-tcp-address=127.0.0.1:4160`
`nsqadmin --lookupd-http-address=127.0.0.1:4161`

```bash
curl -d '{"path":"/home/user/image_dump"}' http://127.0.0.1:8000/image/post
got path: {Path:/home/user/image_dump}
image saved with id: 23
image sent to nsq
```

## Deploying

Using Kubernetes to deploy sample application into various clusters.

## How it looks like on the front-end

| ID | Path                         | Person     |
|----|------------------------------|------------|
| 51 | /home/user/image_dump/test33 | Pending... |
| 52 | /home/user/image_dump/test34 | Hannibal   |
| 53 | /home/user/image_dump/test35 | John Doe   |
| 54 | /home/user/image_dump/test36 | Gergely    |
