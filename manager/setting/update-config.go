package setting

import (
	scyna_proto "github.com/scyna/core/proto/generated"
)

func UpdateDefaultConfig(config *scyna_proto.Configuration) {
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
