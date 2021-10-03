#!/bin/zsh

echo -n "Beginning to deploy test environment to Kind cluster..."
kind create cluster --config kube_files/kind_test_cluster.yaml
echo "done."

echo -n "Creating the configmap for the database..."
kubectl create configmap db-bootstrap --from-file=database_setup.sql --wait
echo "done."

echo -n "Creating secret..."
kubectl apply -f kube_files/secret.yaml
echo "done."

echo -n "Creating mysql deployment..."
kubectl apply -f kube_files/mysql.yaml
kubectl wait --for=condition=Ready pod -l app=mysql
echo "done."

echo -n "Applying nsqlookup.yaml..."
kubectl apply -f kube_files/nsqlookup.yaml
kubectl wait --for=condition=Ready pod -l app=nsqlookup
echo "done."

echo "Applying all services."

echo -n "Applying receiver..."
kubectl apply -f receiver/kube_files/receiver.yaml
kubectl wait --for=condition=Ready pod -l app=receiver
echo "done."

echo -n "Applying face recognition and labelling workers..."
kubectl label nodes kind-worker local-pvc=true
kubectl label nodes kind-worker2 local-pvc=true
kubectl apply -f kube_files/face_recognition_pv_known.yaml
kubectl apply -f kube_files/face_recognition_pv_unknown.yaml
kubectl apply -f kube_files/face_recognition_pvc_known.yaml
kubectl apply -f kube_files/face_recognition_pvc_unknown.yaml
kubectl apply -f kube_files/face_recognition_pv_known.yaml
kubectl apply -f kube_files/face_recognition.yaml
kubectl wait --for=condition=Ready pod -l app=face-recog
echo "done."

echo -n "Applying image processor..."
kubectl apply -f image_processor/kube_files/image_processor.yaml
kubectl wait --for=condition=Ready pod -l app=image-processor
echo "done."

echo -n "Applying the frontend..."
kubectl apply -f kube_files/frontend.yaml
kubectl wait --for=condition=Ready pod -l app=frontend
echo "done."

echo "All done. Port-Forward the receiver service to begin testing."