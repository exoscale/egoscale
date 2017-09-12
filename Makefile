VERSION=0.3.0-snapshot
PREFIX?=/usr/local
# GOPATH=$(PWD)/build:$(PWD)
PROGRAM=exo
# GO=env GOPATH=$(GOPATH) go
GO=go
RM?=rm -f
LN=ln -s
MAIN=cmd/exo.go
SRCS=	types.go		\
		error.go		\
		topology.go		\
		groups.go		\
		vm.go			\
		dns.go			\
		request.go		\
		async.go		\
		keypair.go		\
		ip.go			\
		init.go

all: $(PROGRAM)

$(PROGRAM): $(MAIN) $(SRCS)
				$(GO) build github.com/exoscale/egoscale
				$(GO) build -o $(PROGRAM) $(MAIN)

clean:
				$(RM) $(PROGRAM)
				$(GO) clean
