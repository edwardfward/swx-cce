FROM redis:alpine
COPY redis-config.sh /redis-config.sh
RUN apk --no-cache add bash
RUN chmod 700 /redis-config.sh
CMD [ "/redis-config.sh" ]