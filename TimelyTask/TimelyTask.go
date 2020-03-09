package main

import (
	"bytes"
	"os/exec"
	"strings"

	. "github.com/xaxys/oasis/api"
)

var PLUGIN UserPlugin = instance

var instance *TimelyTaskPlugin = &TimelyTaskPlugin{
	PluginBase: PluginBase{
		PluginDescription: PluginDescription{
			Name:    "TimelyTask",
			Author:  "xaxys",
			Version: "0.1.0",
			DefaultConfigFields: map[string]interface{}{
				"version": "0.1.0",
			},
		},
	},
}

type TimelyTaskPlugin struct {
	PluginBase
}

func (p *TimelyTaskPlugin) OnEnable() bool {
	p.GetConfig().AddHandle(p.UpdateTask)
	p.GetLogger().Info("Config change handle registed.")
	p.UpdateTask()
	return true
}

type Task struct {
	Dir      string   `json:"Dir" yaml:"Dir"`
	Time     string   `json:"Time" yaml:"Time"`
	Print    bool     `json:"Print" yaml:"Print"`
	Commands []string `json:"Commands" yaml:"Commands"`
}

type FuncRunnable func()

func (f FuncRunnable) Run() { f() }

func (p *TimelyTaskPlugin) UpdateTask() {
	p.GetLogger().Info("Updating task...")
	p.GetServer().UnregisterPluginTask(p.GetPlugin())
	var taskList []Task
	if err := p.GetConfig().UnmarshalKey("Tasks", &taskList); err != nil {
		p.GetLogger().Errorf("Fail to unmarshal config! Check your config! Details: %v", err)
		return
	}
	for _, task := range taskList {
		commands := strings.Join(task.Commands, ";")
		p.GetServer().RegisterTask(p.GetPlugin(), task.Time, FuncRunnable(func() {
			cmd := exec.Command("/bin/sh", "-c", commands)
			cmd.Dir = task.Dir
			stdout, stderr, err := execCommand(cmd)
			if task.Print {
				p.GetLogger().Info(stdout)
			}
			if stderr != "" {
				p.GetLogger().Warn(stderr)
			}
			if err != nil {
				p.GetLogger().Error(err)
			}
		}))
	}
	p.GetLogger().Infof("Task updated, applied %d tasks.", len(taskList))
}

func execCommand(cmd *exec.Cmd) (string, string, error) {
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	cmd.Start()
	err := cmd.Wait()
	return outBuffer.String(), errBuffer.String(), err
}
