package main

import (
	"net/http"
	"net"
	"os/exec"
	"os"
	"bytes"
	"log"
	"github.com/yelinaung/wifi-name"
	"strings"
	"strconv"
	"fmt"
)

var Listen string

type InterfaceAddress struct {
	name string
	addr []string
}

/**
 * create an dictonary with interface name and all ip addresses
 */
func getIpAdresses() []InterfaceAddress {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic (nil)
	}
	var ip []InterfaceAddress

	for _, iface := range ifaces {
		var inter InterfaceAddress
		var add []string

		flag := iface.Flags
		if flag&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		if flag&net.FlagUp == 0 {
			continue // interface is down
		}

		inter.name = iface.Name
		addrs, err := iface.Addrs()
		if err != nil {
			// error but we can continue please check the log file
			log.Print(err)
			continue
		}
		for _, addr := range addrs{
			var ip net.IP
			switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue // no ip address or loopback address
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			add = append(add, ip.String())

		}
		inter.addr = add
		ip = append(ip, inter)
	}
	return ip
}

/**
 * return the state of a defined service
 * @param name is the name of the service
 * @param logfile is the path to the logfile
 */
func stateOfService(name string) int {
	var ret int
	var out bytes.Buffer
	cmd := exec.Command("/bin/systemctl", "status" , name)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Start()
	if err != nil {
		log.Print(err)
		return -1
	}
	err = cmd.Wait()
	if err != nil {
		log.Print(err)
		return -1
	}

	str := strings.Split(out.String(), "\n");
	if len(str) <= 1 {
		// service not found or execution failure
		log.Print("service not found or execution failure")
		return -1
	}
	for _, element := range str {
		if strings.Contains(element, "Active: ") {
			if strings.Contains(element, "active (") {
				ret = 1
			} else if strings.Contains(element, "inactive (") {
				ret = 0
			} else {
				ret = -1
			}
		}
	}

	return ret
}

/**
 * answer the request 
 */
func status(w http.ResponseWriter, r *http.Request){
	ip := getIpAdresses()
	var message string
	message = "# TYPE status_ip_address gauge\n"
	for _, element := range ip {
		for _, el := range element.addr {
			message += "status_ip_address{interface=\"" + element.name + "\", ipAddress=\"" + el + "\"} 1\n"
		}
	}

	message += "# TYPE status_wifi_ssid gauge\n"
	message += "status_wifi_ssid{ssid=\"" + wifiname.WifiName() + "\"} 1\n"

	message += "# TYPE status_service gauge\n"
	message += "status_service{name=\"sshd\"} " + strconv.Itoa(stateOfService("sshd")) + "\n";
	message += "status_service{name=\"prometheus\"} " + strconv.Itoa(stateOfService("prometheus")) + "\n";
	message += "status_service{name=\"grafana-server\"} " + strconv.Itoa(stateOfService("grafana-server")) + "\n";
	message += "status_service{name=\"pixiecore\"} " + strconv.Itoa(stateOfService("pixiecore")) + "\n";
	message += "status_service{name=\"waitron\"} " + strconv.Itoa(stateOfService("waitron")) + "\n";

	w.Write([]byte(message))
}

/**
 * print the help message
 */
func help() {
	fmt.Println("Help dialog:\n")
	fmt.Println("Options:")
	fmt.Println("\thelp: Print this help dialog")
	fmt.Println("\tlisten: listen on this address")
	fmt.Println("Usage:")
	fmt.Println("\tstatus_exporter listen=localhost:9999")
	fmt.Println("\tstatus_exporter help")
}

func main() {
	args := os.Args[1:]
	if len(args) > 2 || len(args) == 0 {
		help()
		os.Exit(1)
	}

	for i := 0; i < len(args); i++ {
		if strings.EqualFold(args[i], "help") {
			help()
			os.Exit(0)
		} else if strings.HasPrefix(args[i], "listen="){
			com := strings.Split(args[i], "=")
			Listen = com[1]
		}
	}

	http.HandleFunc("/", status)
	err := http.ListenAndServe(Listen, nil)
	if err != nil {
		panic(err)
	}
}
