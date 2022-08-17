FROM node:lts

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /root/

RUN apt-get update
RUN apt-get install -y \
    build-essential \
    ca-certificates \
    openssl \ 
    libasound2-dev \
    npm \
    uuid-dev \
    wget \
  && update-ca-certificates

RUN wget https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
RUN tar -xf go1.18.3.linux-amd64.tar.gz -C /usr/local/bin
ENV PATH="/usr/local/bin/go/bin:${PATH}"

RUN wget https://github.com/pjsip/pjproject/archive/refs/tags/2.12.tar.gz
RUN tar -xf 2.12.tar.gz -C /root/

WORKDIR /root/pjproject-2.12/
RUN ./configure
RUN make dep -j8
RUN make -j8
RUN make install

WORKDIR /root/ui/

COPY ui .
RUN npm install --silent --production
RUN npm run build

WORKDIR /root/server/

COPY server .
RUN npm install --silent --production
RUN npm run build

WORKDIR /root
RUN cp -R ./ui/build ./public

EXPOSE 4000

CMD ["server/build/caller_back"]