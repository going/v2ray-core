package controllers

import "v2ray.com/core/watchman/database/connector"

var defaultConnector = connector.RegisterConnector(connector.Connectors.Default)
