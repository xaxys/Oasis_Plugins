# Oasis_Plugins

This is plugin repository for Oasis.

Maintained by xaxys.

You can find the Oasis plugin server there https://github.com/xaxys/Oasis

# Plugins

### ConsoleFormatter

A formatter for console prints.

It allows you to custom the print format and time format.

##### screenshot

![img.jpg](https://i.loli.net/2020/03/09/aqzcb7BtPxVSLWK.jpg)

##### config

`{{some letters and variable(such as %Msg%)}}`: use`{{` `}}`to warp your variable so the contents won't exist if the variable doesn't exist. 

```yaml
logformat: '{{[%Time%]}}{{[%Level%]}}{{[%Linenum%]}}: {{[%Plugin%] }}{{%Msg%}}'
timeformat: 2006/01/02 15:04:05.000
version: 0.1.2
```

### TimelyTask

A plugin for you to easily set timely task.

It will auto add / remove tasks when you modified the config

##### screenshot

![img.jpg](https://i.loli.net/2020/03/09/uJlK5CS4np1jOqh.jpg)

##### config

The example contains a task to execute `ls` and `echo helloworld` in the 5th second every minutes.

```yaml
version: 0.1.0
Tasks:
- Dir:
  Time: '5 * * * * *'
  Print: true
  Commands:
  - 'ls'
  - 'echo helloworld'
```

