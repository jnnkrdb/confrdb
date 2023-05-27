# ConfRDB Operator

```bash
operator-sdk init --domain jnnkrdb.de --repo github.com/jnnkrdb/configrdb
```


```bash
operator-sdk create api --group globals --version v1beta2 --kind GlobalConfig --resource --controller
operator-sdk create api --group globals --version v1beta2 --kind GlobalSecret --resource --controller
```

```bash
operator-sdk build docker.io/jnnkrdb/confrdb:v0.0.1
docker push docker.io/jnnkrdb/confrdb:v0.0.1
```