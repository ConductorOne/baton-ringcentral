FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-ringcentral"]
COPY baton-ringcentral /