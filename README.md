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

## 🚀 Requirements

Firewall:
 - open TCP port 8787 (default customizable port)

Packages:
 - lvm2
 - targetcli
 - target-restore

## Docs, Installation and usage

Please check "docs" directory

## Looking for the CSI Driver?

Please check the [Helm charts here](https://github.com/virer/virium-helm-repo/tree/main)

## License

Viriumd is released under Apache License Version 2.0