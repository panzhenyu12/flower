FROM ubuntu:trusty

LABEL maintainer="zhenyupan" \
      developer="pczy"

RUN mkdir /thor/

ADD thor /thor/
ADD thor.json /thor/
ADD templates/ /thor/templates/
ADD deepface.sql /thor/

WORKDIR /thor/

RUN ls -lh
CMD ./thor
