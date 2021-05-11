FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD ./deploy/amd64* /
CMD ["chmod", "+x", "/RandomNumberGo_bot"]
CMD ["/RandomNumberGo_bot"]