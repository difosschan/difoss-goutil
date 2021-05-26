package http_util

import (
	"bytes"
	"context"
	"difoss-goutil/frame"
	"difoss-goutil/util"
	"errors"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func NewFileUploadRequest(url, filedName string, path string, params map[string]string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(filedName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	// 其他参数列表写入 body
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}

func HttpGetFromIP(url, ipv4Addr string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				lAddr, err := net.ResolveTCPAddr(network, "["+ipv4Addr+"]:0")
				if err != nil {
					return nil, err
				}
				// 被请求的地址
				rAddr, err := net.ResolveTCPAddr(network, addr)
				if err != nil {
					return nil, err
				}
				conn, err := net.DialTCP(network, lAddr, rAddr)
				if err != nil {
					return nil, err
				}
				deadline := time.Now().Add(10 * time.Second)
				conn.SetDeadline(deadline)
				return conn, nil
			},
			ResponseHeaderTimeout: frame.DefaultResponseHeaderTimeout,
		},
	}
	return client.Do(req)
}

func CheckServer(address string, timeout time.Duration) (ok bool, err error) {
	_, err = net.DialTimeout("tcp", address+":443", timeout)
	return err == nil, err
}

func CheckServerByInterface(interfaceName, address string, timeout time.Duration) (ok bool, err error) {
	localAddresses, e := util.GetIPByInterfaceName(interfaceName)
	if e != nil {
		return false, e
	}
	var ipv4LocalAddr util.Address
	for _, localAddr := range localAddresses {
		if localAddr.GetIpType() == util.IpTypeV4 {
			ipv4LocalAddr = *localAddr
			break
		}
		return false, errors.New("no ip4v address under network interface")
	}
	return CheckServerByLocalAddr(&ipv4LocalAddr, address, timeout)
}

func CheckServerByLocalAddr(ipv4LocalAddr *util.Address, address string, timeout time.Duration) (ok bool, err error) {
	tcpAddr, e := ipv4LocalAddr.ToTcpAddr()
	if e != nil {
		return false, e
	}
	d := net.Dialer{
		Timeout:   timeout,
		LocalAddr: tcpAddr,
	}
	if _, err := d.Dial("tcp", address+":443"); err != nil {
		return false, err
	}
	return true, nil
}
