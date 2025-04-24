# 📦 Viriumd Installation & Configuration Guide

**Viriumd** is a simple storage provisioning API built on top of Linux LVM and iSCSI (via `targetcli-fb`). This guide explains how to install and configure it.

---

## 🚀 Requirements

Before installing Viriumd, ensure the following tools and packages are available on your system:

- `lvm2`
- `targetcli-fb`
- `iscsi-initiator-utils` (on client side)

Open the following ports on your Firewall:

- TCP port 8787 (Virumd API)
- TCP port 3260 (iSCSI target)

Have free space inside your chosen Volume Group

## ⚙️ Installation

To install Virium you may download prebuilt binary from github or use rpmbuild to build from source(SPEC file provided in this repo).

```
curl --output viriumd.tar.gz -L https://github.com/virer/viriumd/releases/download/v0.2.7/viriumd_0.2.7_linux_amd64.tar.gz
tar -vzxf viriumd.tar.gz
```

## 🔧 Configuration

The configuration file is expected to be located at:

/etc/viriumd/virium.yaml

Here is a sample configuration with default values:

vg_name: "vg_data"
port: "8787"
iqn: "iqn.2025-04.net.virer.virium"
target_portal: "192.168.0.147:3260"
api_username: "virium_api_username"
api_password: "virium_api_password"

🔐 Parameters Explained
Name	Description
vg_name	The volume group name to use for LVM provisioning (free space expected to create new logical volume)
port	The port that viriumd API listens on
iqn	The iSCSI qualified name for your target
target_portal	IP address and port where the iSCSI target is exposed (so most probably the same IP where you deployed Viriumd)
api_username	Basic authentication username for API access
api_password	Basic authentication password for API access

## 🧪 Running the API

After configuration, you can start viriumd manually:

./viriumd -config /etc/viriumd/virium.yaml

Or enable it as a service (unit file in this repo: scripts/rpmbuild/SOURCE/viriumd.service):

sudo systemctl enable viriumd
sudo systemctl start viriumd

### Check status:

sudo systemctl status viriumd

### Logs:

sudo journalctl -u viriumd -f

## 🔐 Securing the API

Viriumd uses HTTP Basic Authentication. Make sure to:

    Change the default credentials

    Use a secure network

    (Optional) Use a reverse proxy (e.g. Nginx) to enable HTTPS  (future plan is to ship this option inside viriumd)


## 📬 API Usage

Once viriumd is running, you can interact with it via its REST API on:

http://<target_ip>:<port>/api/volumes/create

Authentication is required via HTTP basic auth using the configured credentials.    

### Examples

### Create volume

```
curl -X POST http://localhost:8787/api/volumes/create \
  -H "Content-Type: application/json" \
  -u "virium_api_username:virium_api_password" \
  -d '{"initiator_name":"iqn.2025-04.net.virer.virium.test","capacity":10737418240}' # 10 GiB
```

### Delete volume

```
curl -X DELETE http://localhost:8787/api/volumes/delete \
    -H "Content-Type: application/json" \
    -u "virium_api_username:virium_api_password" \
    -d '{"volume_id":"47eb27cd-6824-4977-90fc-c62a21b11dfb"}'
```

## 🧹 Troubleshooting

    No such VG: Ensure vg_name exists and is active.

    iSCSI not working: Check targetcli configuration(targetli ls) and firewalls.

    Permissions: The viriumd binary may require sudo to interact with LVM or iSCSI, or you can adjust the system permissions accordingly.

### ✨ Tips

To avoid to saturate LVM snapshots, please configure LVM snapshot autoextend:
```
sed -i 's/\(.*\)# snapshot_autoextend_threshold = 70/\1snapshot_autoextend_threshold = 70/g' /etc/lvm/lvm.conf 
```