#!/usr/bin/env python3

import face_recognition
import sys
import os
import re
import click
from os.path import basename

def image_files_in_folder(folder):
    return [os.path.join(folder, f) for f in os.listdir(folder) if re.match(r'.*\.(jpg|jpeg|png)', f, flags=re.I)]


@click.command()
@click.argument('path_to_unknown')
def main(path_to_unknown):
    print("Checking image: %s" % path_to_unknown)
    # gather all known people
    # the name of the file is the id of the person
    # compare to all and return if match is found
    images = image_files_in_folder('known_people')
    unknown_image = face_recognition.load_image_file(path_to_unknown)
    unknown_encoding = face_recognition.face_encodings(unknown_image)[0]
    for image in images:
        known_image = face_recognition.load_image_file(image)
        known_encoding = face_recognition.face_encodings(known_image)[0]
        results = face_recognition.compare_faces([known_encoding], unknown_encoding)
        if results:
            print("matched id: %s" % basename(image))
            break
    print("-1")

if __name__ == "__main__":
    main()
