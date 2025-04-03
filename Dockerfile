FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-pulumi-cloud"]
COPY baton-pulumi-cloud /