FROM golang:1.10

RUN apt-get update && apt-get install -y --no-install-recommends \
	    python \
        virtualenv

# install asdf
RUN curl -L https://github.com/asdf-vm/asdf/archive/v0.5.0.tar.gz > asdf.tar.gz && \
    tar -xzvf asdf.tar.gz && \
    mv asdf-0.5.0 /asdf && \
    rm asdf.tar.gz

ENV PATH /asdf/shims:/asdf/bin:$PATH
