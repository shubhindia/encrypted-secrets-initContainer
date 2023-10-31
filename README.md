# InitConainer for encrypted-secrets

## Introduction
#### Note: This is just a POC, plan is to have a mutating webhook which will automatically inject the initContainer in the pod when necessary annotations are present.


This is a initContainer which will decrypt the secrets and mount them in the pod. This is useful when you want to use encrypted secrets in your pod and don't want to store the decrypted secrets in the cluster.

## How it works
This initContainer will decrypt the secrets and write them to `/tmp/opt/conf/{key}`. Since, the same volume is mounted in the app container, these files can be read by the app container and secrets can be used.