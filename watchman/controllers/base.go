package controllers

import "github.com/xtls/xray-core/watchman/database/connector"

var defaultConnector = connector.RegisterConnector(connector.Connectors.Default)
