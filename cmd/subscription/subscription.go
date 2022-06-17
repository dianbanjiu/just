package subscription

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"just/cmd/socks"
	"just/cmd/vmess"
	"just/config"
	"just/utils"
	"net/http"
	"os"
	"strings"
)

var SubCmd = &cobra.Command{
	Use:   "sub",
	Short: "get subscription details",
	Long:  "get subscription details",
	Run:   handleSubscription,
}

func init() {
	var printVar bool
	var path string
	var writeVar bool
	SubCmd.Flags().BoolVarP(&printVar, "print", "p", false, "only print subscription detail to terminal. this is default flag")
	SubCmd.Flags().StringVarP(&path, "config", "c", "", "the clash config path")
	SubCmd.Flags().BoolVarP(&writeVar, "write", "w", false, "whether write the new subscription to clash config file")
}

func handleSubscription(cmd *cobra.Command, args []string) {
	isNeedWrite := cmd.Flags().Changed("write")
	if isNeedWrite {
		clashConfigPath, err := cmd.Flags().GetString("config")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if clashConfigPath == "" {
			clashConfigPath = config.CFG.GetString("clash_config_path")
		}
		if clashConfigPath == "" {
			fmt.Fprintln(os.Stderr, "please set the clash config path via -c flag or put it in just's config.yaml")
			return
		}
		writeConfigToFile(clashConfigPath)
	} else {
		printSubscription()
	}
}
func printSubscription() {
	u := config.CFG.GetString("subscription_url")
	hostList, _, err := getSubscription(u)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	var temp = make([]map[string]interface{}, len(hostList))
	raw, err := json.Marshal(hostList)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err = json.Unmarshal(raw, &temp); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for i, host := range temp {
		fmt.Fprintln(os.Stdout, fmt.Sprintf("Server %d:", i))
		for name, value := range host {
			fmt.Fprintln(os.Stdout, fmt.Sprintf("    %s: %v", name, value))
		}

		fmt.Fprintln(os.Stdout)
	}
}

func initClashConfig(configPath string) (*viper.Viper, error) {
	clashConfig := viper.New()
	clashConfig.SetConfigFile(configPath)
	err := clashConfig.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return clashConfig, nil
}

type ProxyGroup struct {
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type" yaml:"type"`
	Proxies  []string `json:"proxies" yaml:"proxies"`
	Url      string   `json:"url" yaml:"url"`
	Interval int      `json:"interval" yaml:"interval"`
}

func writeConfigToFile(configPath string) {
	clashConfig, err := initClashConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	u := config.CFG.GetString("subscription_url")
	hostList, nameList, err := getSubscription(u)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	clashConfig.Set("proxies", hostList)
	var proxyGroups = make([]ProxyGroup, 0)
	err = clashConfig.UnmarshalKey("proxy-groups", &proxyGroups)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for i := 0; i < len(proxyGroups); i++ {
		proxyGroups[i].Proxies = nameList
	}
	clashConfig.Set("proxy-groups", proxyGroups)
	err = clashConfig.WriteConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

}

// GetSubscription 获取订阅信息
func getSubscription(u string) ([]interface{}, []string, error) {
	// 1. 从地址获取数据
	sub, err := getRawSubscriptionFromUrl(u)
	if err != nil {
		return nil, nil, err
	}

	// 2. base64 解密
	rawHostInfo, err := utils.Base64Decode(sub)
	if err != nil {
		return nil, nil, err
	}

	rawHostList := bytes.Split(rawHostInfo, []byte("\n"))
	var hostList = make([]interface{}, 0)
	var nameList = make([]string, 0)
	for _, rawHost := range rawHostList {
		protocol := strings.ToLower(getProtocolFromRawUrl(rawHost))
		var temp interface{}
		switch protocol {
		case "ss":
			temp, err = socks.ParseRawSocks(string(rawHost[5:]))
			if err != nil {

				fmt.Fprintf(os.Stderr, "parse raw ss failed, err is %s\n", err.Error())
				continue
			}

			nameList = append(nameList, temp.(socks.SocksInfo).Name)
		case "vmess":
			temp, err = vmess.ParseRawVmess(string(rawHost[8:]))
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse raw vmess failed, err is %s\n", err.Error())
				continue
			}
			nameList = append(nameList, temp.(vmess.VmessInfo).Name)
		}

		hostList = append(hostList, temp)
	}
	return hostList, nameList, nil
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
