# echopilot

A repo for building a containerpilot based dumb echo server in order to test out library insertion and external service integrations

## Prerequisites

- Install docker: <https://docs.docker.com/install/>
- Get devbox: https://www.jetpack.io/devbox
- Generate a self-signed cert if you do not have one (useful for testing) (see below)

## Building and running locally
```bash
# The intention is to make a number of these step automatic in the future, but for now you need to be aware that
# these might be necessary.

# enter the devbox shell for ease. This could take a while the first time, but will be quick afterwards
devbox shell

# Generate some self-signed certs for yourself if they are not already present.
openssl req -x509 -newkey rsa:4096 -nodes -subj "/C=US/ST=California/L=Who knows/O=Shady Interprises/OU=self signers/CN=www.example.com" -sha256 -keyout configs/ssl/key.pem -out configs/ssl/cert.pem -days 365

# Build the protobuf (only needed if you changed them really)
buf generate --config proto/buf.yaml --template proto/buf.gen.yaml proto/

# Rebuild the templates (only needed if you changed them)
templ generate

# download all you dependencies
go get -v ./...

# build the binary
go build

# run the binary. I have a config to make your life easier.
./echopilot serve --config etc/echopilot.json

# Now you can navigate to https://127.0.0.1:1443/ and test it out!

# Probably exit your devbox shell before you do something else and forget...
exit
```

## Building and running the docker container
TODO: Fix this up. I don't think it will work as-is right now, but it's close.

```bash
# Generate some self-signed certs for yourself
openssl req -x509 -newkey rsa:4096 -nodes -subj "/C=US/ST=California/L=Who knows/O=Shady Interprises/OU=self signers/CN=www.example.com" -sha256 -keyout configs/ssl/key.pem -out configs/ssl/cert.pem -days 365

# Build the container image
docker build -t echopilot .

# Run the container while mounting the certificates
docker run --rm -p 3000:3000 --mount type=bind,src="$(pwd)/configs/echopilot.json",target="/etc/echopilot/echopilot.json" --mount type=bind,src="$(pwd)/configs/ssl/cert.pem",target="/etc/echopilot/cert.pem" --mount type=bind,src="$(pwd)/configs/ssl/key.pem",target="/etc/echopilot/key.pem" echopilot

# Curl an endpoint
curl https://localhost:3000/echo.v1.EchoService/EchoString -k --data '{"content": "testing"}' --header "Content-Type: application/json"

# Checkout the basic UI in your browser at https://localhost:3000/ui
# A Ctrl-C in the terminal running the container will halt the server gracefully.
```

docker run

## TODO

Previously, this repo used vagrant to spin up a VM with a number of associated
services designed to collect metrics and logs. Since VMs on mac M1 silicon is...
difficult... I'm going to instead create a docker compose file to orchestrate an
entire environment.

## Testing

I broke all the tests and haven't fixed them yet! Don't expect them to work right now!

Just a classic `go test` will do it. Note that you must install golang, buf, etc.
for this to work since it has to build locally.

## Building

### Generating a self-signed certificate

You should really NOT do this for anything going to production or being exposed to
the internet, but for testing it is really nice to just have a simple self-signed
cert. They're easy to make too!

For this, we will be creating the cert and key in the configs/ssh/ directory.
That path is in the .gitignore file, so the certs we create will not end up
on github.

```bash
openssl req -x509 -newkey rsa:4096 -nodes -subj "/C=US/ST=California/L=Who knows/O=Shady Interprises/OU=self signers/CN=www.example.com" -sha256 -keyout configs/ssl/key.pem -out configs/ssl/cert.pem -days 365
```

### Protoc build (outdated, see below)

In the past, we had to run the protoc build ourselves to get the protobuf
definition, rest gateway, and swagger doc:

```bash
protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. --grpc-gateway_out=logtostderr=true:$GOPATH/src --go_out=plugins=grpc:$GOPATH/src api/echo/echo.proto; cp api/echo/echo.swagger.json ./dist/swagger.json
```

If for some reason you want to avoid the buf tool, you will still need to do
something similar! I'm no longer updating that line though, so it may not
work in the future...

### Protobuf code generation with the buf tool

Now, however, we have the wonderful buf tool. In addition to generating out
code for us, it can also lint our protobuf and look for breaking changes.

The Dockerfile currently only builds it. If you want to lint or test for
breaking changes, you will have to do that separately. You can look at the
Dockerfile to see anything that needs to be installed in order to do so, though.

At some point I should produce an image to do these tests automatically.

### Local build

If you really want to build locally, the docker file has all the steps
and needed installs. I _could_ make it into a bash script as well, but
then I would need to duplicate changes to both so I probably will not do that.

### Docker builds

To instead build the docker container simply go to the echopilot directory and run

```bash
# This will probably require sudo on linux
docker build -t echopilot .
```

Then run via

```bash
# This will require a sudo if you are running on linux most likely.
docker run --rm -p 3000:3000 --mount type=bind,src="$(pwd)/configs/echopilot.json",target="/etc/echopilot/echopilot.json" --mount type=bind,src="$(pwd)/configs/ssl/cert.pem",target="/etc/echopilot/cert.pem" --mount type=bind,src="$(pwd)/configs/ssl/key.pem",target="/etc/echopilot/key.pem" echopilot
```

## Working with the elm part

This is actually really simple. Elm just compiles down to either a single .js or a single html that we can serve as static concent.

- Working directory: ui/
- Build: elm make src/Main.elm && mv index.html ../web/
- That's literally all there is to it

Note that the dockerfile builds the elm stuff for you, so you will not need
to do this yourself. If you DO want to try it, keep in mind that you will need
to install elm. See the Dockerfile for a command-line way to do so.

## Sources (in addition to logs of stack overflow and stuff I'm sure I lost track of)

These libraries are used to make our lives a little easier:

- <https://templ.guide>
- <https://github.com/go-chi/chi>
- <https://github.com/charmbracelet/log>
- <https://github.com/caarlos0/env>


general go patterns and tools:

- <https://github.com/thockin/go-build-template>
- <https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831>
- <https://github.com/golang-standards/project-layout>

testing:

- <https://blog.alexellis.io/golang-writing-unit-tests/>
- <https://github.com/benbjohnson/testing>
- <https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c>

signal handling:

- <https://github.com/benbjohnson/testing>
- <https://www.openmymind.net/Golang-Hot-Configuration-Reload/>
- <https://gravitational.com/blog/golang-ssh-bastion-graceful-restarts/>
- <https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2>

systemd hyjinks:

- <https://github.com/coreos/go-systemd>

docker stuff:

- <https://hub.docker.com/r/phusion/baseimage/>

protobuf and buf

- <https://grpc.io/docs/protoc-installation/>
- <https://github.com/bufbuild/buf>
- <https://docs.buf.build/tutorials/getting-started-with-buf-cli>

connect (grpc-compatable rpc protocol)

- <https://connect.build/docs/go/getting-started/>

go generate (just to make the build nicer):

- <https://eli.thegreenplace.net/2021/a-comprehensive-guide-to-go-generate/>

elm (no longer used. Replaced with htmx):

- <https://guide.elm-lang.org/>
- <https://package.elm-lang.org/>
- <https://elmprogramming.com/fetching-data-using-get.html>

grpc (no longer actively used, but good background reading):

- <https://github.com/grpc/grpc-go/tree/master/examples>
- <https://github.com/grpc-ecosystem/grpc-gateway>
- <https://github.com/grpc-ecosystem/grpc-opentracing>
- <https://levelup.gitconnected.com/grpc-basics-part-2-rest-and-swagger-53ec2417b3c4>
- <https://github.com/scottyw/grpc-example>
- <https://github.com/swagger-api/swagger-ui>
- <https://medium.com/@ribice/serve-swaggerui-within-your-golang-application-5486748a5ed4>
- <https://stackoverflow.com/questions/70643183/how-am-i-supposed-to-use-protoc-gen-go-grpc>

Things to investigate for the future(?):

- <https://gist.github.com/rivo/f96ad8710b54a49180a314ec4d68dbfb>
- <https://grisha.org/blog/2014/06/03/graceful-restart-in-golang/>

If I used your resources, no matter how small, thank you! It would have been beyond my patience to figure all of these things out by myself.
