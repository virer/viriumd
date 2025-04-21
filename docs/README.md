### Tips
to avoid to saturate snapshots: configure LVM snapshot autoextend
$ sed -i 's/\(.*\)# snapshot_autoextend_threshold = 70/\1snapshot_autoextend_threshold = 70/g' /etc/lvm/lvm.conf 
