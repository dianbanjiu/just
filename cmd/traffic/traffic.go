package traffic

import (
	"encoding/json"
	"fmt"
	"just/config"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var TrafficCmd = &cobra.Command{
	Use:   "usage",
	Short: "Print traffic usage",
	Long:  "Print traffic usage, contain data usage and reset time",
	Run:   printTraffic,
}

var trafficUrl string

func init() {
	TrafficCmd.Flags().StringVarP(&trafficUrl, "url", "u", "", "justmysocks traffic api address")
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
		resetPrint = fmt.Sprintf("Will reset at %d of next month", traffic.BwResetDayOfMonth)
	case day < traffic.BwResetDayOfMonth:
		resetPrint = fmt.Sprintf("Will reset at %d of this month", traffic.BwResetDayOfMonth)
	default:
		resetPrint = "Already Reset today"
	}
	_, _ = fmt.Fprintf(os.Stdout, "      All: %s\n     Used: %s\nRemaining: %s\n\n%s\n", TrafficTransToHumanRead(traffic.MonthlyBwLimitB), TrafficTransToHumanRead(traffic.BwCounterB), TrafficTransToHumanRead(traffic.MonthlyBwLimitB-traffic.BwCounterB), resetPrint)
}

func TrafficTransToHumanRead(kb int64) string {
	if kb == 0 {
		return ""
	}
	var first, last int64
	var level uint8
	for kb/1000 > 0 {
		level++
		first = kb / 1000
		last = kb - first*1000
		kb /= 1000
	}

	var tLevel = map[uint8]string{
		0: "B",
		1: "KB",
		2: "MB",
		3: "GB",
		4: "TB",
	}
	return fmt.Sprintf("%d.%d%s", first, last, tLevel[level])
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
	var client http.Client
	client.Timeout = 10 * time.Second
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trafficStatus trafficStatus
	err = json.NewDecoder(resp.Body).Decode(&trafficStatus)
	return &trafficStatus, err
}
