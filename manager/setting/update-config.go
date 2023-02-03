package setting

import (
	scyna_engine "github.com/scyna/core/engine"
)

func UpdateDefaultConfig(config *scyna_engine.Configuration) {
	// log.Printf("Update config: %+v\n", config)
	// manager.DefaultConfig = config

	// val, _ := json.Marshal(config)
	// var request = scyna.WriteSettingRequest{
	// 	Module: manager.MODULE_CODE,
	// 	Key:    "config",
	// 	Value:  string(val),
	// }
	// var response scyna.Error
	// if err := scyna.CallService(scyna.SETTING_WRITE_URL, &request, &response); err.Code != scyna.OK.Code {
	// 	log.Printf("Update config error: %+v\n", &response)
	// }
}
