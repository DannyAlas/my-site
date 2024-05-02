---
title: "Homelab"
date: 2024-04-01
author: "Danny"
tags: ["Homelab"]
slug: "homelab"
description: "A brief overview of all that resides in the self-hosted 2018 mid-tier consumer-grade-hardware cloud"
---

# My Homelab for 2024


This is more of a living document than a blog post and I'll keep editing and adding to it as things change. 

I thought about doing some version control, but hey, that's what git is for so [here's](https://github.com/DannyAlas/jazzy-cloud) the link.

# The Physical

## Nodes

![Optiplexes](/imgs/homelab/optiplex-7070.png)

| Components                  | Spec |
| :------------------------:    | :-----------------: |
| <i class="fa-solid fa-microchip"></i>     |   1x i7-9700 & 3x i7-8700  |
| <i class="fa-solid fa-memory"></i>        |   93 GiB   |
| <i class="fa-solid fa-hard-drive"></i> |  4x 256 SSD   |
| <i class="fa-solid fa-network-wired"></i> | 4x Gigabit Intel NIC |
| <i class="fa-solid fa-network-wired"></i> | 4x 10Gigabit ConnectX3 |
| <i class="fa-brands fa-windows"></i>   | pve/8.1.3/b46aac3b42da5d15 |

### Zeus
![Zeus](/imgs/homelab/zeus.png)
( I know, I know, I'm not good at naming things... )
|                Components                 |                       Spec                        |
| :---------------------------------------: | :-----------------------------------------------: |
|   <i class="fa-solid fa-microchip"></i>   |                   1x  i7-6700K                    |
|    <i class="fa-solid fa-memory"></i>     |                      32 GiB                       |
|   <i class="fa-solid fa-database"></i>    |      4x Seagate BarraCuda 4Tb Drive (RAIDZ1)      |
|  <i class="fa-solid fa-hard-drive"></i>   |      2x 2Tb random Drives (Pool 2 - RAIDZ1)       |
|    <i class="fa-solid fa-server"></i>     | NVIDIA RTX 3050 (Patched drivers for transcoding) |
| <i class="fa-solid fa-network-wired"></i> |               1x Gigabit Intel NIC                |
| <i class="fa-solid fa-network-wired"></i> |              1x 10Gigabit ConnectX3               |
| <i class="fa-brands fa-windows"></i> |  pve/8.1.3/b46aac3b42da5d15 |

I run the 4 Dell Optiplexes as [Proxmox](https://www.proxmox.com/en/proxmox-virtual-environment/overview) nodes only, I bought them used for a really good deal.

The, original lab was mainly just Zeus but now he's the granddaddy of all my servers. He's been through more disasters than a Greek tragedy. Still, he acts an extra node in the Prox cluster, my NAS and and runs one of the K8s masters (mainly for graphics related processes). I used to run TrueNas Scale for the NAS but moved to just a instance Debian running ZFS, as I didn't need most of the things TrueNas provided. 


![Proxmox Web Ui](/imgs/homelab/prox-webui.png)


Thus it's a **5 node cluster**, with local lvm thins on a each host. I'm experimenting with Ceph and HA but don't have enough OSDs for now to be production ready. Availability is mainly L7 based, but there are a few VMs in the Prox HA.

I run [Debian](https://www.debian.org/) wherever possible and also use it as my daily driver for my PC and laptop.

I use two separate bridges for networking, one with the gigabit cards for all VMs, LXCs (these are being moved to K8s), WebUI, and corosync. Then the 10G cards are on a separate bridge, no LACP. This is reserved for the Kubernetes cluster and Ceph/storage bandwidth. I setup spanning tree on the Mikrotik's so I can loose one cable and still be able to reach every node.


## Networking

![Homelab Network](/imgs/homelab/network-diagram.png)

### Router/Firewall

I have a fiber gigabit uplink to my ISP going into my main **firewall** [Protectli Vault FW4B](https://protectli.com/product/fw4b/) running PfSense in HA. I run the secondary instance on Zeus and even though pfsync is enabled, the WAN switch currently only terminates to the Protectli. Thus, I don't route any traffic through Zeus at the moment but if needed I could manually move over the WAN link.

By default everything is rejected both ways, I only open ports to the metallb IPs for my Istio Gateways and I control ingress/egress though them for all hosted apps. I also run Snort and have done a lot of *finagling* to get it working nicely, though it's never perfect. There are also three VPNs running, a Wireguard for personal uses and an OpenVPN server setup because a [specific type of network traffic](https://support.torproject.org/abuse/what-about-ddos/#:~:text=But%20because%20Tor%20only%20transports%20correctly%20formed%20TCP%20streams%2C%20not%20all%20IP%20packets%2C%20you%20cannot%20send%20UDP%20packets%20over%20Tor) only allows TCP streams which [Wiregaurd doesn't support](https://www.wireguard.com/known-limitations/#:~:text=WireGuard%20explicitly%20does%20not%20support%20tunneling%20over%20TCP). And lastly a Wiregaurd for my Work **and only work** traffic. My goal in the future is to get another fully dedicated firewall box for proper HA.

#### SW1 Core switch : MikroTik CRS309-1G-8S+IN

![](/imgs/homelab/mikrotik-CRS309-1G-8S+IN.webp)

Nothing fancy, just running SwitchOS.

#### SW2: W.I.P. about to be CRS326-24G-2S+RM

coming very soon...

#### AP

![](/imgs/homelab/WAX214PA.png)

Only one Access Point with 3 SSID, LAN, GUESTS and IOT. 

### VLANS

* LAN : Most *trusted things* are in here.
* IOT : smart tv, smart light bulbs, and anything I do not trust.
* GUEST : This one is for the GUEST Wifi and only provides WAN access.
* LARGE : PVE, PBS and the ZFS NAS and other storage is here
* ADMIN : Switches, AP, PBS & PVE Webui, SSH, secure VPN traffic, etc. 
* CLUSTER: Corosync

# Base services

## DNS...

I use two PiHoles for DNS filtering / blacklisting, two [PowerDNS Recursors](https://doc.powerdns.com/recursor/) and two [PowerDNS Authoritative Servers](https://doc.powerdns.com/authoritative/) connected to my PostgreSQL DB in the K8S cluster.

![](/imgs/homelab/home-dns-infra-dark-background.png)

Kea DHCP in my PfSense router pushes the two PiHoles to my clients, servers are using the re-cursors directly.

The PowerDNS solution is pretty overkill for my needs but it was a good learning experiance and if one of my VMs is ever down, at least my DNS stays running :). It does let me mess around a lot with routing and gather a good bit of metrics! I've also been using [PowerDNS Admin](https://github.com/PowerDNS-Admin/PowerDNS-Admin) as a nice way to manage the records.

## Git + Ops

Forgejo, Jenkins, and Harbor! All running in the K8S cluster. I'm slowly moving my projects to being hosted on Forgejo and just mirrored to GitHub.

## Kubernetes

Set up through some custome Ansible playbooks which was crucial while learning cause sure messed up a lot. Everything is configured with FluxCD (Lens is great b.t.w.). Initially I used Rancher + Helm charts for deployments but as complexity grew, I quickly learned why GitOps is nice.

This is still a W.I.P... I'll have a detailed write up about everything running in the future but as a simple overview for what hosted right now.

- Metallb
- Istio
- NFS Provisioner
- Longhorn
- Cloud Native Postgres (2x)
- KeyDB
- Kafka
- Minio
- Prometheus
- Grafana
- Loki
- Promtail
- Keycloak
- Hashicorp Vault
- *Arrs stack + Qbt + Plex (I use my own operator for this, based on [this wonderful project](https://github.com/kubealex/k8s-mediaserver-operator), I might release mine as public soon as it configures with Postgress for everything, even Qbt)
- Forgejo
- Harbor
- Airflow
- Flyte
- Jenkins
- Then all of my personal web services (this site, APIs, scrapers, etc. I'll also go into these more in the future)
