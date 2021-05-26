package util

import (
	"fmt"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
	"github.com/whyrusleeping/multiaddr-filter"
	"net"
	"difoss-goutil/util/enum"
)

// DEPRECATED:
type CIDRAddress string

func (s CIDRAddress) ToAddress() (*Address, error) {
	return NewAddress(string(s))
}

func GetIPByInterfaceName(interfaceNames ...string) ([]*Address, error) {
	ips := make([]*Address, 0)
	for _, iName := range interfaceNames {
		netInterface, err := net.InterfaceByName(iName)
		if err != nil {
			return nil, err
		}
		addresses, err := GetIPByInterface(netInterface)
		if err != nil {
			return ips, err
		}
		ips = append(ips, addresses...)
	}
	return ips, nil
}

func GetIPByInterface(netInterfaces ... *net.Interface) ([]*Address, error) {
	ips := make([]*Address, 0)
	for _, netInterface := range netInterfaces {
		addresses, err := netInterface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, netAddr := range addresses {
			addr, err := NewAddress(netAddr.String())
			if err != nil {
				return ips, err
			}
			if addr != nil {
				ips = append(ips, addr)
			}
		}
	}
	return ips, nil
}

func GetInterfaceName(flags ...net.Flags) (names []string, e error) {
	interfaces, err := GetInterface(flags...)
	if err != nil {
		return names, err
	}
	for _, iFace := range interfaces {
		names = append(names, iFace.Name)
	}
	return
}

func GetInterface(flags ...net.Flags) ([]net.Interface, error) {
	is, err := net.Interfaces()
	if err != nil {
		return nil, nil
	}
	var filter net.Flags
	for _, f := range flags {
		filter |= f
	}
	var result []net.Interface
	for _, netInterface := range is {
		if filter != 0 {
			if netInterface.Flags&filter != filter {
				continue // skip the network interface whose flags do not match the filter
			}
		}
		result = append(result, netInterface)
	}
	return result, nil
}

func GetIPs() (map[string][]*Address, error) {
	m := make(map[string][]*Address, 0)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			if addr, err := NewAddress(v.String()); err == nil {
				if addrSlice, ok := m[i.Name]; ok {
					m[i.Name] = append(addrSlice, addr)
				} else {
					m[i.Name] = []*Address{addr}
				}
			}
		}
	}
	return m, nil

}

type IpType uint8

var ipTypeWrapper *enum.Wrapper

const (
	IpTypeUnknown = iota
	IpTypeV4
	IPTypeV6
)

func init() {
	ipTypeWrapper = enum.NewEnumWrapper("IpType",
		enum.Item{Value: int(IpTypeUnknown), Name: "unknown"},
		enum.Item{Value: int(IpTypeV4), Name: "ipv4"},
		enum.Item{Value: int(IPTypeV6), Name: "ipv6"},
	)
}

func (t IpType) String() string {
	return ipTypeWrapper.GetName(int(t))
}

func GetIpType(ip net.IP) IpType {
	if ip.To4() != nil {
		return IpTypeV4
	}
	if ip.To16() != nil {
		return IPTypeV6
	}
	return IpTypeUnknown
}

type Address struct {
	net.IP
	*net.IPNet
}

func NewAddress(CIDR string) (*Address, error) {
	if ip, ipNet, err := net.ParseCIDR(CIDR); err != nil {
		return nil, err
	} else {
		return &Address{
			IP:    ip,
			IPNet: ipNet,
		}, nil
	}
}

func (a Address) GetIpType() IpType {
	return GetIpType(a.IP)
}

func (a Address) ConvertIPNet() (string, error) {
	return mask.ConvertIPNet(a.IPNet)
}

func (a Address) String() (CIDR string) {
	ones, _ := a.IPNet.Mask.Size()
	return fmt.Sprintf("%s/%d", a.IP.String(), ones)
}

func (a Address) MarshalJSON() ([]byte, error) {
	s := `"` + a.String() + `"`
	return []byte(s), nil
}

func (a *Address) UnmarshalJSON(data []byte) error {
	if addr, e := NewAddress(string(data)); e != nil {
		return e
	} else {
		a.IP = addr.IP
		a.IPNet = addr.IPNet
		return nil
	}
}

func (a Address) ToMultiaddr(returnErr ... *error) multiaddr.Multiaddr {
	mAddr, e := manet.FromIP(a.IP)
	if len(returnErr) != 0 {
		// use 1st only
		returnErr[1] = &e
	}
	return mAddr
}

func (a *Address) ToTcpAddr(portOrNone ...int) (*net.TCPAddr, error) {
	port := 0
	if len(portOrNone) != 0 {
		port = portOrNone[0]
	}
	return net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", a.IP.String(), port))
}
