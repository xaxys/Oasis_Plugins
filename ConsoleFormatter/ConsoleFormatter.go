package main

import (
	"encoding/json"
	"strings"
	"time"
	"unsafe"

	"github.com/gookit/color"
	. "github.com/xaxys/oasis/api"
)

var PLUGIN UserPlugin = instance

var instance *ConsoleFormatterPlugin = &ConsoleFormatterPlugin{
	PluginBase: PluginBase{
		PluginDescription: PluginDescription{
			Name:    "ConsoleFormatter",
			Author:  "xaxys",
			Version: "0.1.2",
			DefaultConfigFields: map[string]interface{}{
				"version":    "0.1.2",
				"TimeFormat": "2006/01/02 15:04:05.000",
				"LogFormat":  "{{[%Time%]}}{{[%Level%]}}{{[%Linenum%]}}: {{[%Plugin%] }}{{%Msg%}}",
			},
		},
	},
}

type ConsoleFormatterPlugin struct {
	PluginBase
}

func (p *ConsoleFormatterPlugin) OnLoad() bool {
	p.GetServer().RegisterFormatter(instance)
	p.GetLogger().Info("ConsoleFormatter registed!")
	return true
}

func (p *ConsoleFormatterPlugin) OnEnable() bool {
	return true
}

func (p *ConsoleFormatterPlugin) OnDisable() bool {
	return true
}

const ISO8601TimeFormat = "2006-01-02T15:04:05.000Z0700"

type LogEntry struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
	Plugin  string `json:"plugin"`
	Linenum string `json:"linenum"`
}

func (p *ConsoleFormatterPlugin) Format(msg string) string {
	var l LogEntry
	var fs string
	err := json.Unmarshal(*(*[]byte)(unsafe.Pointer(&msg)), &l)
	if err != nil {
		return msg
	}

	s := p.GetConfig().GetString("LogFormat")
	if l.Time != "" {
		t, _ := time.Parse(ISO8601TimeFormat, l.Time)
		l.Time = t.Format(p.GetConfig().GetString("TimeFormat"))
	}
	if l.Level != "" {
		switch l.Level {
		case "DEBUG":
			l.Level = color.Gray.Render(l.Level)
		case "INFO":
			l.Level = color.White.Render(l.Level)
		case "WARN":
			l.Level = color.Yellow.Render(l.Level)
		case "ERROR":
			l.Level = color.Red.Render(l.Level)
		}
	}
	lastp := 0
	for i, j := nextPositon(s, 1); i != -1 && j != -1; i, j = nextPositon(s, j) {
		fs += s[lastp:i]
		f := false
		field := s[i+2 : j-2]
		if l.Level != "" && strings.Contains(field, "%Level%") {
			field = strings.ReplaceAll(field, "%Level%", l.Level)
			f = true
		}
		if l.Time != "" && strings.Contains(field, "%Time%") {
			field = strings.ReplaceAll(field, "%Time%", l.Time)
			f = true
		}
		if l.Msg != "" && strings.Contains(field, "%Msg%") {
			field = strings.ReplaceAll(field, "%Msg%", l.Msg)
			f = true
		}
		if l.Plugin != "" && strings.Contains(field, "%Plugin%") {
			field = strings.ReplaceAll(field, "%Plugin%", l.Plugin)
			f = true
		}
		if l.Linenum != "" && strings.Contains(field, "%Linenum%") {
			field = strings.ReplaceAll(field, "%Linenum%", l.Linenum)
			f = true
		}
		if f {
			fs += field
		}
		lastp = j
	}
	if lastp < len(s) {
		fs += s[lastp:]
	}
	fs = fs + "\n"

	return fs
}

func nextPositon(s string, start int) (i int, j int) {
	for i = start; i < len(s); i++ {
		if s[i-1] == '{' && s[i] == '{' {
			break
		}
	}
	for j = i + 1; j < len(s); j++ {
		if s[j-1] == '}' && s[j] == '}' {
			break
		}
	}
	if i == len(s) {
		i = 0
	}
	if j == len(s) {
		j = -2
	}
	return i - 1, j + 1
}
