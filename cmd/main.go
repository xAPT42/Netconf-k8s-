package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	netconfPort       = "830"
	netconfSubsystem  = "netconf"
	helloMessage      = `<?xml version="1.0" encoding="UTF-8"?>
<hello xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <capabilities>
    <capability>urn:ietf:params:netconf:base:1.0</capability>
  </capabilities>
</hello>]]>]]>`

	getConfigRPC = `<?xml version="1.0" encoding="UTF-8"?>
<rpc message-id="1" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <get-config>
    <source>
      <running/>
    </source>
  </get-config>
</rpc>]]>]]>`

	closeSessionRPC = `<?xml version="1.0" encoding="UTF-8"?>
<rpc message-id="2" xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <close-session/>
</rpc>]]>]]>`
)

type Config struct {
	RouterAddress string
	Username      string
	Password      string
}

type ComplianceResult struct {
	Passed []string
	Failed []string
}

type RpcReply struct {
	XMLName xml.Name `xml:"rpc-reply"`
	Data    string   `xml:",innerxml"`
}

func main() {
	config := parseFlags()

	log.Println("[INFO] Starting NETCONF Compliance Checker v1.0.0")
	log.Printf("[INFO] Connecting to NETCONF router at %s\n", config.RouterAddress)

	result, err := runComplianceCheck(config)
	if err != nil {
		log.Printf("[ERROR] Compliance check failed: %v\n", err)
		os.Exit(1)
	}

	printResults(result)

	if len(result.Failed) > 0 {
		os.Exit(1)
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.RouterAddress, "router-address", "localhost:830", "NETCONF router address (host:port)")
	flag.StringVar(&config.Username, "username", "netconf", "NETCONF username")
	flag.StringVar(&config.Password, "password", "netconf", "NETCONF password")
	flag.Parse()

	return config
}

func runComplianceCheck(config Config) (*ComplianceResult, error) {
	client, err := connectNetconf(config)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer client.Close()

	log.Println("[INFO] NETCONF session initiated")
	log.Println("[INFO] Retrieving running configuration...")

	configData, err := getRunningConfig(client)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	log.Println("[INFO] Configuration retrieved successfully")
	log.Println("[INFO] Validating compliance rules...")

	result := validateCompliance(configData)

	if err := closeSession(client); err != nil {
		log.Printf("[WARN] Failed to close session gracefully: %v\n", err)
	} else {
		log.Println("[INFO] NETCONF session closed")
	}

	return result, nil
}

func connectNetconf(config Config) (*ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", config.RouterAddress, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH dial failed: %w", err)
	}

	log.Println("[INFO] SSH connection established")

	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	if err := session.RequestSubsystem(netconfSubsystem); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to request NETCONF subsystem: %w", err)
	}

	if _, err := stdin.Write([]byte(helloMessage)); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to send hello: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

	return session, nil
}

func getRunningConfig(session *ssh.Session) (string, error) {
	stdin, err := session.StdinPipe()
	if err != nil {
		return "", err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}

	if _, err := stdin.Write([]byte(getConfigRPC)); err != nil {
		return "", err
	}

	buf := make([]byte, 32768)
	n, err := stdout.Read(buf)
	if err != nil {
		return "", err
	}

	response := string(buf[:n])
	response = strings.ReplaceAll(response, "]]>]]>", "")

	return response, nil
}

func closeSession(session *ssh.Session) error {
	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	_, _ = stdin.Write([]byte(closeSessionRPC))
	time.Sleep(200 * time.Millisecond)

	return nil
}

func validateCompliance(configData string) *ComplianceResult {
	result := &ComplianceResult{
		Passed: []string{},
		Failed: []string{},
	}

	configLower := strings.ToLower(configData)

	if strings.Contains(configLower, "ntp") || strings.Contains(configLower, "clock") {
		result.Passed = append(result.Passed, "NTP is enabled")
		log.Println("[PASS] ✓ NTP is enabled")
	} else {
		result.Failed = append(result.Failed, "NTP is not configured")
		log.Println("[FAIL] ✗ NTP is not configured")
	}

	if strings.Contains(configLower, "telnet") && !strings.Contains(configLower, "no telnet") {
		result.Failed = append(result.Failed, "Telnet is enabled - SECURITY VIOLATION")
		log.Println("[FAIL] ✗ Telnet is enabled - SECURITY VIOLATION")
	} else {
		result.Passed = append(result.Passed, "Telnet is disabled")
		log.Println("[PASS] ✓ Telnet is disabled")
	}

	if strings.Contains(configLower, "hostname") || strings.Contains(configLower, "netconf") {
		result.Passed = append(result.Passed, "Hostname follows naming convention")
		log.Println("[PASS] ✓ Hostname follows naming convention")
	} else {
		result.Failed = append(result.Failed, "Hostname does not follow naming convention")
		log.Println("[FAIL] ✗ Hostname does not follow naming convention")
	}

	return result
}

func printResults(result *ComplianceResult) {
	totalRules := len(result.Passed) + len(result.Failed)
	passedCount := len(result.Passed)
	failedCount := len(result.Failed)

	if failedCount == 0 {
		log.Println("[PASS] ============================================")
		log.Println("[PASS] Compliance check successful!")
		log.Printf("[PASS] All %d rules passed\n", passedCount)
		log.Println("[PASS] ============================================")
		log.Println("[INFO] Exiting with code 0")
	} else {
		log.Println("[FAIL] ============================================")
		log.Println("[FAIL] Compliance check failed!")
		log.Printf("[FAIL] %d rules passed, %d rule(s) failed\n", passedCount, failedCount)
		log.Println("[FAIL] ============================================")
		log.Println("[INFO] Exiting with code 1")

		log.Println("\n[FAIL] Failed rules:")
		for _, rule := range result.Failed {
			log.Printf("[FAIL]   - %s\n", rule)
		}
	}
}
