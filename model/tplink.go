package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/burningsunrise/tplink-exporter/parser"
)

type Tplink struct {
	Data struct {
		Tid               string    `json:"_tid_"`
		UsrLvl            int       `json:"usrLvl"`
		PwdNeedChange     int       `json:"pwdNeedChange"`
		Eight02xSta       float64   `json:"_802x_sta"`
		BlVersion         string    `json:"bl_version"`
		ContactInfo       string    `json:"contact_info"`
		DevLoc            string    `json:"dev_loc"`
		DevName           string    `json:"dev_name"`
		DhcpRelaySta      float64   `json:"dhcp_relay_sta"`
		FanFlag           float64   `json:"fan_flag"`
		FanSpeed          string    `json:"fan_speed"`
		FanSta            float64   `json:"fan_sta"`
		FwVersion         string    `json:"fw_version"`
		HwVersion         string    `json:"hw_version"`
		IgmpSnoopingSta   float64   `json:"igmp_snooping_sta"`
		JumboFrameSta     float64   `json:"jumbo_frame_sta"`
		MacAddress        string    `json:"mac_address"`
		MaxTemp           float64   `json:"max_temp"`
		MldSnoopingSta    float64   `json:"mld_snooping_sta"`
		RunTime           string    `json:"run_time"`
		SeNumber          string    `json:"se_number"`
		SerialPortSetting float64   `json:"serial_port_setting"`
		SnmpSta           float64   `json:"snmp_sta"`
		SntpSta           float64   `json:"sntp_sta"`
		SpanningTreeSta   float64   `json:"spanning_tree_sta"`
		SSHSta            float64   `json:"ssh_sta"`
		SysDescription    string    `json:"sys_description"`
		SysTime           string    `json:"sys_time"`
		TelnetSta         float64   `json:"telnet_sta"`
		TemSta            float64   `json:"tem_sta"`
		Temperature       float64   `json:"temperature"`
		WebSta            float64   `json:"web_sta"`
		Memory            []float64 `json:"memory"`
		Cpu               []float64 `json:"cpu"`
	} `json:"data"`
	Ports []struct {
		DuplexCfg      float64 `json:"duplexCfg"`
		DuplexLink     float64 `json:"duplexLink"`
		FlowControl    float64 `json:"flowControl"`
		Include        float64 `json:"include"`
		Lines          float64 `json:"lines"`
		LinkStatus     float64 `json:"linkStatus"`
		MediaType      float64 `json:"mediaType"`
		Port           string  `json:"port"`
		SpeedCfg       float64 `json:"speedCfg"`
		SpeedLink      float64 `json:"speedLink"` //0 and 1 = 0m, 2 = 100m, 3 = 1000m
		State          float64 `json:"state"`
		Type           float64 `json:"type"`
		BroadcastRx    float64 `json:"broadcastRx"`
		MulticastRx    float64 `json:"multicastRx"`
		UnicastRx      float64 `json:"unicastRx"`
		BroadcastTx    float64 `json:"broadcastTx"`
		MulticastTx    float64 `json:"multicastTx"`
		UnicastTx      float64 `json:"unicastTx"`
		OversizePktsTx float64 `json:"oversizePktsTx"`
		ErrorsTx       float64 `json:"errorsTx"`
		PktsTx         float64 `json:"pktsTx"`
		BytesTx        float64 `json:"bytesTx"`
		Pkts64         float64 `json:"Pkts64"`
		Pkts65         float64 `json:"Pkts65"`
		Pkts128        float64 `json:"Pkts128"`
		Pkts256        float64 `json:"Pkts256"`
		Pkts512        float64 `json:"Pkts512"`
		Pkts1023       float64 `json:"Pkts1023"`
		UndersizePkts  float64 `json:"undersizePkts"`
		ErrorsRx       float64 `json:"errorsRx"`
		OversizePktsRx float64 `json:"oversizePktsRx"`
		PktsRx         float64 `json:"pktsRx"`
		BytesRx        float64 `json:"bytesRx"`
		Pvid           float64 `json:"pvid"`
		IngressCheck   float64 `json:"ingress_check"`
		FrameType      float64 `json:"frame_type"`
		Lag            string  `json:"lag"`
		Vlans          []struct {
			Key    float64 `json:"key"`
			Name   string  `json:"name"`
			VlanID float64 `json:"vlanId"`
		} `json:"vlans"`
		Macvlan []struct {
			Key      string  `json:"key"`
			Mac      string  `json:"mac"`
			Note     string  `json:"note"`
			VlanID   float64 `json:"vlanId"`
			VlanName string  `json:"vlanName"`
		} `json:"macvlan"`
	} `json:"ports"`
	Errorcode int  `json:"errorcode"`
	Success   bool `json:"success"`
	Timeout   bool `json:"timeout"`
	DnsName   string
}

func (t *Tplink) Login(y parser.YamlConfig, c *http.Client) error {

	url := fmt.Sprintf("https://%s/data/login.json", t.DnsName)
	payload := strings.NewReader(
		fmt.Sprintf(
			"{\"username\":\"%s\",\"password\":\"%s\",\"operation\":\"write\"}",
			y.User, y.Password),
	)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &t)
	return nil
}

func (t *Tplink) SwitchSystem(c *http.Client) error {
	url := fmt.Sprintf("https://%s/data/systemSummaryConfig.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"operation\":\"read\",\"tab\":\"unit1\"}")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &t)
	return nil
}

func (t *Tplink) SwitchPorts(c *http.Client) error {
	url := fmt.Sprintf("https://%s/data/port.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"operation\":\"load\",\"special\":\"display\",\"tab\":\"unit1\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var jbody string = strings.ReplaceAll(string(body), "data", "ports")
	json.Unmarshal([]byte(jbody), &t)
	for index, port := range t.Ports {
		switch port.SpeedLink {
		case 2:
			t.Ports[index].SpeedLink = 100
		case 3:
			t.Ports[index].SpeedLink = 1000
		case 1:
			t.Ports[index].SpeedLink = 0
		default:
			t.Ports[index].SpeedLink = 0
		}
	}
	return nil
}

func (t *Tplink) SwitchPortStatistics(c *http.Client) error {
	for index, portInfo := range t.Ports {
		var jsonMap map[string]interface{}
		url := fmt.Sprintf("https://%s/data/trafficMonitorCfgDetailModel.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
		payload := strings.NewReader(fmt.Sprintf("{\"operation\":\"read\",\"port\":\"%s\"}", portInfo.Port))
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("Content-Type", "application/json")
		res, err := c.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &jsonMap)
		dataMap := jsonMap["data"].(map[string]interface{})
		for j, ma := range dataMap {
			switch ma.(type) {
			case string:
				if s, err := strconv.ParseFloat(strings.Replace(ma.(string), ",", "", -1), 64); err == nil {
					dataMap[j] = s
				}
			}
		}

		theData, err := json.Marshal(dataMap)
		if err != nil {
			return err
		}
		json.Unmarshal(theData, &t.Ports[index])
	}
	return nil
}

func (t *Tplink) SwitchPortVlans(c *http.Client) error {
	for index, portInfo := range t.Ports {
		url := fmt.Sprintf("https://%s/data/vlanPortDetailCfg.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
		payload := strings.NewReader(fmt.Sprintf("{\"operation\":\"load\",\"port\":\"%s\"}", portInfo.Port))
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("Content-Type", "application/json")
		res, err := c.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		var jbody string = strings.ReplaceAll(string(body), "data", "vlans")
		json.Unmarshal([]byte(jbody), &t.Ports[index])
	}
	return nil
}

func (t *Tplink) SwitchPortVlanCfg(c *http.Client) error {
	var jsonMap map[string][]interface{}
	url := fmt.Sprintf("https://%s/data/vlanPortCfg.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"operation\":\"load\",\"tab\":\"unit1\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &jsonMap)
	for index, portInfo := range t.Ports {
		for _, x := range jsonMap["data"] {
			if portInfo.Port == x.(map[string]interface{})["key"] {
				t.Ports[index].Pvid = x.(map[string]interface{})["pvid"].(float64)
				t.Ports[index].Lag = x.(map[string]interface{})["lag"].(string)
				t.Ports[index].IngressCheck = x.(map[string]interface{})["ingress_check"].(float64)
				t.Ports[index].FrameType = x.(map[string]interface{})["frame_type"].(float64)
			}
		}
	}
	return nil
}

func (t *Tplink) SwitchMacVlanCfgModel(c *http.Client) error {
	var jsonMap map[string]interface{}
	url := fmt.Sprintf("https://%s/data/vlanMacCfgModel.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"operation\":\"read\",\"tab\":\"unit1\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &jsonMap)
	ports := strings.Split(jsonMap["data"].(map[string]interface{})["ports"].(string), ",")
	for index, portInfo := range t.Ports {
		for _, x := range ports {
			if portInfo.Port == x {
				t.switchMacVlanCfg(c, index)
			}
		}
	}
	return nil
}

func (t *Tplink) switchMacVlanCfg(c *http.Client, index int) error {
	url := fmt.Sprintf("https://%s/data/vlanMacCfg.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"operation\":\"load\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var jbody string = strings.ReplaceAll(string(body), "data", "macvlan")
	json.Unmarshal([]byte(jbody), &t.Ports[index])
	return nil
}

func (t *Tplink) SwitchMemory(c *http.Client) error {
	url := fmt.Sprintf("https://%s/data/memoryInfo.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"unit\":\"unit1\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &t)
	return nil
}

func (t *Tplink) SwitchCpu(c *http.Client) error {
	url := fmt.Sprintf("https://%s/data/cpuInfo.json?_tid_=%s&usrLvl=%d", t.DnsName, t.Data.Tid, t.Data.UsrLvl)
	payload := strings.NewReader("{\"unit\":\"unit1\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &t)
	return nil
}
