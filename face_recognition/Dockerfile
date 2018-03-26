FROM python:3.6-slim

LABEL Author="Gergely Brautigam"

RUN apt-get update && apt-get install -y --fix-missing \
    build-essential \
    cmake \
    gfortran \
    git \
    wget \
    curl \
    graphicsmagick \
    libgraphicsmagick1-dev \
    libatlas-dev \
    libavcodec-dev \
    libavformat-dev \
    libgtk2.0-dev \
    libjpeg-dev \
    liblapack-dev \
    libswscale-dev \
    pkg-config \
    python3-dev \
    python3-numpy \
    python3-pip \
    software-properties-common \
    zip \
    && apt-get clean && rm -rf /tmp/* /var/tmp/*

RUN cd ~ && \
    mkdir -p dlib && \
    git clone -b 'v19.9' --single-branch https://github.com/davisking/dlib.git dlib/ && \
    cd  dlib/ && \
    python3 setup.py install --yes USE_AVX_INSTRUCTIONS

RUN python3 -m pip install face_recognition
RUN python3 -m pip install grpcio
RUN python3 -m pip install grpcio-tools

COPY face_pb2_grpc.py /root
COPY face_pb2.py /root
COPY identifier.py /root

VOLUME [ "/unknown_people", "/known_people" ]

EXPOSE 50051

CMD [ "python3", "/root/identifier.py" ]
