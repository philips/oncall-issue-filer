FROM gcr.io/google-appengine/golang

COPY . /go/src/app
RUN go-wrapper install
