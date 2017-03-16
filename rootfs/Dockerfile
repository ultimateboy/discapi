FROM quay.io/deis/base:v0.3.5

# Add user and group
RUN adduser --system \
	--shell /bin/bash \
	--disabled-password \
	--home /opt/discapi \
	--group \
	discapi

COPY . /

# Fix some permission since we'll be running as a non-root user
RUN chown -R discapi:discapi /opt/discapi \
    && chmod +x /opt/discapi/sbin/*

USER discapi
WORKDIR /opt/discapi/sbin/
CMD ["discapi"]

ARG VERSION
ARG BUILD_DATE
ENV VERSION $VERSION
ENV BUILD_DATE $BUILD_DATE
ENV PATH $PATH:/opt/discapi/sbin