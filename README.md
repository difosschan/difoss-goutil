# difoss-goutil

**difoss-goutil**, is a utility tools that help you easier to build applications in go.

It contain following feature:

- configure file

  Re-encapsulate [viper](https://github.com/spf13/viper), sub-feature as follow:
  
  - [x] Perception of changes from RAM to file: use CommonConfigurableUnit.NotifyConfigModified to notify it.
  - [ ] Perception of changes from file to RAM (TODO).

- logger

  Re-encapsulate [zap](https://github.com/uber-go/zap) and [lumberjack](https://github.com/natefinch/lumberjack)
