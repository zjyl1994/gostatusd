package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
)

const backgroundStatDuration = 3 * time.Second

var excludeInterfaceNamePrefix = []string{"lo", "tun", "docker", "veth", "br-", "vmbr", "vnet", "kube"}

var backgroundStat struct {
	DiskRead       uint64
	DiskReadSpeed  uint64
	DiskWrite      uint64
	DiskWriteSpeed uint64

	NetRx       uint64
	NetTx       uint64
	NetTotalIn  uint64
	NetTotalOut uint64

	NetMonthIn  uint64
	NetMonthOut uint64
	NetMonth    string

	CPUPercent float64

	Timestamp int64
}

type recordSaveFormat struct {
	In    uint64 `json:"in"`
	Out   uint64 `json:"out"`
	Month string `json:"month"`
}

func initWorker() error {
	// init network
	in, out, err := getNetInOut()
	if err != nil {
		return err
	}
	backgroundStat.NetTotalIn = in
	backgroundStat.NetTotalOut = out
	// init monthly record
	save, err := ioutil.ReadFile(recordFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			emptySave, err := json.Marshal(recordSaveFormat{Month: time.Now().Format("2006-01")})
			if err != nil {
				return err
			}
			return ioutil.WriteFile(recordFile, emptySave, 0644)
		}
		return err
	}
	var record recordSaveFormat
	err = json.Unmarshal(save, &record)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	backgroundStat.NetMonthIn = record.In
	backgroundStat.NetMonthOut = record.Out
	backgroundStat.NetMonth = record.Month
	// set time
	backgroundStat.Timestamp = time.Now().Unix()
	return nil
}

func backgroundWorker() {
	for {
		// calcuate cpu percent
		if percent, err := cpu.Percent(backgroundStatDuration, false); err == nil {
			backgroundStat.CPUPercent = percent[0]
		}
		now := time.Now().Unix()
		duration := uint64(now - backgroundStat.Timestamp)
		currentMonth := time.Now().Format("2006-01")
		if duration == 0 {
			continue
		}
		// disk speed
		counters, err := disk.IOCounters()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		var readBytes, writeBytes uint64
		for _, c := range counters {
			readBytes += c.ReadBytes
			writeBytes += c.WriteBytes
		}

		backgroundStat.DiskReadSpeed = (readBytes - backgroundStat.DiskRead) / duration
		backgroundStat.DiskWriteSpeed = (writeBytes - backgroundStat.DiskWrite) / duration

		backgroundStat.DiskRead = readBytes
		backgroundStat.DiskWrite = writeBytes

		// network speed / traffic monthly
		in, out, err := getNetInOut()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		diffIn := in - backgroundStat.NetTotalIn
		diffOut := out - backgroundStat.NetTotalOut
		backgroundStat.NetTotalIn = in
		backgroundStat.NetTotalOut = out

		backgroundStat.NetRx = diffIn / duration
		backgroundStat.NetTx = diffOut / duration

		if currentMonth == backgroundStat.NetMonth {
			backgroundStat.NetMonthIn += diffIn
			backgroundStat.NetMonthOut += diffOut
		} else {
			backgroundStat.NetMonthIn = diffIn
			backgroundStat.NetMonthOut = diffOut
			backgroundStat.NetMonth = currentMonth
		}

		save, _ := json.Marshal(recordSaveFormat{
			In:    backgroundStat.NetMonthIn,
			Out:   backgroundStat.NetMonthOut,
			Month: backgroundStat.NetMonth,
		})

		err = ioutil.WriteFile(recordFile, save, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		backgroundStat.Timestamp = now
	}
}

func getNetInOut() (netIn uint64, netOut uint64, err error) {
	nv, err := net.IOCounters(true)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range nv {
		if matchPrefix(v.Name, excludeInterfaceNamePrefix) {
			continue
		}
		netIn += v.BytesRecv
		netOut += v.BytesSent
	}
	return netIn, netOut, nil
}

func matchPrefix(s string, prefixs []string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
