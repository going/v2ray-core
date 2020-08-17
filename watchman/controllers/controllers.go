package controllers

import (
	"context"
	"fmt"
	"time"

	"v2ray.com/core/watchman/proto"

	"github.com/jmoiron/sqlx"
	"v2ray.com/core/watchman/database/connector"
	"v2ray.com/core/watchman/utils"
)

var Agent = &agentsController{
	Connector: defaultConnector,
}

const (
	updateUserTrafficStmt = "UPDATE `user` SET u = u + %d, d = d + %d WHERE id = ? "
	userTrafficLogStmt    = "INSERT INTO `user_traffic_log` (`id`, `user_id`, `u`, `d`, `node_id`, `rate`, `traffic`, `log_time`) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?) "
	nodeOnlineLogStmt     = "INSERT INTO `ss_node_online_log` (`id`, `node_id`, `online_user`, `log_time`) VALUES (NULL, ?, ?, ?) "
	aliveIPStmt           = "INSERT INTO `alive_ip` (`id`, `nodeid`,`userid`, `ip`, `datetime`) VALUES (NULL, ?, ?, ?, ?) "
	nodeHeartBeatStmt     = "UPDATE `ss_node` SET `node_heartbeat` = UNIX_TIMESTAMP(), `node_bandwidth` = `node_bandwidth` + %d WHERE id = ? "
)

// deviceAgentsController represents controller for 'device_agents'.
type agentsController struct {
	connector.Connector
}

// List Address interface{} by input
func (c *agentsController) GetAccounts(ctx context.Context, outputs interface{}) error {
	// stmt := "SELECT u.* from user u LEFT JOIN (SELECT userid, COUNT(DISTINCT ip) AS cnt FROM alive_ip WHERE datetime >= UNIX_TIMESTAMP(DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 MINUTE)) GROUP BY userid) u2 ON u.id = u2.userid WHERE u.enable = 1 AND u.uuid IS NOT NULL AND u.class_expire >= CURRENT_TIMESTAMP AND u.class > 1 AND u.transfer_enable > 0 AND (u2.cnt IS NULL OR u2.cnt <= u.node_connector);"
	stmt := "SELECT u.* from user u LEFT JOIN (SELECT userid, COUNT(DISTINCT ip) AS cnt FROM alive_ip WHERE datetime >= UNIX_TIMESTAMP(DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 MINUTE)) GROUP BY userid) u2 ON u.id = u2.userid WHERE u.enable = 1 AND u.uuid IS NOT NULL AND u.class_expire >= CURRENT_TIMESTAMP AND u.class > 1 AND u.transfer_enable > 0;"
	return c.Invoke(ctx, func(db connector.Q) error {
		return db.SelectContext(ctx, outputs, stmt) // nolint: errcheck
	})
}

// // List Address interface{} by input
func (c *agentsController) UpdateAccountTraffics(ctx context.Context, nodeId int64, account *proto.UserModel) error {
	return c.Invoke(ctx, func(tx *sqlx.Tx) (err error) {

		now := time.Now().Unix()

		tx.MustExec(tx.Rebind(fmt.Sprintf(updateUserTrafficStmt, account.Traffics.Uploads, account.Traffics.Downloads)), account.ID) // nolint: errcheck

		tx.MustExec(tx.Rebind(userTrafficLogStmt), account.ID, account.Traffics.Uploads, account.Traffics.Downloads, nodeId, 1, utils.GetDetectedSize(account.Traffics.Uploads+account.Traffics.Downloads), now) // nolint: errcheck

		tx.MustExec(tx.Rebind(nodeOnlineLogStmt), nodeId, account.Traffics.Clients, now) // nolint: errcheck

		for _, ip := range account.Traffics.IPs {
			tx.MustExec(tx.Rebind(aliveIPStmt), nodeId, account.ID, ip, now) // nolint: errcheck
		}

		return err
	})

}

// List Address interface{} by input
func (c *agentsController) Heartbeat(ctx context.Context, nodeId, traffic int64) error {
	return c.Invoke(ctx, func(db connector.Q) error {
		_, err := db.ExecContext(ctx, fmt.Sprintf(nodeHeartBeatStmt, traffic), nodeId) // nolint: errcheck
		return err
	})
}
