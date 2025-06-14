Name:           viriumd
Version:        0.2.9
Release:        1%{?dist}
Summary:        Viriumd - CSI Storage API server using iSCSI and LVM

License:        MIT
URL:            https://github.com/virer/viriumd
Source0:        https://github.com/virer/viriumd/viriumd-%{version}.tar.gz
BuildArch:      x86_64
BuildRequires:  golang
Requires:       targetcli
Requires:       lvm2

%global debug_package %{nil}

%description
Viriumd is an API server that manages LVM and iSCSI-based volumes for use with Kubernetes CSI drivers. It exposes a RESTful API for volume creation, deletion, and snapshots.

%prep
%setup -q

%build
go build -o viriumd ./cmd/viriumd  

%install
install -D -m 0755 viriumd %{buildroot}/usr/bin/viriumd
install -D -m 0644 config/viriumd.service %{buildroot}/usr/lib/systemd/system/viriumd.service
install -D -m 0644 config/virium.yaml %{buildroot}/etc/viriumd/virium.yaml

%files
%license LICENSE
%doc docs/README.md
/usr/bin/viriumd
/etc/viriumd/virium.yaml
/usr/lib/systemd/system/viriumd.service

%post
# Enable the service on install
%systemd_post viriumd.service

%preun
%systemd_preun viriumd.service

%postun
%systemd_postun_with_restart viriumd.service

%changelog
* Thu Apr 24 2025 Sebastien Caps <virer@hotmail.com> - 0.2.9-1
- Updated source files and fix prep for GHA

* Thu Apr 24 2025 Sebastien Caps <virer@hotmail.com> - 0.2.7-1
- Updated source version and change doc file

* Sat Apr 19 2025 Sebastien Caps <virer@hotmail.com> - 0.1.0-1
- Initial RPM release