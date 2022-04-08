package collector

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/burningsunrise/tplink-exporter/model"
	"github.com/burningsunrise/tplink-exporter/parser"

	"github.com/panjf2000/ants"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type tplinkCollector struct {
	txPackets          *prometheus.Desc
	rxPackets          *prometheus.Desc
	speed              *prometheus.Desc
	vlans              *prometheus.Desc
	memory             *prometheus.Desc
	cpu                *prometheus.Desc
	rxBadPackets       *prometheus.Desc
	txBadPackets       *prometheus.Desc
	broadcastRxPackets *prometheus.Desc
	broadcastTxPackets *prometheus.Desc
	multicastTxPackets *prometheus.Desc
	multicastRxPackets *prometheus.Desc
	unicastTxPackets   *prometheus.Desc
	unicastRxPackets   *prometheus.Desc
	generalInfo        *prometheus.Desc
}

func NewTplinkCollector() *tplinkCollector {
	return &tplinkCollector{
		txPackets: prometheus.NewDesc("port_tx_metric",
			"Shows tx packets on the hosts port",
			[]string{"portnum", "host"}, nil,
		),
		rxPackets: prometheus.NewDesc("port_rx_metric",
			"Shows rx packets on the hosts port",
			[]string{"portnum", "host"}, nil,
		),
		speed: prometheus.NewDesc("port_speed_metric",
			"Shows the hosts port speed",
			[]string{"portnum", "host"}, nil,
		),
		vlans: prometheus.NewDesc("port_vlans_metric",
			"Shows the vlans on port number",
			[]string{"vlanname", "vlanid", "host", "port"}, nil),
		memory: prometheus.NewDesc("switch_memory_metric",
			"Shows the specific switch memory",
			[]string{"host", "macaddress"}, nil),
		cpu: prometheus.NewDesc("switch_cpu_metric",
			"Shows the specific switch cpu",
			[]string{"host", "macaddress"}, nil),
		rxBadPackets: prometheus.NewDesc("port_badrx_metric",
			"Shows bad rx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		txBadPackets: prometheus.NewDesc("port_badtx_metric",
			"Shows bad tx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		broadcastRxPackets: prometheus.NewDesc("port_broadcastrx_metric",
			"Shows broadcast rx packets the hosts port",
			[]string{"portnum", "host"}, nil),
		broadcastTxPackets: prometheus.NewDesc("port_broadcasttx_metric",
			"Shows broadcast tx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		multicastTxPackets: prometheus.NewDesc("port_multicasttx_metric",
			"Shows multicast tx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		multicastRxPackets: prometheus.NewDesc("port_multicastrx_metric",
			"Shows multicast rx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		unicastTxPackets: prometheus.NewDesc("port_unicasttx_metric",
			"Shows unicast tx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		unicastRxPackets: prometheus.NewDesc("port_unicastrx_metric",
			"Shows unicast rx packets on the hosts port",
			[]string{"portnum", "host"}, nil),
		generalInfo: prometheus.NewDesc("switch_generalinfo_metric",
			"Shows general information about the switch with temperature as a metric",
			[]string{"devloc", "sysdesc", "host", "hwversion", "fmversion", "macaddress",
				"systime", "runtime", "serialnum"}, nil),
	}
}


func (collector *tplinkCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.txPackets
	ch <- collector.rxPackets
	ch <- collector.speed
	ch <- collector.vlans
	ch <- collector.memory
	ch <- collector.cpu
	ch <- collector.rxBadPackets
	ch <- collector.txBadPackets
	ch <- collector.broadcastRxPackets
	ch <- collector.broadcastTxPackets
	ch <- collector.unicastRxPackets
	ch <- collector.unicastTxPackets
	ch <- collector.multicastRxPackets
	ch <- collector.multicastTxPackets
	ch <- collector.generalInfo
}

func (collector *tplinkCollector) Collect(ch chan<- prometheus.Metric) {
	collection := probeDevices()

	for _, c := range collection {
		ch <- prometheus.MustNewConstMetric(collector.memory, prometheus.GaugeValue, float64(c.Data.Memory[0]), c.DnsName,
			c.Data.MacAddress)
		ch <- prometheus.MustNewConstMetric(collector.cpu, prometheus.GaugeValue, float64(c.Data.Cpu[0]), c.DnsName,
			c.Data.MacAddress)
		ch <- prometheus.MustNewConstMetric(collector.generalInfo, prometheus.GaugeValue, float64(c.Data.Temperature),
			c.Data.DevLoc, c.Data.SysDescription, c.DnsName, c.Data.HwVersion, c.Data.FwVersion, c.Data.MacAddress,
			c.Data.SysTime, c.Data.RunTime, c.Data.SeNumber)
		for _, p := range c.Ports {
			port := strings.Split(p.Port, "/")[2]
			var vlanName []string
			var vlanId []float64
			if num, err := strconv.ParseFloat(port, 64); err == nil {
				ch <- prometheus.MustNewConstMetric(collector.rxPackets, prometheus.GaugeValue, float64(p.PktsRx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.txPackets, prometheus.GaugeValue, float64(p.PktsTx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.speed, prometheus.GaugeValue, float64(p.SpeedLink),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.rxBadPackets, prometheus.GaugeValue, float64(p.ErrorsRx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.txBadPackets, prometheus.GaugeValue, float64(p.ErrorsTx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.broadcastRxPackets, prometheus.GaugeValue, float64(p.BroadcastRx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.broadcastTxPackets, prometheus.GaugeValue, float64(p.BroadcastTx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.unicastRxPackets, prometheus.GaugeValue, float64(p.UnicastRx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.unicastTxPackets, prometheus.GaugeValue, float64(p.UnicastTx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.multicastRxPackets, prometheus.GaugeValue, float64(p.MulticastRx),
					port, c.DnsName)
				ch <- prometheus.MustNewConstMetric(collector.multicastTxPackets, prometheus.GaugeValue, float64(p.MulticastTx),
					port, c.DnsName)
				// Vlans
				for _, vl := range p.Vlans {
					vlanName = append(vlanName, vl.Name)
					vlanId = append(vlanId, vl.VlanID)
				}
				ch <- prometheus.MustNewConstMetric(collector.vlans, prometheus.GaugeValue, num, strings.Join(vlanName, ","),
					strings.Trim(strings.Replace(fmt.Sprint(vlanId), " ", ",", -1), "[]"), c.DnsName, port)
			}
		}
	}
}

func probeDevices() []model.Tplink {
	y := parser.YamlConfig{}
	y.GetConfig()

	defer ants.Release()
	var wg sync.WaitGroup
	log.WithFields(log.Fields{
		"status": "probing",
	}).Info("scanning all devices")
	collection := []model.Tplink{}
	client := model.HttpClient()

	p, _ := ants.NewPoolWithFunc(20, func(i interface{}) {
		defer wg.Done()

		tplink := model.Tplink{DnsName: i.(string)}

		if err := tplink.Login(y, client); err != nil {
			log.WithFields(log.Fields{
				"login": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchSystem(client); err != nil {
			log.WithFields(log.Fields{
				"switchsystem": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchPorts(client); err != nil {
			log.WithFields(log.Fields{
				"switchports": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchPortStatistics(client); err != nil {
			log.WithFields(log.Fields{
				"portstats": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchPortVlans(client); err != nil {
			log.WithFields(log.Fields{
				"portvlans": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchPortVlanCfg(client); err != nil {
			log.WithFields(log.Fields{
				"portvlancfg": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchMacVlanCfgModel(client); err != nil {
			log.WithFields(log.Fields{
				"macvlancfg": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchMemory(client); err != nil {
			log.WithFields(log.Fields{
				"memory": i.(string),
			}).Error(err)
			return
		}
		if err := tplink.SwitchCpu(client); err != nil {
			log.WithFields(log.Fields{
				"cpu": i.(string),
			}).Error(err)
			return
		}
		collection = append(collection, tplink)
	})

	defer p.Release()
	for _, device := range y.Devices {
		wg.Add(1)
		_ = p.Invoke(device)
	}

	wg.Wait()

	log.WithFields(log.Fields{
		"devices": len(y.Devices),
		"status":  "finished",
	}).Info("waiting for next iteration")
	return collection
}
