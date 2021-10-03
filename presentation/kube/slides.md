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
- <uil-screw /> **GRPC** - where does it come in?
- <uil-hourglass /> **Demo Time** - anything possible on a webpage
- <uil-rocket /> **Refactoring Exercise** - show me some code

<br>
<br>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Doctor_Who_Gopher_Woman.png" class="h-29"  alt="woman-who"/>
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
- Docker

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

# Technologies

<div class="absolute right-30px bottom-30px">
  <img src="imgs/STAR_TREK_GOPHER.png" class="h-29"  alt="star-trek"/>
</div>

---

## GRPC

- Why?

<v-click>

Because of how we interact with the Python FaceRecognition library.

</v-click>

- Use Case

<v-click>

Versatile control over the API and strict contracts on how to implement things. 

</v-click>

- buf.build

<v-click>

GRPC really shines once you implement and use buf.build.

</v-click>

- Benefits
- Drawbacks

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



---

## NSQ

> A realtime distributed messaging platform

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHERCON_ICELAND.png" class="h-29"  alt="iceland"/>
</div>

---

## Go

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Doctor_Who_Gopher.png" class="h-29"  alt="doctor-who"/>
</div>

---

## Face recognition

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

---

## Data flow

---

# Distributed System

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GoDZILLA.png" class="h-29"  alt="godzilla"/>
</div>
---

Definition of a distributed system:

> Multiple "something" ... linked together through the network to appear as one.

There has to be a sync point.

- How distributed is it?
- Where is the sync point?
- What level of consistency does it provide?
- Where are the pain points and possible resource contests?

---

# End
Thank you for listening!

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHER_MIC_DROP.png" class="h-29"  alt="mic"/>
</div>
