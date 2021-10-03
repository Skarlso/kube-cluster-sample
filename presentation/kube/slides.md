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
  <img src="imgs/GO_BUILD.png" class="h-29" />
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
- <uil-rocket /> **Refactoring Exercise** - shuddup show me some code

<br>
<br>

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Doctor_Who_Gopher_Woman.png" class="h-29" />
</div>

<style>
h1 {
  background-color: #2B90B6;
  background-image: linear-gradient(45deg, #4EC5D4 10%, #146b8c 20%);
  background-size: 100%;
  -webkit-background-clip: text;
  -moz-background-clip: text;
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
  <img src="imgs/GOPHER_VIKING.png" class="h-29" />
</div>

---

# Introduction

- Name: Gergely Brautigam
- Work: Weaveworks
- Twitter: https://twitter.com/Skarlso
- Github: https://github.com/Skarlso
- Website: https://gergelybrautigam.com

<div class="absolute right-30px bottom-30px">
  <img src="imgs/DOCTOR_STRANGE_GOPHER.png" class="h-29" />
</div>

---

# Technologies

<div class="absolute right-30px bottom-30px">
  <img src="imgs/STAR_TREK_GOPHER.png" class="h-29" />
</div>

---

## GRPC

- Why?
- Use Case
- buf.build
- Benefits
- Drawbacks

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GopherLink.png" class="h-29" />
</div>

---

## Kubernetes

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Kubernetes_Gophers.png" class="h-29" />
</div>

---

## NSQ

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHERCON_ICELAND.png" class="h-29" />
</div>

---

## Go

<div class="absolute right-30px bottom-30px">
  <img src="imgs/Doctor_Who_Gopher.png" class="h-29" />
</div>

---

## Face recognition

<div class="absolute right-30px bottom-30px">
  <img src="imgs/BATMAN_GOPHER.png" class="h-29" />
</div>


---

# Architecture

<img style="position: fixed; top: 0px;" src="imgs/kube_architecture.png" class="m-29 h-100" />

<div class="absolute right-30px bottom-30px">
  <img src="imgs/RickAndMorty.png" class="h-29" />
</div>

---

# Distributed System

- Why, At what level, Eventual Consistency model, where is the sync

Hint: Database and the file storage on the image.

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GoDZILLA.png" class="h-29" />
</div>

---

# End
Thank you for listening!

<div class="absolute right-30px bottom-30px">
  <img src="imgs/GOPHER_MIC_DROP.png" class="h-29" />
</div>
