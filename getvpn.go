/*
 * This file is part of the VPNGateFetcher distribution (https://github.com/prben/vpngatefetcher).
 * Copyright (c) 2020 Pr Ben.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/jianfeng6909/vpngate"
)

var (
	openvpnBin, proto string
	latency           int64
	minSpeed          int
)

func init() {
	flag.StringVar(&openvpnBin, "openvpnBin", "/usr/sbin/openvpn", "Mention path of OpenVPN binary")
	flag.StringVar(&proto, "proto", "any", "Preferred Protocol - tcp/udp/any")
	flag.IntVar(&minSpeed, "minSpeed", 40, "Minimum speed of VPN (Mbps)")
	flag.Int64Var(&latency, "latency", 40, "Minimum latency of VPN (ms)")
	flag.Parse()

	if proto != "any" && proto != "udp" && proto != "tcp" {
		flag.Usage()
	}

	if _, err := os.Stat(openvpnBin); err != nil {
		panic(err)
	}
}

func main() {
	var ping time.Duration
	var speed int
	var vpnLowestLatency *vpngate.VPN

	c := &http.Client{}
	ping = time.Millisecond * time.Duration(latency)
	speed = minSpeed * 1000000 // Mbps

	vpns, err := vpngate.Get(c)
	if err != nil {
		panic(err)
	}

	fmt.Println("looking for low latent endpoints...")
	for _, vpn := range vpns {
		if proto == "any" {
			if vpn.Ping <= ping && vpn.Speed > speed {
				if ok := rawconnect(vpn.IP, strconv.Itoa(vpn.Port), vpn.Proto); ok {
					vpnLowestLatency = vpn
				}
			}
		} else {
			if vpn.Ping <= ping && vpn.Speed > speed && vpn.Proto == proto {
				if ok := rawconnect(vpn.IP, strconv.Itoa(vpn.Port), vpn.Proto); ok {
					vpnLowestLatency = vpn
				}
			}
		}

	}

	if vpnLowestLatency == nil {
		panic("no vpn found")
	}

	fmt.Printf("Found one low latent vpn: %s[%s/%v] in %s\n", vpnLowestLatency.IP, vpnLowestLatency.Proto, vpnLowestLatency.Port, vpnLowestLatency.Country)
	err = ioutil.WriteFile("/tmp/vpn.config", []byte(vpnLowestLatency.OpenVPN()), 0644)
	if err != nil {
		panic(err)
	}

	openvpn(openvpnBin)
}

func rawconnect(host string, port string, proto string) bool {

	conn, err := net.DialTimeout(proto, net.JoinHostPort(host, port), 3*time.Second)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
	}

	return true

}

func openvpn(openvpnBin string) {
	var out bytes.Buffer

	cmd := exec.Command(openvpnBin, "--config", "/tmp/vpn.config")
	cmd.Stdout = &out

	fmt.Printf("Starting OpenVPN...\n")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	fmt.Printf("OpenVPN: err: %q, output: %q\n", err, cmd.Stdout)
}
