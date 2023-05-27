# ----------------------------------------------- GO-BUILD
# BUILDING STAGE
FROM golang:1.20-alpine 
RUN apk add --no-cache curl make bash
RUN export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.28.1
RUN curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.28.1/operator-sdk_linux_amd64
RUN chmod +x operator-sdk_linux_amd64 
RUN mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk