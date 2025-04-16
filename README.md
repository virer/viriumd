# Virium

Virium is a storage API built on top of LVM(Logical volume management) and iSCSI(Internet Small Computer Systems Interface) using targetcli-fb.
Virum is not scalable, nor highly-available* storage solution

*a more or less HA solution would be possible using Linbit DRBD and a virtual IP using keepalived with VRRP protocol.

## Motivations and goals

The existing kubernetes external block storage solution was too complex to set up in my opinion for my homelab goals or have bad/strange implementation. I wanted a simple solution built on top of a strong foundation(LVM+iSCSI targetcli-fb). Also, the challenge to code a solution golang while I'm a golang beginner 

This is not a production ready solution and not intented to be (at least in the first steps).

## Build
$ CGO_ENABLED=0 GOOS=linux go build -o tmp/viriumd

## Usage:

$ export VG_NAME=vg_data
$ ./viriumd

## Example

### Create volume
curl -X POST http://localhost:8787/api/volumes/create \
  -H "Content-Type: application/json" \
  -d '{"name":"test-vol","capacity":10737418240}' # 10 GiB

# Delete volume
curl -X DELETE http://localhost:8787/api/volumes/delete \
    -H "Content-Type: application/json" \
    -d '{"VolumeID":"test-vol"}'

