from __future__ import print_function

import grpc

import face_pb2
import face_pb2_grpc


def run():
    # Test Identifier
    channel = grpc.insecure_channel('localhost:50051')
    stub = face_pb2_grpc.IdentifyStub(channel)
    response = stub.Identify(face_pb2.IdentifyRequest(image_path='unknown.jpg'))
    print("Identifier client received: %s" % response.image_name)

    # Test HealthCheck
    health = face_pb2_grpc.HealthCheckStub(channel)
    response = health.HealthCheck(face_pb2.Empty())
    print("Response from health received %r" % response.ready)

if __name__ == '__main__':
    run()