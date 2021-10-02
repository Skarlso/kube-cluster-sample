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

<img style="position: fixed; bottom: 0px; right: 0px;" src="imgs/GO_BUILD.png" class="m-29 h-29" />

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

<img style="position: fixed; bottom: 0px; right: 0px;" src="imgs/Doctor_Who_Gopher_Woman.png" class="m-29 h-29" />

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

<img style="position: fixed; bottom: 0px; right: 0px;" src="imgs/GOPHER_VIKING.png" class="m-29 h-29" />