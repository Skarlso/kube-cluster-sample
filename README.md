# Kubernetes Sample Cluster Applications

A sample Micro-Service application for my blog post about deploying with Kubernetes located here: [Deploying a Distributed Application with Kubernetes](https://skarlso.github.io/2018/03/15/kubernetes-distributed-application/).

## Db

Database is MySQL which can run anywhere but needs to be accessible by all services. In order to debug the database,
first, determine what pod it runs on with:

```bash
kubectl get pods
```

Then run:

```bash
kubectl port-forward mysql-77976765f9-j25xl 3306:3306
```

... to expose the database to localhost. Once that is done, simply access it with a local mysql client:

```bash
mysql -ptcp -h127.0.0.1 -P3306 -uroot -ppassword123
```

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

Images have a `status` field. The field simply takes care of tracking the state of an image which can be these three ATM:

```go
const (
    // PENDING -- not yet send to face recognition service
    PENDING Status = iota
    // PROCESSED -- processed by face recognition service; even if no person was found for the image
    PROCESSED
    // FAILEDPROCESSING -- for whatever reason the processing failed and this image is flagged for a retry
    FAILEDPROCESSING
)
```

If an image is failed processing those can be later reconciled / retried, by a cron job for example, scavenging for images which failed processing.

To generate the protobuf code, run:

```bash
protoc -I facerecog/ facerecog/face.proto --go_out=plugins=grpc:facerecog
```

## Face Recognition

Using protobuf and gRPC, the face recognition service talks to the image processor service. This ensures the flexibility of exchanging the face recognition service to whatever implementation is available at any given point in time.

To generate the protobuf code, run:

```bash
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. face.proto
```

## Circuit Breaker

The circuit breaker is a rudimentary circuit health check device. If it receives a pre-configured number of failed attempts at calling the back-end face recognition service, it will break the circuit for a period of time. After that period elapses it will re-try with a `Ping` to see if the service is at health. If the ping comes back as green, it opens the flow again.

Ping:

```bash
2018/03/15 11:04:36 timeout over. running ping.
2018/03/15 11:04:36 backend still not functioning. extending break.
2018/03/15 11:04:36 max sending try count of 3 reached. sending not allowed for 9.999999735s time period.
```

Normal run:

```bash
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

| ID | Path                         | Person     | Status |
|----|------------------------------|------------|--------|
| 51 | /home/user/image_dump/test33 | Pending... | 0      |
| 52 | /home/user/image_dump/test34 | Hannibal   | 1      |
| 53 | /home/user/image_dump/test35 | John Doe   | 1      |
| 54 | /home/user/image_dump/test36 | Gergely    | 1      |
| 55 | /home/user/image_dump/test37 | Pending... | 2      |

## Mysql

Create the config map for the database which bootstraps the db.

```bash
kubectl create configmap db-bootstrap --from-file=dbinit/database_setup.sql
```

## Troubleshooting

If a container is stuck on creating `kubectl describe pods` lists all the last actions.

For example:

```
Events:
  Type    Reason                 Age   From               Message
  ----    ------                 ----  ----               -------
  Normal  Scheduled              2m    default-scheduler  Successfully assigned nsqlookup-c9dc7574c-dtnjv to minikube
  Normal  SuccessfulMountVolume  2m    kubelet, minikube  MountVolume.SetUp succeeded for volume "default-token-bvtbm"
  Normal  Pulling                1m    kubelet, minikube  pulling image "nsqio/nsq"
```

We can see that this pod is pulling an image for a container, which can be large, so it takes a while.

# Testing with Kind

There is a project under sigs called [Kind](https://github.com/kubernetes-sigs/kind). I'm using it for testing.

To test this deployment, you'll need a storage resource.

Apply label to the nodes in order for the PVC to work.

```
kubectl label nodes <your-node-name> local-pvc=true
```

Create a PV which will describe a volume resource on the cluster. Then create a claim which will claim it for the service.

Both can be found under kube_files called `face_recognition_pv{c}_{un}known.yaml`. 

# Slides

Slides is using [slidev](https://sli.dev).

```shell
cd presentation/kube
npm i
npm run dev
```

This should serve the `slides.md` located in that folder.
