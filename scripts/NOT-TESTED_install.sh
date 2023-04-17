# install script for new droplet. (assumes root)

# Add charm.sh apt repository
mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | tee /etc/apt/sources.list.d/charm.list
apt update

# Install charm applications
# - softserve (git server w/ TUI)
# - wishlist (ssh gateway)
# - gum (enables better interavtive shell-scripts)
sudo apt install -y softserve wishlist gum

## mount digital ocean volume
# Create a mount point for your volume:
$ mkdir -p /mnt/volume_fra1_01

# Mount your volume at the newly-created mount point:
$ mount -o discard,defaults,noatime /dev/disk/by-id/scsi-0DO_Volume_volume-fra1-01 /mnt/volume_fra1_01

# Change fstab so the volume will be mounted after a reboot
$ echo '/dev/disk/by-id/scsi-0DO_Volume_volume-fra1-01 /mnt/volume_fra1_01 ext4 defaults,nofail,discard 0 0' | sudo tee -a /etc/fstab
