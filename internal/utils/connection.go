package utils

import "net"

// SubnetToIPAndMask - Функция производит конвертацию из вида xxx.xxx.xxx.xxx/xx
// в представление в виде начального адреса и маски подсети
//
// прим:
// 1.1.1.1/32 -> 1.1.1.1, 255.255.255.255
// 1.0.0.0/16 -> 1.0.0.0, 255.255.0.0
func SubnetToIPAndMask(subnet string) (ip string, mask string, err error) {
	ipValue, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", "", err
	}
	mask = net.IP(ipnet.Mask).String()
	return ipValue.String(), mask, nil
}
