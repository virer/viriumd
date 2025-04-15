## Usage:

$ export VG_NAME=vg_data
$ ./virium

## Create volume
curl -X POST http://localhost:8787/api/volumes/create \
  -H "Content-Type: application/json" \
  -d '{"name":"test-vol","capacity":10737418240}' # 10 GiB

# Delete volume
curl -X DELETE http://localhost:8787/api/volumes/delete \
    -H "Content-Type: application/json" \
    -d '{"VolumeID":"test-vol"}'

