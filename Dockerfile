FROM cfcommunity/bosh2-errand-resource:v0.1.2
ENV EXPECTED_SHA1 "7b89fb83064f40bd2e8adfc96c88f4b0b2bd6f34  check \
278fafe63fbf8af276d65640f4bafc94579cd742  in \
957d664beb1b976351c4db192abf1f6d08189fea  out \
"

RUN cd /opt/resource && export CURRENT_SHA=$(sha256sum *) && echo "Computed SHA256: " && printenv CURRENT_SHA \
    printenv CURRENT_SHA1 |sha256sum -c -

#FROM concourse/buildroot:base
#MAINTAINER https://github.com/cloudfoundry-community/bosh2-errand-resource
#
#ADD check /opt/resource/check
#ADD in /opt/resource/in
#ADD out /opt/resource/out
#
#RUN chmod +x /opt/resource/*



#FROM golang:1.21.6 AS builder
#ENV BUILT_BINARIES_DIR=/go/src/github.com/cloudfoundry-community/bosh2-errand-resource/build
#
#ADD . /go/src/github.com/cloudfoundry-community/bosh2-errand-resource
#WORKDIR /go/src/github.com/cloudfoundry-community/bosh2-errand-resource
#RUN ./ci/tasks/build.sh
#
#FROM alpine:3.19.1 AS resource
#COPY --from=builder /go/src/github.com/cloudfoundry-community/bosh2-errand-resource/build /opt/resource
#RUN chmod +x /opt/resource/*