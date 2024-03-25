FROM public.ecr.aws/zinclabs/zinc:latest

ENV ZINC_DATA_PATH="/data" \
    ZINC_FIRST_ADMIN_USER="admin" \
    ZINC_FIRST_ADMIN_PASSWORD="admin"

EXPOSE 4080

VOLUME /data

CMD ["zincsearch"]
