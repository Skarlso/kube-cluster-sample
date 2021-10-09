#!/bin/zsh

create_test_environment() {
  echo -n "Beginning to deploy test environment to Kind cluster..."
  kind create cluster --config kind_test_cluster.yaml
  echo "done."

  echo -n "Creating the configmap for the database..."
  kubectl create configmap db-bootstrap --from-file=dbinit/database_setup.sql
  echo "done."

  echo -n "Creating secret..."
  kubectl apply -f kube_files/secret.yaml
  echo "done."

  echo -n "Creating mysql deployment..."
  kubectl apply -f kube_files/mysql.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=mysql
  echo "done."

  echo -n "Applying nsqlookup.yaml..."
  kubectl apply -f kube_files/nsqlookup.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=nsqlookup
  echo "done."

  echo -n "Applying nsqd.yaml..."
  kubectl apply -f kube_files/nsqd.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=nsqd
  echo "done."

  echo "Applying all services."

  echo -n "Applying receiver..."
  kubectl apply -f receiver/kube_files/receiver.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=receiver
  echo "done."

  echo -n "Applying face recognition..."
  kubectl apply -f kube_files/face_recognition.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=face-recog
  echo "done."

  echo -n "Applying image processor..."
  kubectl apply -f image_processor/kube_files/image_processor.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=image-processor
  echo "done."

  echo -n "Applying the frontend..."
  kubectl apply -f kube_files/frontend.yaml
  kubectl wait --for=condition=Ready --timeout 60s pod -l app=frontend
  echo "done."

  echo -n "Port-Forwarding receiver on 8000..."
  kubectl port-forward deployment/receiver-deployment 8000:8000 &
  echo "done."

  echo "All done. Port-Forward the receiver service to begin testing."
}

delete_test_cluster() {
  echo -n "Removing test cluster..."
  kind delete cluster --name kube-facerecog-test-cluster
  echo "done."
}

param=${1-create}

case "${param}" in
create)
    create_test_environment
    ;;
delete)
    delete_test_cluster
    ;;
*)
    echo "Usage: manage_test_environment {[create],delete}"
    ;;
esac
