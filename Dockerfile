FROM concourse/buildroot:base
MAINTAINER https://github.com/cloudfoundry-community/bosh2-errand-resource

ADD check /opt/resource/check
ADD in /opt/resource/in
ADD out /opt/resource/out

RUN chmod +x /opt/resource/*
