# Virium

Virium is a storage API built on top of LVM(Logical volume management) and iSCSI(Internet Small Computer Systems Interface) using targetcli-fb.
Virum is not scalable, nor a highly-available* storage solution

*a more or less HA solution would be possible using Linbit DRBD and a virtual IP using keepalived with VRRP protocol.

## Motivations and goals

The existing kubernetes external block storage solution was too complex to set up in my opinion for my homelab goals or have bad/strange implementation. I wanted a simple solution built on top of a strong foundation(LVM+iSCSI targetcli-fb). Also, the challenge to code a solution golang while I'm a golang beginner 

This is not a production ready solution and not intented to be (at least in the first steps).

## Requierements

Firewall:
 open TCP port 8787 (default customizable port)


## Build
$ ./buid.sh

## Usage

$ ./viriumd -v=2

## Example

### Create volume
$ curl -X POST http://localhost:8787/api/volumes/create \
  -H "Content-Type: application/json" \
  -u "virium_api_username:virium_api_password" \
  -d '{"initiator_name":"iqn.2025-04.net.virer.virium.test","capacity":10737418240}' # 10 GiB

### Delete volume
$ curl -X DELETE http://localhost:8787/api/volumes/delete \
    -H "Content-Type: application/json" \
    -u "virium_api_username:virium_api_password" \
    -d '{"volume_id":"47eb27cd-6824-4977-90fc-c62a21b11dfb"}'

