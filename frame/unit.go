package frame

type Configurable interface {
	LoadConfig(name string, configPtr interface{}, modified ... chan<- *ConfigItem) error
	NotifyConfigModified(items ...ConfigItem)
	GetConfigName() string
}

type ConfigItem struct {
	Name  string
	Value interface{}
}

type CommonConfigurableUnit struct {
	CfgName   string
	CfgModChs [] chan<- *ConfigItem
}

func (x *CommonConfigurableUnit) LoadConfig(name string, configPtr interface{}, modified ... chan<- *ConfigItem) error {
	x.CfgName = name
	x.CfgModChs = modified
	// cfg, ok := configPtr.(*Config)
	// todo: deal with modified notification here...
	return nil
}

func (x *CommonConfigurableUnit) NotifyConfigModified(items ...ConfigItem) {
	for _, ch := range x.CfgModChs {
		if len(items) == 0 {
			ch <- &ConfigItem{
				Name: x.CfgName,
			}
		} else {
			for _, item := range items {
				ch <- &ConfigItem{
					// add self.CfgName in front of it to prevent out-of-range notification.
					Name:  x.CfgName + "." + item.Name,
					Value: item.Value,
				}
			}
		}
	}
}

func (x *CommonConfigurableUnit) GetConfigName() string {
	return x.CfgName
}
