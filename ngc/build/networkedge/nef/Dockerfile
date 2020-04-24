# SPDX-License-Identifier: Apache-2.0
# Copyright © 2019 Intel Corporation

FROM centos:7.6.1810

RUN rm -rf /etc/yum.repos.d/*
RUN cd /etc/yum.repos.d && curl -o CentOS-Base.repo "http://mirrors.aliyun.com/repo/Centos-7.repo"
RUN cd /etc/yum.repos.d && curl -o epel.repo "http://mirrors.aliyun.com/repo/epel-7.repo"
RUN yum clean all

ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
ENV no_proxy=localhost,127.0.0.1,nefservice,afservice

RUN yum install -y wget\
	net-tools \ 
	nmap \
        && yum clean all

WORKDIR /root
COPY ./nef . 
COPY ./entrypoint.sh .

RUN chmod +x ./nef
RUN chmod +x ./entrypoint.sh

CMD ["./entrypoint.sh"]

