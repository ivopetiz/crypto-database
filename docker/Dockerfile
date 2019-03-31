FROM ubuntu:18.04

#ENV GRAFANA_VERSION 5.2.2
ENV INFLUXDB_VERSION 1.6.0

# Prevent some error messages
ENV DEBIAN_FRONTEND noninteractive

# Install all prerequisites
RUN	apt-get -y update && apt-get install -y apt-utils wget supervisor curl apt-utils \
	&& apt-get -y clean \
	&& rm -rf /var/lib/apt/lists/*

# ---------------- #
#   Installation   #
# ---------------- #

# Install Grafana to /src/grafana
#RUN		mkdir -p src/grafana && cd src/grafana && \
#			wget -nv https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana-${GRAFANA_VERSION}.linux-amd64.tar.gz -O grafana.tar.gz && \
#			tar xzf grafana.tar.gz --strip-components=1 && rm grafana.tar.gz

# Install InfluxDB
RUN	wget -nv https://dl.influxdata.com/influxdb/releases/influxdb_${INFLUXDB_VERSION}_amd64.deb && \
		dpkg -i influxdb_${INFLUXDB_VERSION}_amd64.deb && rm influxdb_${INFLUXDB_VERSION}_amd64.deb

# ----------------- #
#   Configuration   #
# ----------------- #

# Configure InfluxDB
COPY influxdb/config.toml /etc/influxdb/config.toml
COPY influxdb/run.sh /usr/local/bin/run_influxdb
# These two databases have to be created. These variables are used by set_influxdb.sh and set_grafana.sh
ENV PRE_CREATE_DB altcoin
ENV INFLUXDB_HOST localhost:8086
ENV INFLUXDB_DATA_USER machine
ENV INFLUXDB_DATA_PW passwd
#ENV INFLUXDB_GRAFANA_USER grafana
#ENV INFLUXDB_GRAFANA_PW grafana
ENV ROOT_PW root

# Configure Grafana
#ADD     ./grafana/config.ini /etc/grafana/config.ini
#ADD		grafana/run.sh /usr/local/bin/run_grafana
#ADD		./configure.sh /configure.sh
#ADD		./set_grafana.sh /set_grafana.sh
#ADD		./set_influxdb.sh /set_influxdb.sh
#RUN     /configure.sh
#
COPY ./supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# ----------- #
#   Cleanup   #
# ----------- #

RUN	apt-get autoremove -y wget curl
#
# ---------------- #
#   Expose Ports   #
# ---------------- #

# Grafana
#EXPOSE	3000

# InfluxDB Admin server
EXPOSE 8083

# InfluxDB HTTP API
EXPOSE 8086

# InfluxDB HTTPS API
EXPOSE 8084

# -------- #
#   Run!   #
# -------- #

CMD	["/usr/bin/supervisord"]