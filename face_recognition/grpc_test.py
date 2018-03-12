from __future__ import print_function

import grpc

import face_pb2
import face_pb2_grpc


def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = face_pb2_grpc.IdentifyStub(channel)
    response = stub.Identify(face_pb2.IdentifyRequest(image_path='unknown.jpg'))
    print("Identifier client received: %d" % response.person_id)


if __name__ == '__main__':
    run()