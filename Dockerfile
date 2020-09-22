FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD RandomNumberGo_bot /
ADD config.json /
ADD resources.xml /
CMD ["chmod", "+x", "/RandomNumberGo_bot"]
CMD ["/RandomNumberGo_bot"]