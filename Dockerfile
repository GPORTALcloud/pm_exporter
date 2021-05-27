FROM fedora:34

COPY pm_exporter /usr/local/bin/pm_exporter
RUN chmod +x /usr/local/bin/pm_exporter

ENTRYPOINT [ "/usr/local/bin/pm_exporter" ]
