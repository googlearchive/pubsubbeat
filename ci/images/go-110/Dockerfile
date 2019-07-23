FROM golang:1.12

RUN apt-get update 
RUN apt-get install -y --no-install-recommends python3.6
RUN apt-get install -y --no-install-recommends pip3-python
RUN pip3 install virtualenv

# install asdf
RUN curl -L https://github.com/asdf-vm/asdf/archive/v0.5.0.tar.gz > asdf.tar.gz && \
    tar -xzvf asdf.tar.gz && \
    mv asdf-0.5.0 /asdf && \
    rm asdf.tar.gz

ENV PATH /asdf/shims:/asdf/bin:$PATH
