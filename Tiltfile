# Deploy: tell Tilt what YAML to deploy
k8s_yaml('db-bootstrap',
    './kube_files/secret.yaml',
    './kube_files/mysql.yaml',
    './kube_files/nsqlookup.yaml',
    './kube_files/nsqd.yaml',
    './kube_files/face_recognition.yaml',
    './kube_files/frontend.yaml',
    './receiver/kube_files/receiver.yaml',
    './image_processor/kube_files/image_processor.yaml')

# Build: tell Tilt what images to build from which directories
docker_build('skarlso/image-processor', 'image_processor')
docker_build('skarlso/face-recognition', 'face_recognition')
docker_build('skarlso/frontend', 'frontend')
docker_build('skarlso/receiver', 'receiver')
# ...

# Watch: tell Tilt how to connect locally (optional)
k8s_resource('receiver', port_forwards=8000)
k8s_resource('frontend', port_forwards=8081)
