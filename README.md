### docker build
``` docker build . -t <image>:<tag> ```

### or pull 
``` docker pull riet/aliyun-rds-exporter ```

### run
```
docker run \ 
  -d \ 
  --name aliyun-slb-exporter \
  -e ACCESS_KEY_ID=<aliyun ak> \
  -e ACCESS_KEY_SECRET=<aliyun ak sk> \
  -e REGION_ID=<region id> \
  -p 10001:10001 \
  riet/aliyun-rds-exporter 
```

visit http://localhost:10001/metrics