from concurrent import futures
import time
import face_recognition
import sys
import os
import re
import click
from os.path import basename
import grpc
import face_pb2
import face_pb2_grpc

_ONE_DAY_IN_SECONDS = 60 * 60 * 24

class Identifer(face_pb2_grpc.IdentifyServicer):

    def Identify(self, request, context):
        path = request.image_path
        image_name = self.identify(path)
        return face_pb2.IdentifyResponse(image_name = image_name)


    def image_files_in_folder(self, folder):
        return [os.path.join(folder, f) for f in os.listdir(folder) if re.match(r'.*\.(jpg|jpeg|png)', f, flags=re.I)]


    def identify(self, path_to_unknown):
        if len(path_to_unknown) < 1:
            return "none"
        print("Checking image: %s" % path_to_unknown)
        known_people = os.getenv('KNOWN_PEOPLE', 'known_people')
        print("Known people images location is: %s" % known_people)
        images = self.image_files_in_folder(known_people)
        unknown_image = face_recognition.load_image_file(path_to_unknown)
        unknown_encoding = face_recognition.face_encodings(unknown_image)[0]
        for image in images:
            known_image = face_recognition.load_image_file(image)
            known_encoding = face_recognition.face_encodings(known_image)[0]
            results = face_recognition.compare_faces([known_encoding], unknown_encoding)
            if results:
                return basename(image)
        return "none"


class HealthChecker(face_pb2_grpc.HealthCheckServicer):
    def HealthCheck(self, request, context):
        return face_pb2.HealthCheckResponse(ready=True)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    face_pb2_grpc.add_IdentifyServicer_to_server(Identifer(), server)
    face_pb2_grpc.add_HealthCheckServicer_to_server(HealthChecker(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()
