FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-pulumi"]
COPY baton-pulumi /