package settings

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/kardianos/osext"
)

type pair struct {
	Left  string `toml:"left"`
	Right string `toml:"right"`
}

type staticCfg struct {
	VirtualRoot string `toml:"vstatic"`
	LocalRoot   string `toml:"lstatic"`
	CompressDef string `toml:"compress"`
}

type serverCfg struct {
	Port string `toml:"port"`
	Host string `toml:"host"`
}
type fileSyncCfg struct {
	Pairs []pair `toml:"pair"`
}

type templateCfg struct {
	Home         string `toml:"home"`
	DelimesLeft  string `toml:"ldelime"`
	DelimesRight string `toml:"rdelime"`
	Charset      string `toml:"charset"`
	Reload       bool   `toml:"reload"`
}
type defaultVar struct {
	AppName string `toml:"appname"`
}

type adminCfg struct {
	Passwd string `toml:"passwd"`
}

type logCfg struct {
	Path   string `toml:"path"`
	Format string `toml:"format"`
	File   string `toml:"-"`
}

type setting struct {
	Static      staticCfg   `toml:"static"`
	Server      serverCfg   `toml:"server"`
	Template    templateCfg `toml:"template"`
	DefaultVars defaultVar  `toml:"defaultvars"`
	Admin       adminCfg    `toml:"admin"`
	Log         logCfg      `toml:"log"`
	Filesync    fileSyncCfg `toml:"filesync"`
}

var (
	Folder        string
	settingStruct = new(setting)

	//GlobalSettings
	Static      staticCfg
	Server      serverCfg
	Filesync    fileSyncCfg
	Template    templateCfg
	DefaultVars defaultVar
	Admin       adminCfg
	Log         logCfg
)

var lock = new(sync.Mutex)

func init() {
	var err error
	Folder, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
}

func Init() {
	cfgFile := getAbs("./settings/settings.toml") // "D:\\Dev\\gopath\\src\\github.com\\Felamande\\filesync.v2\\settings\\settings.toml"
	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}
	toml.Unmarshal(b, settingStruct)

	Static = settingStruct.Static
	Server = settingStruct.Server
	Filesync = settingStruct.Filesync
	Template = settingStruct.Template
	DefaultVars = settingStruct.DefaultVars
	Admin = settingStruct.Admin
	Log = settingStruct.Log
}

func getAbs(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(Folder, path)
	}
	return path
}

// func readConfig(ConfigFile string) *SavedConfig {
// 	fmt.Println(ConfigFile)
// 	config := new(SavedConfig)
// 	data, err := ioutil.ReadFile(ConfigFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = ymlread.Unmarshal(data, config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return config

// }
