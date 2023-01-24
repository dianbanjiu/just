package subscription

import (
	"bufio"
	"bytes"
	"fmt"
	"just/cmd/socks"
	"just/cmd/vmess"
	"just/config"
	"just/utils"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var SubCmd = &cobra.Command{
	Use:   "sub",
	Short: "get subscription details",
	Long:  "get subscription details",
	Run:   handleSubscription,
}

var (
	clashConfigPath string
	writeVar        bool
	subscriptionUrl string
)

func init() {
	SubCmd.Flags().StringVarP(&clashConfigPath, "config", "c", "", "the clash config path")
	SubCmd.Flags().BoolVarP(&writeVar, "write", "w", true, "whether write the new subscription to clash config file")
	SubCmd.Flags().StringVarP(&subscriptionUrl, "url", "u", "", "subscription url")
}

func handleSubscription(cmd *cobra.Command, args []string) {
	isNeedWrite := cmd.Flags().Changed("write")
	if isNeedWrite {
		if clashConfigPath == "" {
			clashConfigPath = config.CFG.GetString("clash_config_path")
		}
		if clashConfigPath == "" {
			_, _ = fmt.Fprintln(os.Stderr, "please set the clash config path via -c flag or put it in just's config.yaml")
			return
		}
		writeConfigToFile(clashConfigPath)
	}
}

type ProxyGroup struct {
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type" yaml:"type"`
	Proxies  []string `json:"proxies" yaml:"proxies"`
	Url      string   `json:"url" yaml:"url"`
	Interval int      `json:"interval" yaml:"interval"`
}

func writeConfigToFile(configPath string) {
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "clash config file doesn't exist, and create new failed. err: %v", err.Error())
		return
	}
	defer file.Close()

	if subscriptionUrl == "" {
		subscriptionUrl = config.CFG.GetString("subscription_url")
	}
	if subscriptionUrl == "" {
		_, _ = fmt.Fprintln(os.Stderr, "subscription url cannot be empty")
		return
	}
	list, err := getSubscription(subscriptionUrl)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	templ, err := template.ParseFiles("clash.yaml")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	err = templ.Execute(file, list)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}

// GetSubscription 获取订阅信息
func getSubscription(u string) ([]map[string]string, error) {
	// 1. 从地址获取数据
	sub, err := getRawSubscriptionFromUrl(u)
	if err != nil {
		return nil, err
	}

	// 2. base64 解密
	rawHostInfo, err := utils.Base64Decode(sub)
	if err != nil {
		return nil, err
	}
	rawHostList := bytes.Split(rawHostInfo, []byte("\n"))

	var hostList = make([]map[string]string, 0)
	for _, rawHost := range rawHostList {
		protocol := strings.ToLower(getProtocolFromRawUrl(rawHost))
		var temp map[string]string
		switch protocol {
		case "ss":
			temp, err = socks.ParseRawSocks(string(rawHost[5:]))
			if err != nil {

				_, _ = fmt.Fprintf(os.Stderr, "parse raw ss failed, err is %s\n", err.Error())
				continue
			}
		case "vmess":
			temp, err = vmess.ParseRawVmess(string(rawHost[8:]))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "parse raw vmess failed, err is %s\n", err.Error())
				continue
			}
		}

		hostList = append(hostList, temp)
	}
	return hostList, nil
}

// getProtocolFromRawUrl 获取代理的协议信息，如 ss、vmess
func getProtocolFromRawUrl(raw []byte) string {
	var protocol = make([]byte, 0)
	for _, b := range raw {
		if b == ':' {
			break
		}
		protocol = append(protocol, b)
	}
	return string(protocol)
}

// getRawSubscriptionFromUrl 从订阅地址获取加密之后的订阅信息
func getRawSubscriptionFromUrl(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	line, _, err := reader.ReadLine()
	return string(line), err
}
