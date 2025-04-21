# Virium

**Viriumd** is a lightweight storage API built on top of **LVM** (Logical Volume Manager) and **iSCSI** (Internet Small Computer Systems Interface) using `targetcli-fb`.

> ⚠️ Virium is **not** a scalable or highly available storage solution by default. It is designed for **homelab, testing, or demo** environments.

That said, a *somewhat* highly available setup could be achieved using tools like **Linbit DRBD** for replication and **keepalived** with **VRRP** for virtual IP failover.

## ✨ Motivation & Goals

I created Virium because I found most Kubernetes external storage drivers too complex or oddly implemented for my homelab needs. I wanted:

- A **simple solution**
- Built on **solid, well-known tools** (LVM + iSCSI via `targetcli-fb`)
- An opportunity to **learn and practice Golang** from scratch

This project is **not production-ready** and is **not intended to be**, at least in its early stages.

## 🔧 Features (current or planned)

- REST API for managing LVM volumes and exposing them via iSCSI
- CSI driver support for Kubernetes integration
- Snapshot support via LVM
- Simple YAML-based configuration
- Systemd + RPM packaging

## 🚫 Not Included (by design)

- High-availability or failover (though possible with DRBD + keepalived)
- Multi-node storage orchestration

## 🧪 Ideal Use Cases

- Homelabs
- Learning/experimentation
- Demo environments
- Lightweight Kubernetes clusters with external CSI driver

## Requierements

Firewall:
 - open TCP port 8787 (default customizable port)

Packages:
 - lvm2
 - targetcli
 - target-restore

## Build
$ ./scripts/buid.sh

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
