# Network Bootstrap

KWS_Core supports two network backends: **host bridge** and **OVN**. The mode is selected at build time.

```bash
make build NETWORK=bridge   # host bridge (default fallback)
make build NETWORK=ovn      # OVN/OVS overlay (default)
```

`NetworkMode` is injected via `-ldflags` and controls which libvirt XML interface configuration is generated for new VMs.

---

## Host Bridge (`NETWORK=bridge`)

VMs attach to `virbr0`, libvirt's built-in NAT bridge. No additional infrastructure is required.

### How it works

```
VM (virtio NIC)
    │
virbr0  (e.g., 192.168.122.1/24, linux managed)
    │
iptables MASQUERADE
    │
Host physical NIC → external network
```

- libvirt creates and manages `virbr0` automatically on install.
- VMs receive IPs via dnsmasq (libvirt's built-in DHCP).
- Outbound traffic is NAT'd through the host.
- No inter-host VM routing — all VMs are isolated to a single host.

### Setup

No extra steps. `make conf` + `make build NETWORK=bridge` is sufficient.

### Pros

- Zero infrastructure overhead — works out of the box after libvirt install.
- Simple to debug: standard Linux bridge, visible with `brctl show`.
- No dependency on OVS/OVN build or cluster state.

### Cons

- VMs are NAT'd: no direct inbound access without port forwarding rules.
- Single-host only: VMs on different hosts cannot communicate directly.
- No SDN features: no logical networks, ACLs, or distributed routing.
- Not suitable for multi-node clusters.

---

## OVN (`NETWORK=ovn`)

VMs attach to `br-int`, an Open vSwitch integration bridge managed by OVN. Provides a distributed overlay network across multiple physical hosts using Geneve tunnels.

### Architecture

```
                    ┌─────────────────────────────────────────┐
                    │            Control Node                  │
                    │                                          │
                    │  ovn-northd ── OVN NB DB (:6641)        │
                    │       │                                  │
                    │       └──── OVN SB DB (:6642)           │
                    │                                          │
                    │  br-ext ── eno1 (physical NIC)          │
                    └──────────────────┬──────────────────────┘
                                       │ Geneve tunnel (UDP 6081)
              ┌────────────────────────┼────────────────────────┐
              │                        │                         │
   ┌──────────┴──────────┐  ┌──────────┴──────────┐            ...
   │     Worker Node      │  │     Worker Node      │
   │                      │  │                      │
   │  VM ── br-int        │  │  VM ── br-int        │
   │         │            │  │         │            │
   │  ovn-controller      │  │  ovn-controller      │
   │         │            │  │         │            │
   │  br-ext ── ens3      │  │  br-ext ── ens3      │
   └──────────────────────┘  └──────────────────────┘
```

**OVN NB DB**: Logical network intent (logical switches, routers, ACLs).
**OVN SB DB**: Physical binding state (which VM is on which chassis).
**ovn-northd**: Translates logical intent → physical flow rules.
**ovn-controller**: Per-host agent, programs OVS flow tables.
**br-int**: Integration bridge. All VM NICs attach here.
**br-ext**: Uplink bridge. Bridges the physical NIC for external traffic.

For Integration, You should use ovn-go-cms package.
refer in Knet-easy-cloud organization.

VM traffic path (cross-host):
```
VM NIC → br-int → Geneve encapsulation → physical NIC → network → remote host → br-int → destination VM
```

MTU is set to **1450** on VM interfaces (1500 − 50 byte Geneve overhead).

### Setup

#### 1. Build OVN/OVS from source (once per host)

```bash
make ovn-install
```

Clones `ovn-org/ovn`, builds OVS submodule, installs to `/usr/local`.

#### 2. Bootstrap control node

```bash
make ovn-cluster IP=<node-ip> DNS=<gateway-ip> [IFACE=eno1]
```

- Creates OVN NB/SB databases.
- Starts `openvswitch`, `ovn-nb`, `ovn-sb`, `ovn-northd`, `ovn-controller` as systemd services.
- Creates `br-ext` and attaches `IFACE` (default: `eno1`) as uplink.
- Configures `20-br-ext.network` for systemd-networkd.


install ovn-go-cms

`git clone https://github.com/easy-cloud-Knet/ovn-go-cms.git` 

#### 3. Bootstrap worker nodes

```bash
make ovn-worker IP=<node-ip> DNS=<gateway-ip> [CTL_IP=10.5.15.39] [IFACE=ens3]
```

- Starts `openvswitch` and `ovn-controller` as systemd services.
- Creates `br-ext`, attaches `IFACE` (default: `ens3`) as uplink.
- Registers with the control node at `CTL_IP:6642`.
- Sets `ovn-encap-type=geneve` and `ovn-encap-ip` for tunnel endpoint.

#### 4. Build KWS_Core in OVN mode

```bash
make build                  # NETWORK=ovn is the default
```

#### Verification

```bash
# Control node
sudo ovn-sbctl show         # should list all registered chassis
sudo ovs-vsctl show

# Worker node
sudo systemctl status openvswitch ovn-controller
```

### Pros

- Multi-host VM networking: VMs on different physical nodes communicate transparently.
- Full SDN capability: logical switches, routers, ACLs, NAT via OVN northbound API.
- Flat overlay: VMs get routable overlay IPs regardless of physical topology.
- Scalable: adding a worker node only requires running `make ovn-worker`.

### Cons

- Significant infrastructure overhead: OVN must be built from source, three databases and two daemons per control node, one agent per worker.
- Operational complexity: cluster state must be healthy for VM networking to function; a failed SB DB stalls all flow updates.
- Geneve overhead: effective MTU is 1450 (50 bytes lower than standard Ethernet), which can cause issues with applications assuming 1500 MTU.
- Build time: `make ovn-install` compiles OVS and OVN from source, which takes 10–20 minutes.
- Requires `systemd-networkd`: the `br-ext` interface is configured via `.network` files; `systemd-networkd-wait-online` is masked in base images to avoid 2-minute boot timeout.

---

## Comparison

|                          | Host Bridge          | OVN                    |
| ------------------------ | -------------------- | ---------------------- |
| Setup complexity         | Low                  | High                   |
| Multi-host VM networking | No                   | Yes                    |
| SDN features             | No                   | Yes                    |
| External IP for VMs      | NAT only             | Overlay (configurable) |
| VM MTU                   | 1500                 | 1450                   |
| Infrastructure daemons   | 0                    | 5+ (per cluster)       |
| Build time               | —                    | 10–20 min (source)     |
| Suitable for             | Single-host dev/test | Multi-node production  |
