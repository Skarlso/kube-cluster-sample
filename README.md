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

Face recognition now works, thanks to: [Face Recognition](https://github.com/ageitgey/face_recognition).

Running the python script results in:

```bash
❯ python3 identifier.py unkown.jpg
Checking image: unkown.jpg
matched id: hannibal_1.jpg
```

A person can have many images associated with him/her. The images are saved for a person in the db and linked to it.
Once the Face Recognition service identifies an image it will send back the name of the image. A joined query produces the name of the person which the front-end displays.

## Face Recognition

Using protobuf and gRPC, the face recognition service talks to the image processor service. This ensures the flexibility of exchanging the face recognition service to whatever implementation is available at any given point in time.

## Circuit Breaker

```bash
2018/03/14 21:03:01 timeout over. opening circuit.
2018/03/14 21:03:01 could not send image: rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp [::1]:50051: connect: connection refused"
2018/03/14 21:03:03 Processing image id:  9
2018/03/14 21:03:03 circuit breaker try count:  1
2018/03/14 21:03:03 could not send image: rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp [::1]:50051: connect: connection refused"
2018/03/14 21:03:05 Processing image id:  10
2018/03/14 21:03:05 circuit breaker try count:  2
2018/03/14 21:03:05 could not send image: rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp [::1]:50051: connect: connection refused"
2018/03/14 21:03:05 maximum try of 3 sends reached. disabling for 10s time period.
2018/03/14 21:03:07 Processing image id:  11
2018/03/14 21:03:07 circuit breaker try count:  3
2018/03/14 21:03:07 max sending try count of 3 reached. sending not allowed for 7.998267708s time period.
2018/03/14 21:03:09 Processing image id:  12
2018/03/14 21:03:09 circuit breaker try count:  3
2018/03/14 21:03:09 max sending try count of 3 reached. sending not allowed for 5.995015753s time period.
2018/03/14 21:03:11 Processing image id:  13
2018/03/14 21:03:11 circuit breaker try count:  3
2018/03/14 21:03:11 max sending try count of 3 reached. sending not allowed for 3.994391825s time period.
```

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
