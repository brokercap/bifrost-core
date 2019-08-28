package http_manager

import "github.com/brokercap/bifrost_core/sdk/plugin/storage"

func SetToServer(toServerKey string,pluginName string,Uri string) bool {
	storage.SetToServerInfo(toServerKey,storage.ToServer{PluginName:pluginName,ConnUri:Uri,Notes:""})
	return true
}