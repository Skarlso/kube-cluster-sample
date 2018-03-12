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
        person_id = self.identify(path)
        return face_pb2.IdentifyResponse(person_id=person_id)


    def image_files_in_folder(self, folder):
        return [os.path.join(folder, f) for f in os.listdir(folder) if re.match(r'.*\.(jpg|jpeg|png)', f, flags=re.I)]


    def identify(self, path_to_unknown):
        print("Checking image: %s" % path_to_unknown)
        # gather all known people
        # the name of the file is the id of the person
        # compare to all and return if match is found
        images = self.image_files_in_folder('known_people')
        unknown_image = face_recognition.load_image_file(path_to_unknown)
        unknown_encoding = face_recognition.face_encodings(unknown_image)[0]
        for image in images:
            known_image = face_recognition.load_image_file(image)
            known_encoding = face_recognition.face_encodings(known_image)[0]
            results = face_recognition.compare_faces([known_encoding], unknown_encoding)
            if results:
                base = basename(image)
                id = os.path.splitext(base)[0]
                return int(id)
        return "-1"


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    face_pb2_grpc.add_IdentifyServicer_to_server(Identifer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()
