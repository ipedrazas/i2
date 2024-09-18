package prxmx

import (
	"context"
	"fmt"
	"strings"

	"github.com/luthermonson/go-proxmox"
)

type Cluster struct {
	ApiURL string
	User   string
	Pass   string
	Client *proxmox.Client
}

type Node struct {
	Name    string
	IP      []string
	Uptime  Uptime
	Running bool
}

type Application struct {
	Name      string
	PublicUrl string
	Port      int
}

func NewProxmoxClient(apiURL, user, pass string) *proxmox.Client {
	credentials := proxmox.Credentials{
		Username: user,
		Password: pass,
	}
	client := proxmox.NewClient(apiURL,
		proxmox.WithCredentials(&credentials),
	)
	return client
}

func NewCluster(apiURL, user, pass string) *Cluster {
	// credentials := proxmox.Credentials{
	// 	Username: user,
	// 	Password: pass,
	// }
	client := proxmox.NewClient(apiURL,
		// proxmox.WithCredentials(&credentials),
		proxmox.WithAPIToken(user, pass),
	)
	return &Cluster{
		ApiURL: apiURL,
		User:   user,
		Pass:   pass,
		Client: client,
	}
}

func (c *Cluster) GetClusterNodes() ([]string, error) {
	nodes, err := c.Client.Nodes(context.Background())
	if err != nil {
		return nil, err
	}
	nodeNames := make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Node+" "+node.IP)
	}
	return nodeNames, nil
}

func (c *Cluster) getVirtualMachines() ([]*proxmox.VirtualMachine, error) {
	nodes, err := c.Client.Nodes(context.Background())
	if err != nil {
		return nil, err
	}
	VMs := []*proxmox.VirtualMachine{}
	ctx := context.Background()

	for _, nodeStatus := range nodes {
		node, err := c.Client.Node(ctx, nodeStatus.Node)
		if err != nil {
			return nil, err
		}
		vms, err := node.VirtualMachines(ctx)
		if err != nil {
			return nil, err
		}
		for _, vm := range vms {
			VMs = append(VMs, vm)
		}
	}
	return VMs, nil
}
func (c *Cluster) GetVMs(all bool) ([]Node, error) {

	VMs := []Node{}

	allVms, err := c.getVirtualMachines()
	if err != nil {
		return nil, err
	}

	for _, vm := range allVms {
		if vm.Template {
			continue
		}
		node := Node{
			Name:    vm.Name,
			Uptime:  ParseUptime(vm.Uptime),
			Running: vm.Status == "running",
		}
		if node.Running {
			node.IP = getIPs(vm)
		}
		VMs = append(VMs, node)
		// if all {
		// 	VMs = append(VMs, node)
		// } else {
		// 	if node.Running {
		// 		node.IP = getIPs(vm)
		// 		VMs = append(VMs, node)
		// 	}
		// }

	}

	return VMs, nil
}

func getIPs(vm *proxmox.VirtualMachine) []string {
	ifaces, err := vm.AgentGetNetworkIFaces(context.Background())
	if err != nil {
		return nil
	}

	ips := make([]string, 0, len(ifaces))
	for _, iface := range ifaces {
		for _, ip := range iface.IPAddresses {
			if ip.IPAddressType == "ipv4" {
				ips = append(ips, ip.IPAddress)
			}
		}
	}
	return ips
}

type Uptime struct {
	Seconds uint64
	Minutes uint64
	Hours   uint64
	Days    uint64
	Raw     uint64
}

func ParseUptime(uptime uint64) Uptime {
	seconds := uptime % 60
	minutes := (uptime / 60) % 60
	hours := (uptime / 3600) % 24
	days := (uptime / 86400) % 30
	return Uptime{
		Seconds: seconds,
		Minutes: minutes,
		Hours:   hours,
		Days:    days,
		Raw:     uptime,
	}
}

func (c *Node) ToString() string {
	if c.Running {
		return fmt.Sprintf("%s - %s - %s ", c.Name, strings.Join(c.IP, ","), c.Uptime.ToString())
	}
	return fmt.Sprintf("%s - %v", c.Name, c.Running)
}

func (u *Uptime) ToString() string {
	return fmt.Sprintf("%d days, %d h, %d ', %d ''", u.Days, u.Hours, u.Minutes, u.Seconds)
}

func (u *Uptime) ToStringShort() string {
	if u.Days > 0 {
		return fmt.Sprintf("%d days", u.Days)
	}
	if u.Hours > 0 {
		return fmt.Sprintf("%d h", u.Hours)
	}
	if u.Minutes > 0 {
		return fmt.Sprintf("%d min", u.Minutes)
	}
	return fmt.Sprintf("%d s", u.Seconds)
}
