FROM gcr.io/google-appengine/golang

COPY . /go/src/github.com/philips/oncall-issue-filer
COPY . /go/src/app
RUN go-wrapper install
