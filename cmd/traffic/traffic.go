package traffic

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"just/config"
	"net/http"
	"os"
	"time"
)

var TrafficCmd = &cobra.Command{
	Use:   "usage",
	Short: "Print traffic usage",
	Long:  "Print traffic usage, contain data usage and reset time",
	Run:   printTraffic,
}

var trafficUrl string

func init() {
	TrafficCmd.Flags().StringVarP(&trafficUrl, "traffic_url", "u", "", "copy justmysocks traffic api to here")
}
func printTraffic(cmd *cobra.Command, args []string) {
	if trafficUrl == "" {
		trafficUrl = config.CFG.GetString("traffic_url")
	}
	if trafficUrl == "" {
		_, _ = fmt.Fprintln(os.Stderr, "traffic api url cannot be empty! ")
		return
	}
	traffic, err := getTraffic(trafficUrl)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	if traffic == nil {
		_, _ = fmt.Fprintln(os.Stderr, "request success, but get nothing. Please check API or network is correctly.")
		return
	}

	l, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(l)
	day := now.Day()
	var resetPrint string
	switch {
	case day > traffic.BwResetDayOfMonth:
		resetPrint = fmt.Sprintf("流量将在下个月的%d号恢复", traffic.BwResetDayOfMonth)
	case day < traffic.BwResetDayOfMonth:
		resetPrint = fmt.Sprintf("流量将在本月的%d号恢复", traffic.BwResetDayOfMonth)
	default:
		resetPrint = "流量已于今日恢复"
	}
	_, _ = fmt.Fprintf(os.Stdout, "总流量：%dGB\n已使用流量：%dMB\n剩余流量：%dMB\n%s\n",
		traffic.MonthlyBwLimitB/1000/1000/1000,
		traffic.BwCounterB/1000/1000,
		(traffic.MonthlyBwLimitB-traffic.BwCounterB)/1000/1000,
		resetPrint)
}

type trafficStatus struct {
	MonthlyBwLimitB   int64 `json:"monthly_bw_limit_b"`
	BwCounterB        int64 `json:"bw_counter_b"`
	BwResetDayOfMonth int   `json:"bw_reset_day_of_month"`
}

// getTraffic get traffic info
// u is api for justmysocks support
// you can visit justmysocks website and check your subscription service,
// you will find an API tag, copy the path to config.yaml
func getTraffic(u string) (*trafficStatus, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trafficStatus trafficStatus
	err = json.NewDecoder(resp.Body).Decode(&trafficStatus)
	return &trafficStatus, err
}
