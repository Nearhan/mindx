FROM scratch
ADD app /
EXPOSE 5001
CMD ["/app"]