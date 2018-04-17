FROM progrium/busybox
RUN opkg-install nmap bash
ADD main /
EXPOSE 8080
CMD ["/main"]
