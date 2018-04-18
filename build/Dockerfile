FROM progrium/busybox
RUN opkg-install nmap
ADD main /
EXPOSE 8080
CMD ["/main"]
