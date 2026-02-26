#!/bin/bash

# AppArmor configuration script for KWS
# This script configures AppArmor to allow libvirt access to /var/lib/kws/

set -e

echo "Configuring AppArmor for KWS..."

# File 1: /etc/apparmor.d/local/usr.lib.libvirt.virt-aa-helper
VIRT_AA_HELPER="/etc/apparmor.d/local/usr.lib.libvirt.virt-aa-helper"
echo "Updating $VIRT_AA_HELPER..."

if [ ! -f "$VIRT_AA_HELPER" ]; then
    echo "Error: $VIRT_AA_HELPER not found. Please ensure libvirt is properly installed."
    exit 1
fi

# Check if entries already exist
if ! sudo grep -q "/var/lib/kws/" "$VIRT_AA_HELPER" 2>/dev/null; then
    echo "Adding KWS paths to $VIRT_AA_HELPER..."
    sudo tee -a "$VIRT_AA_HELPER" > /dev/null << 'EOF'
# Allow access to KWS storage
  /var/lib/kws/ r,
  /var/lib/kws/** rwk,
EOF
else
    echo "KWS paths already exist in $VIRT_AA_HELPER"
fi

# File 2: /etc/apparmor.d/libvirt/TEMPLATE.qemu
TEMPLATE_QEMU="/etc/apparmor.d/libvirt/TEMPLATE.qemu"
echo "Updating $TEMPLATE_QEMU..."

if [ ! -f "$TEMPLATE_QEMU" ]; then
    echo "Error: $TEMPLATE_QEMU not found. Please ensure libvirt is properly installed."
    exit 1
fi

# Check if entries already exist
if ! sudo grep -q "/var/lib/kws/" "$TEMPLATE_QEMU" 2>/dev/null; then
    echo "Adding KWS paths to $TEMPLATE_QEMU..."
    # Insert before the closing brace of the profile
    sudo sed -i '/^}[[:space:]]*$/i \  # Allow access to KWS storage\n  /var/lib/kws/ r,\n  /var/lib/kws/** rwk,' "$TEMPLATE_QEMU"
else
    echo "KWS paths already exist in $TEMPLATE_QEMU"
fi

# Reload AppArmor profiles
echo "Reloading AppArmor profiles..."
sudo systemctl reload apparmor || true

echo "AppArmor configuration completed successfully!"
