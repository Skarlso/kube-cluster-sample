---
theme: default
layout: cover
# apply any windi css classes to the current slide
class: 'text-center'
# https://sli.dev/custom/highlighters.html
highlighter: shiki
# show line numbers in code blocks
lineNumbers: false
# some information about the slides, markdown enabled
info: |
  Distributed FaceRecognition using Kubernetes and Go
# persist drawings in exports and build
drawings:
  persist: false
---

# Welcome to Distributed Face-recognition with Go

Name: Gergely Brautigam <br>
Github: https://github.com/Skarlso/kube-cluster-sample <br>
Twitter: https://twitter.com/Skarlso <br>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GO_BUILD.png" class="h-29"  alt="build"/>
</div>

---

# Agenda

- <uil-calender /> **Gauge the Audience** - who knows what how much where why
- <uil-user /> **Small Introduction** - theme can be shared and used with npm packages
- <uil-laptop /> **Technologies used** - kubernetes, face-recognition library, Go, nsq, GRPC
- <uil-jackhammer /> **Architecture** - a bit about the architecture, overview
- <uil-ship /> **Kubernetes** - deeper into the rabbit hole with Kubernetes
- <uil-fidget-spinner /> **Distributed Systems** - why is it distributed exactly?
- <uil-spinner /> **User Cases** - what is all of this even useful for?
- <uil-hourglass /> **Demo Time** - anything possible on a webpage
- <uil-rocket /> **Refactoring Exercise** - show me some code

<br>
<br>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Doctor_Who_Gopher_Woman.png" class="h-29" alt="woman-who"/>
</div>

<style>
h1 {
  background-color: #2B90B6;
  background-image: linear-gradient(45deg, #4EC5D4 10%, #146b8c 20%);
  background-size: 100%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  -moz-text-fill-color: transparent;
}
</style>

---

# Gauge the Audience
Who has heard off, knows, is familiar, understands

- Kubernetes
- GRPC
- Distributed Systems
- NSQ

<img style="position: fixed; top: 0;" src="imgs/the_expert.jpeg" class="m-69 h-50"  alt="the-expert"/>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHER_VIKING.png" class="h-29"  alt="viking"/>
</div>

---

# Introduction

- Name: Gergely Brautigam
- Work: Weaveworks
- Twitter: https://twitter.com/Skarlso
- Github: https://github.com/Skarlso
- Website: https://gergelybrautigam.com

<div class="absolute right-30px bottom-30px">
  <img src="imgs/DOCTOR_STRANGE_GOPHER.png" class="h-29"  alt="doctor-strange"/>
</div>

---

# Why?

<v-click>

For the glory of Sontaran empire.

<img style="position: fixed; top: 0;" src="imgs/sontaran.jpeg" class="m-40 h-85"  alt="the-expert"/>

</v-click>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Doctor_Who_Gopher.png" class="h-29"  alt="doctor-who"/>
</div>

---

# Technologies

<div class="absolute right-30px bottom-30px">
  <img src="imgs/STAR_TREK_GOPHER.png" class="h-29"  alt="star-trek"/>
</div>

---

## GRPC

<v-clicks>

- Why?
> Strict contract, easy to use, type safety, discoverability ( more visible than a JSON API which doesn't have a SWAGGER doc ).

- Use Case
> Versatile control over the API and strict contracts on how to implement things.

- buf.build
> protobuf really shines once you implement and use buf.build with re-usable protobuf modules.

- Drawbacks
> Rigid structure, having to provide an SDK, more overhead to implement from the client side ( for dynamic languages having to
> use GRPC generated code instead of putting the JSON payload into a dict or a map ).

</v-clicks>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GopherLink.png" class="h-29"  alt="link"/>
</div>

---

## Kubernetes

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Kubernetes_Gophers.png" class="h-29"  alt="kube"/>
</div>

---

### Let's try to limit this

- Deployments vs StatefulSets
- Network policies (restrict who can talk to what)
- Service discovery
- Secrets and ConfigMaps
- Persistent Volumes and Claims
- Resource Limits
- LoadBalancing, certificate manager and LetsEncrypt (none of that is set up but sounds cool)
- Readiness and Liveness probes

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher1.png" class="h-29"  alt="gopher1"/>
</div>

---

## NSQ

> A realtime distributed messaging processing at scale.

- Versatile
- Easy to use and set up
- Works out of the box
- There are a number of alternatives since then like KubeMQ but NSQ remains strong
- Distributed setup using nsqd and nsqlookupd is amazing and uses little resources

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHERCON_ICELAND.png" class="h-29"  alt="iceland"/>
</div>

---

## Go

The main question... How is Go helping in all of this?

<v-clicks>

- Fast
- Easy concurrency
- Easy to write and begin ( started this project more than 4 years ago )
- Lots of nice libraries and wonderful online documentation partially enforced by the linter ( therefore community standard )
- Circuit breaker ( non-retry, simple blocker )

</v-clicks>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GIRL_GOPHER.png" class="h-29"  alt="doctor-who"/>
</div>

---

## Face recognition

- [face-recognition](https://github.com/ageitgey/face_recognition) library in Python (SORRY)

> it's processing images in pools of 10

- Why Python?

<v-click>

> The gocv library wasn't as developed back when I wrote this as it is now.
> GoCV doesn't provide face recognition. It provides face detection.
> Multiple capabilities of the face recognition library such as:
>   - Realtime face recognition of multiple faces in photos
>   - Recognize faces on a raspberry pi camera
>   - And many more...[link](https://github.com/ageitgey/face_recognition#facial-recognition)

</v-click>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/BATMAN_GOPHER.png" class="h-29"  alt="batman"/>
</div>


---

# Architecture

<div class="absolute right-30px bottom-30px">
  <img src="imgs/RickAndMorty.png" class="h-29"  alt="rick-and-morty"/>
</div>

---

## Microservice Architecture

<img style="position: fixed; top: 0;" src="imgs/kube_architecture.png" class="m-29 h-100"  alt="architecture"/>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher2.png" class="h-29"  alt="gopher2"/>
</div>

---

## Data flow

<v-clicks>

<img style="position: fixed; top: 0;" src="imgs/flow-00.png" class="m-29 h-100"  alt="flow00"/>
<img style="position: fixed; top: 0;" src="imgs/flow-01.png" class="m-29 h-100"  alt="flow01"/>
<img style="position: fixed; top: 0;" src="imgs/flow-02.png" class="m-29 h-100"  alt="flow02"/>
<img style="position: fixed; top: 0;" src="imgs/flow-03.png" class="m-29 h-100"  alt="flow03"/>
<img style="position: fixed; top: 0;" src="imgs/flow-04.png" class="m-29 h-100"  alt="flow04"/>
<img style="position: fixed; top: 0;" src="imgs/flow-05.png" class="m-29 h-100"  alt="flow05"/>
<img style="position: fixed; top: 0;" src="imgs/flow-06.png" class="m-29 h-100"  alt="flow06"/>

</v-clicks>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher3.png" class="h-29"  alt="gopher3"/>
</div>

---

# Distributed System

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GoDZILLA.png" class="h-29"  alt="godzilla"/>
</div>

---

Definition of a distributed system:

> Multiple "something" ... linked together through the network to appear as one.

There has to be a sync point unless you don't have shared state.

- How distributed is it?
- Where is the sync point?
- Transactions...
- What level of consistency does it provide?
- Where are the pain points and possible resource contests?
- Concurrent state management
- Sync methods, quorum, raft, etc...
- Sharding with PlanetScale, AWS Aurora, GCP Cloud Database

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher4.png" class="h-29"  alt="gopher4"/>
</div>

---

# Use cases

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher5.png" class="h-29"  alt="gopher5"/>
</div>

---

## Meeting announcer



<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher6.png" class="h-29"  alt="gopher6"/>
</div>

---

## Dog belly rubbing appointment

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher7.png" class="h-29"  alt="gopher7"/>
</div>


---

## Cat feeding time

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher8.png" class="h-29"  alt="gopher8"/>
</div>


---

## Mood tracker

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher9.png" class="h-29"  alt="gopher9"/>
</div>

---

## Super Simple Architecture

<img style="position: fixed; top: 0;" src="imgs/simple-design.png" class="m-29 h-100"  alt="flow00"/>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher11.png" class="h-29"  alt="gopher9"/>
</div>

---

## Real purpose

<v-clicks>

- show off technoloy
- show off the power of Go
- learn architecture
- learn Kubernetes
- learn more in-depth Go

</v-clicks>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher12.png" class="h-29"  alt="gopher9"/>
</div>

---

# Demo Time

Let's get crackin'!

<div class="absolute right-30px bottom-30px">
  <img src="imgs/SPACEGIRL1.png" class="h-29"  alt="spacegirl"/>
</div>

---

# Refactoring

- What
- Why
- Where
- How

<div class="absolute right-30px bottom-30px">
  <img src="imgs/gopher10.png" class="h-29"  alt="gopher10"/>
</div>

---

# End
Thank you for listening!

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHER_MIC_DROP.png" class="h-29"  alt="mic"/>
</div>
