-- --------------------------------------------------------
-- 主机:                           192.168.10.1
-- 服务器版本:                        PostgreSQL 15.1 on aarch64-openwrt-linux-gnu, compiled by aarch64-openwrt-linux-musl-gcc (OpenWrt GCC 12.3.0 r24106-10cc5fcd00) 12.3.0, 64-bit
-- 服务器操作系统:                      
-- HeidiSQL 版本:                  12.4.0.6659
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES  */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- 导出  表 public.action_info 结构
CREATE TABLE IF NOT EXISTS "action_info" (
	"id" VARCHAR(36) NOT NULL,
	"scene_automation_id" VARCHAR(36) NOT NULL,
	"action_target" VARCHAR(255) NULL DEFAULT NULL,
	"action_type" VARCHAR(10) NOT NULL,
	"action_param_type" VARCHAR(10) NULL DEFAULT NULL,
	"action_param" VARCHAR(50) NULL DEFAULT NULL,
	"action_value" TEXT NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "action_info_scene_automations_fk" FOREIGN KEY ("scene_automation_id") REFERENCES "scene_automations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.action_info 的数据：-1 rows
/*!40000 ALTER TABLE "action_info" DISABLE KEYS */;
/*!40000 ALTER TABLE "action_info" ENABLE KEYS */;

-- 导出  表 public.alarm_config 结构
CREATE TABLE IF NOT EXISTS "alarm_config" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"alarm_level" VARCHAR(10) NOT NULL,
	"notification_group_id" VARCHAR(36) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"enabled" VARCHAR(10) NOT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.alarm_config 的数据：-1 rows
/*!40000 ALTER TABLE "alarm_config" DISABLE KEYS */;
/*!40000 ALTER TABLE "alarm_config" ENABLE KEYS */;

-- 导出  表 public.alarm_history 结构
CREATE TABLE IF NOT EXISTS "alarm_history" (
	"id" VARCHAR(36) NOT NULL,
	"alarm_config_id" VARCHAR(36) NOT NULL,
	"group_id" VARCHAR(36) NOT NULL,
	"scene_automation_id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"content" TEXT NULL DEFAULT NULL,
	"alarm_status" VARCHAR(3) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"create_at" TIMESTAMPTZ NOT NULL,
	"alarm_device_list" JSONB NOT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.alarm_history 的数据：-1 rows
/*!40000 ALTER TABLE "alarm_history" DISABLE KEYS */;
/*!40000 ALTER TABLE "alarm_history" ENABLE KEYS */;

-- 导出  表 public.alarm_info 结构
CREATE TABLE IF NOT EXISTS "alarm_info" (
	"id" VARCHAR(36) NOT NULL,
	"alarm_config_id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"alarm_time" TIMESTAMPTZ NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"content" TEXT NULL DEFAULT NULL,
	"processor" VARCHAR(36) NULL DEFAULT NULL,
	"processing_result" VARCHAR(10) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"alarm_level" VARCHAR(10) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "alarm_info_fk" FOREIGN KEY ("alarm_config_id") REFERENCES "alarm_config" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.alarm_info 的数据：-1 rows
/*!40000 ALTER TABLE "alarm_info" DISABLE KEYS */;
/*!40000 ALTER TABLE "alarm_info" ENABLE KEYS */;

-- 导出  表 public.attribute_datas 结构
CREATE TABLE IF NOT EXISTS "attribute_datas" (
	"id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"key" VARCHAR(255) NOT NULL,
	"ts" TIMESTAMPTZ NOT NULL,
	"bool_v" BOOLEAN NULL DEFAULT NULL,
	"number_v" DOUBLE PRECISION NULL DEFAULT NULL,
	"string_v" TEXT NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	UNIQUE INDEX "attribute_datas_device_id_key_key" ("device_id", "key"),
	CONSTRAINT "attribute_datas_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.attribute_datas 的数据：-1 rows
/*!40000 ALTER TABLE "attribute_datas" DISABLE KEYS */;
REPLACE INTO "attribute_datas" ("id", "device_id", "key", "ts", "bool_v", "number_v", "string_v", "tenant_id") VALUES
	('5fed1ff9-a753-b82d-7f9e-845eb028e036', '112233445566', 'macadress_wifi', '2025-02-13 14:15:14.622858+00', NULL, NULL, 'aabbccddeeff', 'd616bcbb'),
	('60051ab9-4ca6-e0bb-5907-7674260ac1fd', '112233445566', 'macadress_ble', '2025-02-13 14:15:14.622858+00', NULL, NULL, 'aaaaaaccccccc', 'd616bcbb'),
	('7a86ed8d-131d-3ae8-192b-77334c4de2a8', '112233445566', 'DeviceID', '2025-02-13 14:15:14.622858+00', NULL, 1, NULL, 'd616bcbb'),
	('dace08fb-0e9f-9d97-76ec-170221d74e21', '112233445566', 'Direct_connection', '2025-02-13 14:15:14.622858+00', NULL, 1, NULL, 'd616bcbb'),
	('238392cc-3bf4-daac-2b0c-087bad57b0e0', '112233445566', 'Device_type', '2025-02-13 14:15:14.622858+00', NULL, 0, NULL, 'd616bcbb');
/*!40000 ALTER TABLE "attribute_datas" ENABLE KEYS */;

-- 导出  表 public.attribute_set_logs 结构
CREATE TABLE IF NOT EXISTS "attribute_set_logs" (
	"id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"operation_type" VARCHAR(255) NULL DEFAULT NULL,
	"message_id" VARCHAR(36) NULL DEFAULT NULL,
	"data" TEXT NULL DEFAULT NULL,
	"rsp_data" TEXT NULL DEFAULT NULL,
	"status" VARCHAR(2) NULL DEFAULT NULL,
	"error_message" VARCHAR(500) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"user_id" VARCHAR(36) NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "attribute_set_logs_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.attribute_set_logs 的数据：-1 rows
/*!40000 ALTER TABLE "attribute_set_logs" DISABLE KEYS */;
/*!40000 ALTER TABLE "attribute_set_logs" ENABLE KEYS */;

-- 导出  表 public.boards 结构
CREATE TABLE IF NOT EXISTS "boards" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"config" JSON NULL DEFAULT '{}',
	"tenant_id" VARCHAR(36) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"home_flag" VARCHAR(2) NOT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"menu_flag" VARCHAR(2) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.boards 的数据：-1 rows
/*!40000 ALTER TABLE "boards" DISABLE KEYS */;
/*!40000 ALTER TABLE "boards" ENABLE KEYS */;

-- 导出  表 public.casbin_rule 结构
CREATE TABLE IF NOT EXISTS "casbin_rule" (
	"id" BIGINT NOT NULL DEFAULT 'nextval(''casbin_rule_id_seq''::regclass)',
	"ptype" VARCHAR(100) NULL DEFAULT NULL,
	"v0" VARCHAR(100) NULL DEFAULT NULL,
	"v1" VARCHAR(100) NULL DEFAULT NULL,
	"v2" VARCHAR(100) NULL DEFAULT NULL,
	"v3" VARCHAR(100) NULL DEFAULT NULL,
	"v4" VARCHAR(100) NULL DEFAULT NULL,
	"v5" VARCHAR(100) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "idx_casbin_rule" ("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
);

-- 正在导出表  public.casbin_rule 的数据：-1 rows
/*!40000 ALTER TABLE "casbin_rule" DISABLE KEYS */;
/*!40000 ALTER TABLE "casbin_rule" ENABLE KEYS */;

-- 导出  表 public.command_set_logs 结构
CREATE TABLE IF NOT EXISTS "command_set_logs" (
	"id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"operation_type" VARCHAR(255) NULL DEFAULT NULL,
	"message_id" VARCHAR(36) NULL DEFAULT NULL,
	"data" TEXT NULL DEFAULT NULL,
	"rsp_data" TEXT NULL DEFAULT NULL,
	"status" VARCHAR(2) NULL DEFAULT NULL,
	"error_message" VARCHAR(500) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"user_id" VARCHAR(36) NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"identify" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "command_set_logs_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.command_set_logs 的数据：-1 rows
/*!40000 ALTER TABLE "command_set_logs" DISABLE KEYS */;
REPLACE INTO "command_set_logs" ("id", "device_id", "operation_type", "message_id", "data", "rsp_data", "status", "error_message", "created_at", "user_id", "description", "identify") VALUES
	('51720b92-e76a-a9aa-b8ee-273c82f4cfe8', '112233445566', '1', '9456399', '{"method":"Device_restart","params":{"Reboot":1}}', NULL, '1', '', '2025-02-13 14:19:59.803694+00', '11111111-4fe9-b409-67c3-111111111111', '下发命令日志记录', 'Device_restart'),
	('9f11b77d-d2e3-f6b3-5fd6-3559800ffde0', '112233445566', '1', '9456554', '{"method":"Device_restart","params":{"Reboot":1}}', NULL, '1', '', '2025-02-13 14:22:34.446528+00', '11111111-4fe9-b409-67c3-111111111111', '下发命令日志记录', 'Device_restart'),
	('f5663fb6-53c4-e194-4f0b-46c12cbbc9ed', '112233445566', '1', '9456560', '{"method":"Runs_automatically","params":{"Mode":1}}', NULL, '1', '', '2025-02-13 14:22:40.196824+00', '11111111-4fe9-b409-67c3-111111111111', '下发命令日志记录', 'Runs_automatically'),
	('5a31808c-62e5-83ba-847a-773091b66608', '112233445566', '1', '9456566', '{"method":"Device_calibration_time","params":{"sntp":1}}', NULL, '1', '', '2025-02-13 14:22:46.882531+00', '11111111-4fe9-b409-67c3-111111111111', '下发命令日志记录', 'Device_calibration_time'),
	('3bacc988-a164-ef5e-3cdc-eccbd44e6bd7', '112233445566', '1', '9456572', '{"method":"stop","params":{"Mode":0}}', NULL, '1', '', '2025-02-13 14:22:52.091407+00', '11111111-4fe9-b409-67c3-111111111111', '下发命令日志记录', 'stop'),
	('c6036215-0022-8ada-f021-1271cbc79925', '112233445566', '1', '9456602', '{"method":"set_mode","params":{"Designated_azimuth":52,"Mode":1,"Specifies_height":86}}', NULL, '1', '', '2025-02-13 14:23:22.764357+00', '11111111-4fe9-b409-67c3-111111111111', '下发命令日志记录', 'set_mode');
/*!40000 ALTER TABLE "command_set_logs" ENABLE KEYS */;

-- 导出  表 public.data_policy 结构
CREATE TABLE IF NOT EXISTS "data_policy" (
	"id" VARCHAR(36) NOT NULL,
	"data_type" VARCHAR(1) NOT NULL,
	"retention_days" INTEGER NOT NULL,
	"last_cleanup_time" TIMESTAMPTZ NULL DEFAULT NULL,
	"last_cleanup_data_time" TIMESTAMPTZ NULL DEFAULT NULL,
	"enabled" VARCHAR(1) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.data_policy 的数据：-1 rows
/*!40000 ALTER TABLE "data_policy" DISABLE KEYS */;
REPLACE INTO "data_policy" ("id", "data_type", "retention_days", "last_cleanup_time", "last_cleanup_data_time", "enabled", "remark") VALUES
	('b', '2', 15, '2024-06-05 02:02:00.003+00', '2024-05-21 02:02:00.003+00', '1', ''),
	('a', '1', 15, '2024-06-05 02:02:00.003+00', '2024-05-21 02:02:00.101+00', '1', '');
/*!40000 ALTER TABLE "data_policy" ENABLE KEYS */;

-- 导出  表 public.data_scripts 结构
CREATE TABLE IF NOT EXISTS "data_scripts" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(99) NOT NULL,
	"device_config_id" VARCHAR(36) NOT NULL,
	"enable_flag" VARCHAR(9) NOT NULL,
	"content" TEXT NULL DEFAULT NULL,
	"script_type" VARCHAR(9) NOT NULL,
	"last_analog_input" TEXT NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "data_scripts_device_configs_fk" FOREIGN KEY ("device_config_id") REFERENCES "device_configs" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.data_scripts 的数据：-1 rows
/*!40000 ALTER TABLE "data_scripts" DISABLE KEYS */;
/*!40000 ALTER TABLE "data_scripts" ENABLE KEYS */;

-- 导出  表 public.devices 结构
CREATE TABLE IF NOT EXISTS "devices" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NULL DEFAULT NULL,
	"voucher" VARCHAR(500) NOT NULL DEFAULT '',
	"tenant_id" VARCHAR(36) NOT NULL DEFAULT '',
	"is_enabled" VARCHAR(36) NOT NULL DEFAULT '',
	"activate_flag" VARCHAR(36) NOT NULL DEFAULT '',
	"created_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"update_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"device_number" VARCHAR(36) NOT NULL DEFAULT '',
	"product_id" VARCHAR(36) NULL DEFAULT NULL,
	"parent_id" VARCHAR(36) NULL DEFAULT NULL,
	"protocol" VARCHAR(36) NULL DEFAULT NULL,
	"label" VARCHAR(255) NULL DEFAULT NULL,
	"location" VARCHAR(100) NULL DEFAULT NULL,
	"sub_device_addr" VARCHAR(36) NULL DEFAULT NULL,
	"current_version" VARCHAR(36) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT '{}',
	"protocol_config" JSON NULL DEFAULT '{}',
	"remark1" VARCHAR(255) NULL DEFAULT NULL,
	"remark2" VARCHAR(255) NULL DEFAULT NULL,
	"remark3" VARCHAR(255) NULL DEFAULT NULL,
	"device_config_id" VARCHAR(36) NULL DEFAULT NULL,
	"batch_number" VARCHAR(500) NULL DEFAULT NULL,
	"activate_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"is_online" SMALLINT NOT NULL DEFAULT '0',
	"access_way" VARCHAR(10) NULL DEFAULT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"service_access_id" VARCHAR(36) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "devices_unique" ("device_number"),
	UNIQUE INDEX "devices_unique_1" ("voucher"),
	CONSTRAINT "devices_service_access_fk" FOREIGN KEY ("service_access_id") REFERENCES "service_access" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
	CONSTRAINT "fk_device_config_id" FOREIGN KEY ("device_config_id") REFERENCES "device_configs" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT,
	CONSTRAINT "fk_product_id" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.devices 的数据：2 rows
/*!40000 ALTER TABLE "devices" DISABLE KEYS */;
REPLACE INTO "devices" ("id", "name", "voucher", "tenant_id", "is_enabled", "activate_flag", "created_at", "update_at", "device_number", "product_id", "parent_id", "protocol", "label", "location", "sub_device_addr", "current_version", "additional_info", "protocol_config", "remark1", "remark2", "remark3", "device_config_id", "batch_number", "activate_at", "is_online", "access_way", "description", "service_access_id") VALUES
	('8a09d81c-3d13-f159-f968-ec69140da60d', '气象站', '{"username":"27ae0698-99ac-33b4-4ec","password":""}', 'd616bcbb', '', 'active', '2025-02-13 14:01:37.926362+00', '2025-02-13 14:01:37.926362+00', '8a09d81c-3d13-f159-f968-ec69140da60d', NULL, NULL, NULL, '', NULL, NULL, NULL, '{}', '{}', NULL, NULL, NULL, '964d6220-ecbf-a043-1960-85b1a2758cea', NULL, NULL, 0, 'A', NULL, NULL),
	('112233445566', '112233445566', '{"username":"f5ee492d-62af-41ca-c52","password":"d1d55c0"}', 'd616bcbb', '', 'active', '2025-02-13 14:10:59.609943+00', '2025-02-13 14:10:59.609943+00', '112233445566', NULL, NULL, NULL, '112233445566', NULL, NULL, NULL, '{}', '{}', NULL, NULL, NULL, '315d9d82-5c76-3197-4eab-8c0a641ccdc9', NULL, NULL, 0, 'A', NULL, NULL);
/*!40000 ALTER TABLE "devices" ENABLE KEYS */;

-- 导出  表 public.device_configs 结构
CREATE TABLE IF NOT EXISTS "device_configs" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(99) NOT NULL,
	"device_template_id" VARCHAR(36) NULL DEFAULT NULL,
	"device_type" VARCHAR(9) NOT NULL,
	"protocol_type" VARCHAR(36) NULL DEFAULT NULL,
	"voucher_type" VARCHAR(36) NULL DEFAULT NULL,
	"protocol_config" JSON NULL DEFAULT NULL,
	"device_conn_type" VARCHAR(36) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT '{}',
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"other_config" JSON NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "device_configs_device_templates_fk" FOREIGN KEY ("device_template_id") REFERENCES "device_templates" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.device_configs 的数据：2 rows
/*!40000 ALTER TABLE "device_configs" DISABLE KEYS */;
REPLACE INTO "device_configs" ("id", "name", "device_template_id", "device_type", "protocol_type", "voucher_type", "protocol_config", "device_conn_type", "additional_info", "description", "tenant_id", "created_at", "updated_at", "remark", "other_config") VALUES
	('964d6220-ecbf-a043-1960-85b1a2758cea', '气象站模板', 'd391b336-4273-d101-27f8-d1fbedfde866', '1', 'MQTT', 'ACCESSTOKEN', '{"value":null}', NULL, '{}', NULL, 'd616bcbb', '2025-02-11 03:53:06.743747+00', '2025-02-11 03:53:24.657096+00', NULL, '{"online_timeout":5,"heartbeat":0}'),
	('315d9d82-5c76-3197-4eab-8c0a641ccdc9', 'ESP设备配置', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '1', 'MQTT', 'BASIC', '{}', NULL, '{}', NULL, 'd616bcbb', '2025-02-04 04:33:43.11584+00', '2025-02-13 14:29:37.353237+00', NULL, '{"online_timeout":10,"heartbeat":0}');
/*!40000 ALTER TABLE "device_configs" ENABLE KEYS */;

-- 导出  表 public.device_model_attributes 结构
CREATE TABLE IF NOT EXISTS "device_model_attributes" (
	"id" VARCHAR(36) NOT NULL,
	"device_template_id" VARCHAR(36) NOT NULL,
	"data_name" VARCHAR(255) NULL DEFAULT NULL,
	"data_identifier" VARCHAR(255) NOT NULL,
	"read_write_flag" VARCHAR(10) NULL DEFAULT NULL,
	"data_type" VARCHAR(50) NULL DEFAULT NULL,
	"unit" VARCHAR(50) NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "device_model_attributes_unique" ("device_template_id", "data_identifier"),
	CONSTRAINT "device_model_attributes_device_templates_fk" FOREIGN KEY ("device_template_id") REFERENCES "device_templates" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.device_model_attributes 的数据：5 rows
/*!40000 ALTER TABLE "device_model_attributes" DISABLE KEYS */;
REPLACE INTO "device_model_attributes" ("id", "device_template_id", "data_name", "data_identifier", "read_write_flag", "data_type", "unit", "description", "additional_info", "created_at", "updated_at", "remark", "tenant_id") VALUES
	('de8f4c22-fee7-df81-58f7-229b1e5d479e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备 蓝牙mac地址', 'macadress_ble', 'R-只读', 'String', '', '', '[]', '2025-02-06 12:34:46.045286+00', '2025-02-06 12:34:46.045286+00', NULL, 'd616bcbb'),
	('3db44364-c96b-4f69-381a-69d4463165fc', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备 wifi  mac地址', 'macadress_wifi', 'R-只读', 'String', '', '', '[]', '2025-02-06 12:35:48.439086+00', '2025-02-06 12:35:48.439086+00', NULL, 'd616bcbb'),
	('384ccd3d-25cc-0801-09d7-3474ba33fcd3', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备编号', 'DeviceID', 'R-只读', 'Number', '', '设备出厂编号', '[]', '2025-02-10 00:28:57.209706+00', '2025-02-10 00:29:16.160648+00', NULL, 'd616bcbb'),
	('a1cb372d-9e0a-9347-c608-4dc4a952a53e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '是否为直连设备', 'Direct_connection', 'RW-读/写', 'Number', '', '1（是），0（不是）
', '[]', '2025-02-11 01:08:10.07277+00', '2025-02-11 01:08:10.07277+00', NULL, 'd616bcbb'),
	('f903166f-a974-3f0b-6359-7759bc0aab45', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备类型', 'Device_type', 'RW-读/写', 'Number', '', '0（单轴）；1（斜单轴）；2（双轴）
', '[]', '2025-02-11 06:44:04.880618+00', '2025-02-11 06:45:14.710888+00', NULL, 'd616bcbb');
/*!40000 ALTER TABLE "device_model_attributes" ENABLE KEYS */;

-- 导出  表 public.device_model_commands 结构
CREATE TABLE IF NOT EXISTS "device_model_commands" (
	"id" VARCHAR(36) NOT NULL,
	"device_template_id" VARCHAR(36) NOT NULL,
	"data_name" VARCHAR(255) NULL DEFAULT NULL,
	"data_identifier" VARCHAR(255) NOT NULL,
	"params" JSON NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "device_model_commands_unique" ("data_identifier", "device_template_id"),
	CONSTRAINT "device_model_commands_device_templates_fk" FOREIGN KEY ("device_template_id") REFERENCES "device_templates" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.device_model_commands 的数据：8 rows
/*!40000 ALTER TABLE "device_model_commands" DISABLE KEYS */;
REPLACE INTO "device_model_commands" ("id", "device_template_id", "data_name", "data_identifier", "params", "description", "additional_info", "created_at", "updated_at", "remark", "tenant_id") VALUES
	('fc214760-206f-2692-4d90-d51164c9909b', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向东运行', 'set_motion_east', '[{"data_name":"向东运行时间","data_identifier":"east","param_type":"Number","description":"设备向东运行多少秒，设备方位角变小","data_type":"string","enum_config":[],"id":0.966594102891416}]', '', NULL, '2025-02-10 02:11:20.450619+00', '2025-02-10 02:14:51.180783+00', NULL, 'd616bcbb'),
	('b9cbfb02-5ece-d213-7621-c5908f483daa', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向北运行', 'set_motion_north', '[{"data_name":"向北运行时间","data_identifier":"north","param_type":"Number","description":"设备向北方向运行多少秒，设备高度角变大","data_type":"string","enum_config":[],"id":0.17887362350128666}]', '', NULL, '2025-02-10 02:05:53.714792+00', '2025-02-10 02:15:02.470042+00', NULL, 'd616bcbb'),
	('ffd2670f-9199-b5d8-3264-879ad2b7ef7e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '初次上电设置参数', 'set_first_power', '[{"data_name":"设备名称","data_identifier":"Devicename","param_type":"string","description":"设备现场名称编号","data_type":"string","enum_config":[],"id":0.2493764878392466},{"data_name":"蓝牙编号","data_identifier":"bleid","param_type":"String","description":"","enum_config":[],"id":0.4590131617040294},{"data_name":"帆板长","data_identifier":"Long","param_type":"Number","description":"设备光伏帆板长度，单位厘米","enum_config":[],"id":0.5847989918191545},{"data_name":"帆板宽","data_identifier":"Wide","param_type":"Number","description":"设备光伏帆板宽度，单位厘米","enum_config":[],"id":0.0017696418924946222},{"data_name":"设备阵列东西间距","data_identifier":"East_West_spacing","param_type":"Number","description":"设备间东西间隔，单位厘米","enum_config":[],"id":0.310698242363485},{"data_name":"设备阵列南北间距","data_identifier":"North_south_spacing","param_type":"Number","description":"设备间南北间隔，单位厘米","enum_config":[],"id":0.9336989106637725},{"data_name":"设备类型","data_identifier":"Device_type","param_type":"Number","description":"设备类型分三种：\n0（单轴）；1（斜单轴）；2（双轴）\n","enum_config":[],"id":0.5659839822647621},{"data_name":"限制东","data_identifier":"Restriction_East","param_type":"Number","description":"设备方位角电子围栏，最低（范围-180~0）","enum_config":[],"id":0.8699351887804592},{"data_name":"限制西","data_identifier":"Restriction_West","param_type":"Number","description":"设备方位角电子围栏。最大（范围0~180）","enum_config":[],"id":0.942441830846289},{"data_name":"限制南","data_identifier":"Restriction_South","param_type":"Number","description":"设备高度角电子围栏，最小（范围0~90），低于限制北","enum_config":[],"id":0.7213654440138069},{"data_name":"限制北","data_identifier":"Restriction_North","param_type":"string","description":"设备高度角电子围栏，最大（范围0~90）","enum_config":[],"id":0.46858871706865135},{"data_name":"风速预警值","data_identifier":"wind_error","param_type":"Number","description":"设备进入防风模式风速触发值，单位米每秒","enum_config":[],"id":0.0015485381903646012},{"data_name":"解除预警风速值","data_identifier":"wind_anquan","param_type":"Number","description":"解除防风风速大小，单位米每秒","enum_config":[],"id":0.8517626357309545},{"data_name":"电机占空比","data_identifier":"Duty","param_type":"Number","description":"控制电机转速，范围（0~100）","enum_config":[],"id":0.4466173258031818},{"data_name":"回转减速比","data_identifier":"Slewing_reduction","param_type":"Number","description":"设备转一圈，编码器转的圈数","enum_config":[],"id":0.9033323841663345},{"data_name":"跟踪精度","data_identifier":"track_longitude","param_type":"Number","description":"设备自动运行，跟踪太阳位置的误差","enum_config":[],"id":0.6207050303798538},{"data_name":"回位太阳高度角","data_identifier":"Return_sun_altitude","param_type":"Number","description":"设备归位时，太阳高度","enum_config":[],"id":0.24141308771796233},{"data_name":"高度角电机限制电流","data_identifier":"H_current_limiting","param_type":"Number","description":"高度电机正常运行电流安全范围，单位毫安","enum_config":[],"id":0.7169837409795439},{"data_name":"方位角（旋转角）电机限制电流","data_identifier":"A_current_limiting","param_type":"Number","description":"方位电机正常运行电流安全范围，单位毫安","enum_config":[],"id":0.24110668788592782},{"data_name":"设备经度","data_identifier":"Longitude","param_type":"Number","description":"","enum_config":[],"id":0.9715377902582469},{"data_name":"设备纬度","data_identifier":"Latitude","param_type":"Number","description":"","enum_config":[],"id":0.4547693805049029},{"data_name":"高度角电机电流获取延时","data_identifier":"H_delay","param_type":"Number","description":"获取电机电流缓冲时间，单位毫秒","enum_config":[],"id":0.6606283304349359},{"data_name":"方位角（旋转角）电机电流获取延时","data_identifier":"A_delay","param_type":"Number","description":"方位电机电流获取缓冲时间，单位毫秒","enum_config":[],"id":0.6701381265088022},{"data_name":"回味设备高度角","data_identifier":"Return_device_altitude","param_type":"Number","description":"设备归位状态高度角","enum_config":[],"id":0.5860188549500094},{"data_name":"是否为直连设备","data_identifier":"Direct_connection","param_type":"Number","description":"1（是），0（不是）\n","enum_config":[],"id":0.2012483969904828},{"data_name":"是否开启规避","data_identifier":"is_circumvent","param_type":"Number","description":"0（不规避）1（规避\n","enum_config":[],"id":0.45905653060446894},{"data_name":"年","data_identifier":"year","param_type":"Number","description":"例子：2024年就输入24","data_type":"string","enum_config":[],"id":0.3335337204008233},{"data_name":"月","data_identifier":"month","param_type":"Number","description":"1~12月","enum_config":[],"id":0.5998522792038801},{"data_name":"日","data_identifier":"day","param_type":"Number","description":"","enum_config":[],"id":0.733341139709766},{"data_name":"时","data_identifier":"hour","param_type":"Number","description":"","enum_config":[],"id":0.3947520880555486},{"data_name":"分钟","data_identifier":"minute","param_type":"Number","description":"","enum_config":[],"id":0.49268651161172805},{"data_name":"秒","data_identifier":"second","param_type":"Number","description":"","enum_config":[],"id":0.1471325091171114}]', '设备上电需要设置全部参数，设备默认进入调试模式，参数修改后需要重新调整设备运行模式', NULL, '2025-02-10 01:54:26.796393+00', '2025-02-11 06:52:00.412962+00', NULL, 'd616bcbb'),
	('28149e55-0a73-6e6c-0ced-e8e7770995e6', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向西运行', 'set_motion_west', '[{"data_name":"向西运行时间","data_identifier":"west","param_type":"Number","description":"向西运行多少秒，方位角变大","data_type":"string","enum_config":[],"id":0.6433723482906932}]', '', NULL, '2025-02-10 02:13:14.044048+00', '2025-02-10 02:13:47.60848+00', NULL, 'd616bcbb'),
	('5ea5e0d1-3b39-4336-a766-75e7ffbab77e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向南运行', 'set_motion_south', '[{"data_name":"向南运行时间","data_identifier":"south","param_type":"Number","description":"设备向南方向运行多少秒，高度角变小","data_type":"string","enum_config":[],"id":0.6641119283713559}]', '', NULL, '2025-02-10 02:08:37.000437+00', '2025-02-10 02:14:08.892981+00', NULL, 'd616bcbb'),
	('341dd50a-04ea-7787-c871-0690575ec267', '4c2f6999-630f-a358-bb7d-8eb5130a502c', 'ota升级', 'ota_up', '[{"data_name":"ota","data_identifier":"ota_up","param_type":"string","description":"0（停止）1（进行ota升级）\n","data_type":"string","enum_config":[],"id":0.3695510601508296}]', '', NULL, '2025-02-10 02:21:08.083524+00', '2025-02-10 02:21:08.083524+00', NULL, 'd616bcbb'),
	('5db84f1c-b3de-d068-23bc-6517df644a23', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设置运行模式', 'set_mode', '[{"data_name":"运行模式","data_identifier":"Mode","param_type":"Number","description":"0\t停止运行(调试模式)\n1\t自动运行\n2\t指定位置运行\n3\t模拟运行\n4\t防风模式\n5\t防雪模式\n6\t防冰雹模式\n7\t清洗模式（大雨天气）\n","data_type":"string","enum_config":[],"id":0.1907584242966347},{"data_name":"指定方位角","data_identifier":"Designated_azimuth","param_type":"Number","description":"设备运行到指定方位角","enum_config":[],"id":0.9881065632970181},{"data_name":"指定高度角","data_identifier":"Specifies_height","param_type":"Number","description":"设备运行到指定高度角","enum_config":[],"id":0.03491042177275139}]', '修改设备运行模式.停止运行与自动运行模式可不下发指定高度角，指定方位角数据。', NULL, '2025-02-10 02:00:37.061584+00', '2025-02-12 07:36:22.752621+00', NULL, 'd616bcbb'),
	('f7fad6b6-c1bb-c0e1-8b70-77261b81725f', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '运维更改参数', 'O&M_change_parameters', '[{"data_name":"帆板长","data_identifier":"Long","param_type":"Number","description":"设备光伏帆板长度，单位厘米","data_type":"string","enum_config":[],"id":0.11979955135138765},{"data_name":"帆板宽","data_identifier":"Wide","param_type":"Number","description":"设备光伏帆板宽度，单位厘米","enum_config":[],"id":0.595651212611531},{"data_name":"限制东","data_identifier":"Restriction_East","param_type":"Number","description":"设备方位角电子围栏，最低（范围-180~0）","enum_config":[],"id":0.838184601192526},{"data_name":"限制西","data_identifier":"Restriction_West","param_type":"Number","description":"设备方位角电子围栏。最大（范围0~180）","enum_config":[],"id":0.32298658048283135},{"data_name":"限制南","data_identifier":"Restriction_South","param_type":"Number","description":"设备高度角电子围栏，最小（范围0~90），低于限制北","enum_config":[],"id":0.6539817054531625},{"data_name":"限制北","data_identifier":"Restriction_North","param_type":"Number","description":"设备高度角电子围栏，最大（范围0~90）","enum_config":[],"id":0.7613142852799506},{"data_name":"风速预警值","data_identifier":"wind_error","param_type":"Number","description":"防风模式触发风速大小","enum_config":[],"id":0.27055724679828375},{"data_name":"解除预警风速","data_identifier":"wind_anquan","param_type":"Number","description":"解除防风模式时风速在该值以下","enum_config":[],"id":0.2033414962629585},{"data_name":"电机占空比","data_identifier":"Duty","param_type":"Number","description":"电机转速","enum_config":[],"id":0.1470929546594928},{"data_name":"回转减速比","data_identifier":"Slewing_reduction","param_type":"Number","description":"设备转一圈，编码器转的圈数","enum_config":[],"id":0.6221065042052765},{"data_name":"跟踪精度","data_identifier":"track_longitude","param_type":"Number","description":"设备跟踪太阳运动的设备位置误差范围","enum_config":[],"id":0.06207506096382365},{"data_name":"回位太阳高度角","data_identifier":"Return_sun_altitude","param_type":"Number","description":"太阳高度在该值以下，设备进入回位状态","enum_config":[],"id":0.762963270109311},{"data_name":"高度角电机限制电流","data_identifier":"H_current_limiting","param_type":"Number","description":"高度电机安全运行电流范围，单位毫安","enum_config":[],"id":0.21474882169414244},{"data_name":"方位角（旋转角）电机限制电流","data_identifier":"A_current_limiting","param_type":"Number","description":"方位电机安全运行电流范围，单位毫安","enum_config":[],"id":0.8858020182262034},{"data_name":"高度角电机电流获取延时","data_identifier":"H_delay","param_type":"Number","description":"缓冲获取高度电机电流时间，单位毫秒","enum_config":[],"id":0.18670328206494835},{"data_name":"方位角（旋转角）电机电流获取延时","data_identifier":"A_delay","param_type":"Number","description":"方位电机电流获取缓冲时间，单位毫秒","enum_config":[],"id":0.6336412534343263},{"data_name":"设备经度","data_identifier":"Longitude","param_type":"Number","description":"","enum_config":[],"id":0.111952542897787},{"data_name":"设备纬度","data_identifier":"Latitude","param_type":"Number","description":"","enum_config":[],"id":0.4852835393446615},{"data_name":"回位设备高度角","data_identifier":"Return_device_altitude","param_type":"Number","description":"设备回位状态运行到该高度","enum_config":[],"id":0.4717365899964667},{"data_name":"是否开启规避","data_identifier":"is_circumvent","param_type":"string","description":"0（不规避）1（规避\n","enum_config":[],"id":0.35507198415102037},{"data_name":"年","data_identifier":"year","param_type":"Number","description":"例子： 2024年便输入24","enum_config":[],"id":0.11302261916408218},{"data_name":"月","data_identifier":"month","param_type":"Number","description":"","data_type":"string","enum_config":[],"id":0.572883190279778},{"data_name":"日","data_identifier":"day","param_type":"Number","description":"","enum_config":[],"id":0.8549404746480156},{"data_name":"时","data_identifier":"hour","param_type":"Number","description":"","enum_config":[],"id":0.5599423871683895},{"data_name":"分钟","data_identifier":"minute","param_type":"Number","description":"","enum_config":[],"id":0.08247950781460012},{"data_name":"秒","data_identifier":"second","param_type":"Number","description":"","enum_config":[],"id":0.44008578549400346}]', '运维人员后续修改参数，设备默认进入调试模式，参数修改后需要重新调整设备运行模式', NULL, '2025-02-10 07:09:11.21809+00', '2025-02-11 06:51:40.151999+00', NULL, 'd616bcbb');
/*!40000 ALTER TABLE "device_model_commands" ENABLE KEYS */;

-- 导出  表 public.device_model_custom_commands 结构
CREATE TABLE IF NOT EXISTS "device_model_custom_commands" (
	"id" VARCHAR(36) NOT NULL,
	"device_template_id" VARCHAR(36) NOT NULL,
	"buttom_name" VARCHAR(255) NOT NULL,
	"data_identifier" VARCHAR(255) NOT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"instruct" TEXT NULL DEFAULT NULL,
	"enable_status" VARCHAR(10) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR NOT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.device_model_custom_commands 的数据：4 rows
/*!40000 ALTER TABLE "device_model_custom_commands" DISABLE KEYS */;
REPLACE INTO "device_model_custom_commands" ("id", "device_template_id", "buttom_name", "data_identifier", "description", "instruct", "enable_status", "remark", "tenant_id") VALUES
	('b34540de-fc70-1887-10a4-34b126150337', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备停止运行（调试模式）', 'stop', '设备停止运行', '{
"Mode":0

}', 'enable', NULL, 'd616bcbb'),
	('d30d5209-01d6-3300-60f4-32e3be7cb9d6', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备重启', 'Device_restart', '设备重启', '{
"Reboot":1
}', 'enable', NULL, 'd616bcbb'),
	('c22e7f49-7eaa-6268-e673-9c10d90012fa', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备自动运行', 'Runs_automatically', '设备进入自动运行模式', '{
"Mode":1
}', 'enable', NULL, 'd616bcbb'),
	('34dfa37b-7006-81cd-9550-25a09d05affe', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备校准时间', 'Device_calibration_time', '设备校准时间', '{
"sntp":1
}', 'enable', NULL, 'd616bcbb');
/*!40000 ALTER TABLE "device_model_custom_commands" ENABLE KEYS */;

-- 导出  表 public.device_model_custom_control 结构
CREATE TABLE IF NOT EXISTS "device_model_custom_control" (
	"id" VARCHAR(36) NOT NULL,
	"device_template_id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"control_type" VARCHAR NOT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"content" TEXT NULL DEFAULT NULL,
	"enable_status" VARCHAR(10) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "device_model_custom_control_device_templates_fk" FOREIGN KEY ("device_template_id") REFERENCES "device_templates" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.device_model_custom_control 的数据：-1 rows
/*!40000 ALTER TABLE "device_model_custom_control" DISABLE KEYS */;
/*!40000 ALTER TABLE "device_model_custom_control" ENABLE KEYS */;

-- 导出  表 public.device_model_events 结构
CREATE TABLE IF NOT EXISTS "device_model_events" (
	"id" VARCHAR(36) NOT NULL,
	"device_template_id" VARCHAR(36) NOT NULL,
	"data_name" VARCHAR(255) NULL DEFAULT NULL,
	"data_identifier" VARCHAR(255) NOT NULL,
	"params" JSON NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "device_model_events_unique" ("device_template_id", "data_identifier"),
	CONSTRAINT "device_model_events_device_templates_fk" FOREIGN KEY ("device_template_id") REFERENCES "device_templates" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.device_model_events 的数据：2 rows
/*!40000 ALTER TABLE "device_model_events" DISABLE KEYS */;
REPLACE INTO "device_model_events" ("id", "device_template_id", "data_name", "data_identifier", "params", "description", "additional_info", "created_at", "updated_at", "remark", "tenant_id") VALUES
	('1be79015-c06a-b478-1f2b-2f5d1ed9b73b', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备故障', 'Equipment_failure', '[{"data_name":"故障码","data_identifier":"Device_error","read_write_flag":"Number","description":"0\t无故障\n101\t温度读取失败\n102\t电压异常（低于正常值）\n103\tmesh组网失败\n104\t倾角传感器错误\n105\t编码器错误\n106\t超出限制北范围\n107\t超出限制南范围\n108\t超出限制西范围\n109\t超出限制东范围\n110\t风速传感器错误\n111\t高度角电机电流超出\n112\t方位角电机电流超出\n115\t编码器，倾角传感器全部通信错误\n116\tgps获取错误\n117\t4g通信失败\n118\tota升级失败\n","id":0.07517099198017196}]', '', NULL, '2025-02-11 06:53:34.133389+00', '2025-02-11 06:53:34.133389+00', NULL, 'd616bcbb'),
	('155f188d-6207-a8a0-d02e-7565f709c6d3', '4c2f6999-630f-a358-bb7d-8eb5130a502c', 'ota升级', 'OTA_UP', '[{"data_name":"ota升级","data_identifier":"otaup","read_write_flag":"Number","description":"1（进行一次ota升级）\n","id":0.5817395082492451}]', '', NULL, '2025-02-11 06:57:12.043546+00', '2025-02-12 07:32:25.9668+00', NULL, 'd616bcbb');
/*!40000 ALTER TABLE "device_model_events" ENABLE KEYS */;

-- 导出  表 public.device_model_telemetry 结构
CREATE TABLE IF NOT EXISTS "device_model_telemetry" (
	"id" VARCHAR(36) NOT NULL,
	"device_template_id" VARCHAR(36) NOT NULL,
	"data_name" VARCHAR(255) NULL DEFAULT NULL,
	"data_identifier" VARCHAR(255) NOT NULL,
	"read_write_flag" VARCHAR(10) NULL DEFAULT NULL,
	"data_type" VARCHAR(50) NULL DEFAULT NULL,
	"unit" VARCHAR(50) NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "device_model_telemetry_unique" ("device_template_id", "data_identifier"),
	CONSTRAINT "device_model_telemetry_device_templates_fk" FOREIGN KEY ("device_template_id") REFERENCES "device_templates" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.device_model_telemetry 的数据：62 rows
/*!40000 ALTER TABLE "device_model_telemetry" DISABLE KEYS */;
REPLACE INTO "device_model_telemetry" ("id", "device_template_id", "data_name", "data_identifier", "read_write_flag", "data_type", "unit", "description", "additional_info", "created_at", "updated_at", "remark", "tenant_id") VALUES
	('5ec3a8a9-5910-292c-c949-eb8c57b1d64a', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '限制北', 'Restriction_North', 'RW-读/写', 'Number', '度', '设备高度角电子围栏（最大）范围：0~90', '[]', '2025-02-06 11:58:03.904045+00', '2025-02-06 11:58:03.904045+00', NULL, 'd616bcbb'),
	('8795f915-2aa3-7c34-c9f7-0692f22b44f0', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备运行模式', 'Mode', 'RW-读/写', 'Number', '', '0（停止运行）1（自动运行）2（指定位置运行）3（模拟运行），4（防风模式），5（防雪模式），6（防冰雹模式），7（清洗模式）', '[]', '2025-02-04 11:35:52.427059+00', '2025-02-04 11:35:52.427059+00', NULL, 'd616bcbb'),
	('70529138-97c1-2344-b238-52add0de01d7', '4c2f6999-630f-a358-bb7d-8eb5130a502c', 'OTA升级', 'ota_up', 'RW-读/写', 'Number', '', '下方设备进行ota升级 ，1表示执行
', '[]', '2025-02-04 11:37:06.915134+00', '2025-02-04 11:37:06.915134+00', NULL, 'd616bcbb'),
	('abdeb9f4-49c9-7aa1-c2e7-8686c612b41f', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '重启', 'Reboot', 'RW-读/写', 'Number', '', '让设备进行重启（0（停止）1（运行））', '[]', '2025-02-04 11:38:04.873817+00', '2025-02-04 11:38:04.873817+00', NULL, 'd616bcbb'),
	('57f618a7-4ca7-7a57-5a21-88e0b016ddb6', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '风速预警值', 'wind_error', 'RW-读/写', 'Number', '米每秒', '设备进入防风触发风速', '[]', '2025-02-06 11:58:33.034955+00', '2025-02-06 11:58:33.034955+00', NULL, 'd616bcbb'),
	('cd7eb92a-5cec-b8d2-2bd5-9906d7745571', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '风速', 'wind_speed', 'R-只读', 'Number', '米每秒', '风速传感器风速值', '[]', '2025-02-06 12:09:54.612257+00', '2025-02-06 12:09:54.612257+00', NULL, 'd616bcbb'),
	('f6c6ce5d-ee9b-59d7-31e3-7493b9b4f08e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '校准时间', 'sntp', 'RW-读/写', 'Number', '', '让设备向上获取时间，并进行校准
（0（停止）1（运行））', '[]', '2025-02-04 11:35:10.386412+00', '2025-02-04 11:38:36.532661+00', NULL, 'd616bcbb'),
	('1abe8670-0659-cfed-fff0-87397d0fd384', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '解除预警风速', 'wind_anquan', 'RW-读/写', 'Number', '米每秒', '设备接触防风风速', '[]', '2025-02-06 11:58:55.392544+00', '2025-02-06 11:58:55.392544+00', NULL, 'd616bcbb'),
	('46fc0bca-a99d-05ed-eee4-71e9695c4dfb', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向南运行', 'south', 'RW-读/写', 'Number', '秒', '设备向南转动时间', '[]', '2025-02-04 11:30:32.778982+00', '2025-02-04 11:40:46.892772+00', NULL, 'd616bcbb'),
	('e9f73051-7b1e-67bb-63af-174ac9f0803e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向西转动', 'west', 'RW-读/写', 'Number', '秒', '设备向西进行转动时间', '[]', '2025-02-04 11:32:45.294702+00', '2025-02-04 11:40:55.584342+00', NULL, 'd616bcbb'),
	('c093e3c6-1047-6b08-0a23-3f28c27be892', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向东运行', 'east', 'RW-读/写', 'Number', '秒', '设备向北方向转动', '[]', '2025-02-04 11:31:28.566198+00', '2025-02-04 11:41:04.36283+00', NULL, 'd616bcbb'),
	('0c42ae62-2d21-e364-9309-d150659a59c3', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备向北运行', 'north', 'RW-读/写', 'Number', '秒', '设备向北转动时间', '[]', '2025-02-04 11:29:51.84498+00', '2025-02-04 11:41:12.140191+00', NULL, 'd616bcbb'),
	('588799da-e96a-fe3f-b58f-f2ddf0e5b585', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备名称', 'Devicename', 'RW-读/写', 'String', '', '设备现场名称', '[]', '2025-02-06 11:50:35.790829+00', '2025-02-06 11:50:35.790829+00', NULL, 'd616bcbb'),
	('30dc50d6-6b21-9bca-b498-89b33b8af22a', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '蓝牙编号', 'bleid', 'RW-读/写', 'String', '', '蓝牙名称', '[]', '2025-02-06 11:50:56.769979+00', '2025-02-06 11:50:56.769979+00', NULL, 'd616bcbb'),
	('5872433e-0f96-4ea7-cb2f-55d11e7ce441', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '帆板长', 'Long', 'RW-读/写', 'Number', '厘米', '设备太阳能板长度', '[]', '2025-02-06 11:51:29.700846+00', '2025-02-06 11:51:29.700846+00', NULL, 'd616bcbb'),
	('319a61bd-6bd8-b8cb-fffb-39a95f5c3b48', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '帆板宽', 'Wide', 'RW-读/写', 'Number', '厘米', '设备太阳能板宽度', '[]', '2025-02-06 11:51:57.043335+00', '2025-02-06 11:51:57.043335+00', NULL, 'd616bcbb'),
	('ef0c2960-cea0-44dd-a65f-e62fa9fed4d6', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备阵列东西间距', 'East_West_spacing', 'RW-读/写', 'Number', '厘米', '不同设备间东西间距', '[]', '2025-02-06 11:52:23.391006+00', '2025-02-06 11:52:23.391006+00', NULL, 'd616bcbb'),
	('3be5b352-0330-0906-67ce-c2310f48e066', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备阵列南北间距', 'North_south_spacing', 'RW-读/写', 'Number', '厘米', '不同设备间南北间距', '[]', '2025-02-06 11:52:50.879332+00', '2025-02-06 11:52:50.879332+00', NULL, 'd616bcbb'),
	('da5da4aa-7848-5ab5-a85a-afacd1a03782', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '限制东', 'Restriction_East', 'RW-读/写', 'Number', '度', '设备方位角电子围栏（最低）：
范围：-180~0', '[]', '2025-02-06 11:54:51.806676+00', '2025-02-06 11:54:51.806676+00', NULL, 'd616bcbb'),
	('b26d0d62-03c1-2f0e-322f-44fcd8268df8', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '限制西', 'Restriction_West', 'RW-读/写', 'Number', '度', '设备方位角电子围栏（最大），
范围：0~180', '[]', '2025-02-06 11:55:40.096347+00', '2025-02-06 11:55:40.096347+00', NULL, 'd616bcbb'),
	('1668617f-8b0e-29a5-7cb1-4bf93c8a99a4', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '限制南', 'Restriction_South', 'RW-读/写', 'Number', '度', '设备高度角电子围栏（最低）范围：0~90', '[]', '2025-02-06 11:57:18.62805+00', '2025-02-06 11:57:18.62805+00', NULL, 'd616bcbb'),
	('09d2d1e8-bc77-e6df-a38e-10e518b26d53', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '回转减速比', 'Slewing_reduction', 'RW-读/写', 'Number', '', '设备转一圈，编码器需要转的数值', '[]', '2025-02-06 12:00:04.387335+00', '2025-02-06 12:00:04.387335+00', NULL, 'd616bcbb'),
	('ae73c4fc-1c38-5537-a5a1-07cc067dac7b', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '电机占空比', 'Duty', 'RW-读/写', 'Number', '', '电机转速', '[]', '2025-02-06 11:59:35.063322+00', '2025-02-06 12:00:34.936841+00', NULL, 'd616bcbb'),
	('0a6d9444-4fd7-7bf4-e764-eedc3d4fc227', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '跟踪精度', 'track_longitude', 'RW-读/写', 'Number', '度', '设备跟踪太阳缓冲大小', '[]', '2025-02-06 12:01:02.120483+00', '2025-02-06 12:01:02.120483+00', NULL, 'd616bcbb'),
	('17983eee-d418-4a42-f89f-e6f47a4e53b9', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '回位太阳高度角', 'Return_sun_altitude', 'RW-读/写', 'Number', '度', '设备归为状态是太阳所处高度', '[]', '2025-02-06 12:01:33.891624+00', '2025-02-06 12:01:33.891624+00', NULL, 'd616bcbb'),
	('ffb5818e-44e5-06f9-94bf-85753678326d', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '高度角电机限制电流', 'H_current_limiting', 'RW-读/写', 'Number', '毫安', '高度电机安全运行电流大小', '[]', '2025-02-06 12:02:15.880279+00', '2025-02-06 12:02:15.880279+00', NULL, 'd616bcbb'),
	('628cd87c-c2d9-8628-0978-fdb5901642cf', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '方位角（旋转角）电机限制电流', 'A_current_limiting', 'RW-读/写', 'Number', '毫安', '方位电机安全运行电流范围', '[]', '2025-02-06 12:02:41.777944+00', '2025-02-06 12:02:41.777944+00', NULL, 'd616bcbb'),
	('fee37769-96d2-d56f-b397-9be525b90fb4', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '高度角电机电流获取延时', 'H_delay', 'RW-读/写', 'Number', '毫秒', '获取高度电机电流延迟', '[]', '2025-02-06 12:03:19.827041+00', '2025-02-06 12:03:19.827041+00', NULL, 'd616bcbb'),
	('8ced01bf-4518-0ee6-7991-3b311dbde2ab', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备系统时间', 'System_time', 'RW-读/写', 'String', '', '设备自身时间，格式：年/月/日/时/分/秒', '[]', '2025-02-06 12:06:30.776143+00', '2025-02-06 12:06:30.776143+00', NULL, 'd616bcbb'),
	('6faab688-e54c-3bd7-9b22-a6aad697a6b4', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备经度', 'Longitude', 'RW-读/写', 'Number', '度', '', '[]', '2025-02-06 12:03:38.28857+00', '2025-02-06 12:07:06.912674+00', NULL, 'd616bcbb'),
	('88c76485-3848-646e-5210-55755ad779a5', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备标定', 'set_demarcate', 'RW-读/写', 'Number', '', '设备自身进行标定，补偿位置数据
1，该时间进行标定，0默认初始值', '[]', '2025-02-04 11:33:29.250986+00', '2025-02-11 06:29:17.824415+00', NULL, 'd616bcbb'),
	('79519a8c-3be8-7602-f023-2bfb0ba17c3a', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设置中点', 'set_midpoint', 'RW-读/写', 'Number', '', '将设备方位角信息置0
（0（停止）1（运行））', '[]', '2025-02-04 11:34:27.891245+00', '2025-02-11 07:13:47.135515+00', NULL, 'd616bcbb'),
	('55d037fb-dc17-150e-26ac-3d82cc2fa36d', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备纬度', 'Latitude', 'RW-读/写', 'Number', '度', '', '[]', '2025-02-06 12:03:54.380042+00', '2025-02-06 12:07:00.273048+00', NULL, 'd616bcbb'),
	('9fdac67c-8501-04b1-a414-342a915d6215', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '电控板温度', 'temp', 'R-只读', 'Number', '摄氏度', '电控盒温度', '[]', '2025-02-06 12:08:39.912013+00', '2025-02-06 12:08:39.912013+00', NULL, 'd616bcbb'),
	('52aa6fc9-1e37-9ce3-62b7-7f83ff1bea08', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '电池电压', 'Battery_voltage', 'R-只读', 'Number', '毫伏', '内置电池电压', '[]', '2025-02-06 12:09:13.213346+00', '2025-02-06 12:09:13.213346+00', NULL, 'd616bcbb'),
	('c095b174-7279-d375-5b96-08c3e39cc487', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '真实太阳方位角', 'A_true', 'R-只读', 'Number', '度', '', '[]', '2025-02-06 12:10:13.580365+00', '2025-02-06 12:10:13.580365+00', NULL, 'd616bcbb'),
	('d5e63299-d579-f27b-0aff-58cb802edbf0', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '真实太阳高度角', 'H_true', 'R-只读', 'Number', '度', '', '[]', '2025-02-06 12:10:25.915185+00', '2025-02-06 12:10:25.915185+00', NULL, 'd616bcbb'),
	('5e1e76fb-991d-d710-100a-2ba88fb2b5f9', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '公式计算规避后太阳高度角', 'sun_h', 'R-只读', 'Number', '度', '', '[]', '2025-02-06 12:10:40.391607+00', '2025-02-06 12:10:40.391607+00', NULL, 'd616bcbb'),
	('9364764e-ea65-2e4d-ac2e-e197241a34e6', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '公式计算规避后太阳方位角', 'sun_a', 'R-只读', 'Number', '度', '', '[]', '2025-02-06 12:10:56.112494+00', '2025-02-06 12:10:56.112494+00', NULL, 'd616bcbb'),
	('de228fdb-ad90-b31f-e45b-5b480d1d3be7', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '单轴设备固定角', 'Device_fixed', 'R-只读', 'Number', '度', '', '[]', '2025-02-06 12:11:10.296965+00', '2025-02-06 12:11:10.296965+00', NULL, 'd616bcbb'),
	('03ced83c-9447-d17a-04ff-5a898e6093cd', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '公式计算规避后旋转角', 'Rotation_angle', 'R-只读', 'Number', '度', '斜单轴设备数据', '[]', '2025-02-06 12:11:35.969866+00', '2025-02-06 12:11:35.969866+00', NULL, 'd616bcbb'),
	('6a20f4fb-4e3d-0064-b1f2-6a19d41a4267', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备方位角(旋转角)', 'Device_a', 'R-只读', 'Number', '', '', '[]', '2025-02-06 12:17:59.888991+00', '2025-02-06 12:17:59.888991+00', NULL, 'd616bcbb'),
	('d93a1551-7c67-1b76-186a-7b369b6f4da0', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备高度角', 'Device_h', 'R-只读', 'Number', '度', '', '[]', '2025-02-06 12:18:14.206224+00', '2025-02-06 12:18:14.206224+00', NULL, 'd616bcbb'),
	('ae3b43d1-7fa9-45c1-dee6-e7605737a7c8', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '高度电机电流', 'Motor_current_h', 'R-只读', 'Number', '毫安', '', '[]', '2025-02-06 12:20:06.347953+00', '2025-02-06 12:20:06.347953+00', NULL, 'd616bcbb'),
	('81477167-006d-5afa-23d4-ad32bd7c1de2', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '方位角（旋转角）电机电流获取延时', 'A_delay', 'RW-读/写', 'Number', '毫秒', '延迟获取设备方位角电机电流 ', '[]', '2025-02-06 12:07:51.336406+00', '2025-02-06 12:20:31.062373+00', NULL, 'd616bcbb'),
	('2d494b67-0aec-57ba-c33b-27ea39b7b2cf', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '方位电机（旋转角）电流', 'Motor_current_a', 'R-只读', 'Number', '毫安', '', '[]', '2025-02-06 12:25:22.779992+00', '2025-02-06 12:25:22.779992+00', NULL, 'd616bcbb'),
	('02731535-eed5-294e-153c-44dd7790c0bf', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '运行状态', 'Device_status', 'R-只读', 'Number', '', '0	停止运行
1	自动运行
2	指定位置运行
3	模拟运行
4	防风模式
5	防雪模式
6	防冰雹模式
7	清洗模式（大雨天气）
', '[]', '2025-02-06 12:25:53.688864+00', '2025-02-06 12:25:53.688864+00', NULL, 'd616bcbb'),
	('b5208d26-be1f-4ff2-62f0-33c340011911', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '故障码', 'Device_error', 'R-只读', 'Number', '', '0	无故障
101	温度读取失败
102	电压异常（低于正常值）
103	mesh组网失败
104	倾角传感器错误
105	编码器错误
106	超出限制北范围
107	超出限制南范围
108	超出限制西范围
109	超出限制东范围
110	风速传感器错误
111	高度角电机电流超出
112	方位角电机电流超出
115	编码器，倾角传感器全部通信错误
116	gps获取错误
117	4g通信失败
118	ota升级失败
', '[]', '2025-02-06 12:26:22.59962+00', '2025-02-06 12:26:22.59962+00', NULL, 'd616bcbb'),
	('db4e6c38-dbce-7b00-13a7-b3c2fa481a4a', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '回位设备高度角', 'Return_device_altitude', 'RW-读/写', 'Number', '', '归位状态设备运行到的高度角', '[]', '2025-02-06 12:28:06.849402+00', '2025-02-06 12:28:06.849402+00', NULL, 'd616bcbb'),
	('024463aa-9790-c80e-5883-39976e8cbae1', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '指定方位角', 'Designated_azimuth', 'RW-读/写', 'Number', '', '', '[]', '2025-02-06 12:29:11.124833+00', '2025-02-06 12:29:11.124833+00', NULL, 'd616bcbb'),
	('a08946a6-9d92-60e9-88a9-9a7576fc2dfe', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '指定高度角', 'Specifies_height', 'RW-读/写', 'Number', '', '', '[]', '2025-02-06 12:29:24.366288+00', '2025-02-06 12:29:24.366288+00', NULL, 'd616bcbb'),
	('00da271e-b82f-2adc-beec-257bc327f93e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', 'ota升级标志位', 'ota_flag', 'R-只读', 'Number', '', '0（未升级），2（正在升级），3（升级完成）,4(不需要ota升级)，5（升级失败）
', '[]', '2025-02-06 12:29:42.407495+00', '2025-02-06 12:29:42.407495+00', NULL, 'd616bcbb'),
	('7e8d2b86-e339-4cf9-666d-ca374578d8e7', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '系统版本', 'System_version', 'R-只读', 'Number', '', '固件版本号', '[]', '2025-02-06 12:30:00.354424+00', '2025-02-06 12:30:00.354424+00', NULL, 'd616bcbb'),
	('651fc2d5-d073-0961-1d62-012973242997', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备处于mesh网络层数', 'mesh_layer', 'R-只读', 'Number', '', '', '[]', '2025-02-06 12:31:36.805079+00', '2025-02-06 12:31:36.805079+00', NULL, 'd616bcbb'),
	('dd42e452-ad73-f093-9681-c7e7aa466180', '4c2f6999-630f-a358-bb7d-8eb5130a502c', 'mesh网络id', 'mesh_id', 'R-只读', 'Number', '', '', '[]', '2025-02-06 12:32:14.780132+00', '2025-02-06 12:32:14.780132+00', NULL, 'd616bcbb'),
	('d852f73e-ce5b-4c59-377f-9d82c02fe31e', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '根节点mac地址', 'is_root_mac', 'R-只读', 'String', '', '', '[]', '2025-02-06 12:31:48.913624+00', '2025-02-06 12:32:29.207798+00', NULL, 'd616bcbb'),
	('5a1db7dd-8145-6f20-e18a-9a533bd32978', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '是否处于mesh网络', 'is_mesh', 'R-只读', 'Number', '', '0（不是），1（是）
', '[]', '2025-02-06 12:32:51.024252+00', '2025-02-06 12:32:51.024252+00', NULL, 'd616bcbb'),
	('61add2ca-c8ec-5b90-ae09-8eafa39cf6f0', '4c2f6999-630f-a358-bb7d-8eb5130a502c', '设备是否为根节点', 'is_root', 'R-只读', 'Number', '', '0（不是），1（是）
', '[]', '2025-02-06 12:32:00.611266+00', '2025-02-06 12:32:57.49769+00', NULL, 'd616bcbb'),
	('51419432-f0e4-ef43-40b0-ad3e8442be33', 'd391b336-4273-d101-27f8-d1fbedfde866', '风速', 'wind_speed', 'R-只读', 'Number', 'm/s', '', '[]', '2025-02-09 03:23:38.771712+00', '2025-02-09 03:23:38.771712+00', NULL, 'd616bcbb'),
	('3df3a65a-7847-6165-fc71-dd0fc3c9b9a5', 'd391b336-4273-d101-27f8-d1fbedfde866', '风向', 'wind_direction', 'R-只读', 'Number', '度', '实际值（正北方向为0°顺时针增加度数，正东方为90°）', '[]', '2025-02-09 03:25:06.6538+00', '2025-02-09 03:25:15.28915+00', NULL, 'd616bcbb'),
	('a0a3407a-2ba0-1734-be98-359896aa7614', 'd391b336-4273-d101-27f8-d1fbedfde866', '温度', 'temperature', 'R-只读', 'Number', '摄氏度', '', '[]', '2025-02-09 03:26:25.250051+00', '2025-02-09 11:32:04.599574+00', NULL, 'd616bcbb'),
	('dec20597-f044-6ee5-a250-cfaa0bf94bc8', 'd391b336-4273-d101-27f8-d1fbedfde866', '湿度', 'humidity', 'R-只读', 'Number', '%', '', '[]', '2025-02-09 03:25:53.742736+00', '2025-02-09 11:31:46.882417+00', NULL, 'd616bcbb');
/*!40000 ALTER TABLE "device_model_telemetry" ENABLE KEYS */;

-- 导出  表 public.device_templates 结构
CREATE TABLE IF NOT EXISTS "device_templates" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"author" VARCHAR(36) NULL DEFAULT '',
	"version" VARCHAR(50) NULL DEFAULT '',
	"description" VARCHAR(500) NULL DEFAULT '',
	"tenant_id" VARCHAR(36) NOT NULL DEFAULT '',
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"flag" SMALLINT NULL DEFAULT '1',
	"label" VARCHAR(255) NULL DEFAULT NULL,
	"web_chart_config" JSON NULL DEFAULT NULL,
	"app_chart_config" JSON NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"path" VARCHAR(999) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.device_templates 的数据：2 rows
/*!40000 ALTER TABLE "device_templates" DISABLE KEYS */;
REPLACE INTO "device_templates" ("id", "name", "author", "version", "description", "tenant_id", "created_at", "updated_at", "flag", "label", "web_chart_config", "app_chart_config", "remark", "path") VALUES
	('d391b336-4273-d101-27f8-d1fbedfde866', '气象监控站', 'lxc', '1.0.0', '气象传感器', 'd616bcbb', '2025-02-08 08:34:07.158453+00', '2025-02-12 08:16:41.299385+00', 1, '', '[]', '[]', NULL, ''),
	('4c2f6999-630f-a358-bb7d-8eb5130a502c', 'ESP32物模型', 'lxc', 'V1.0', 'ESP32设备物模型', 'd616bcbb', '2025-02-04 04:31:13.266078+00', '2025-02-13 14:28:40.571135+00', 1, '', '[]', '[]', NULL, '');
/*!40000 ALTER TABLE "device_templates" ENABLE KEYS */;

-- 导出  表 public.device_trigger_condition 结构
CREATE TABLE IF NOT EXISTS "device_trigger_condition" (
	"id" VARCHAR(36) NOT NULL,
	"scene_automation_id" VARCHAR(36) NOT NULL,
	"enabled" VARCHAR(10) NOT NULL,
	"group_id" VARCHAR(36) NOT NULL,
	"trigger_condition_type" VARCHAR(10) NOT NULL,
	"trigger_source" VARCHAR(36) NULL DEFAULT NULL,
	"trigger_param_type" VARCHAR(10) NULL DEFAULT NULL,
	"trigger_param" VARCHAR(50) NULL DEFAULT NULL,
	"trigger_operator" VARCHAR(10) NULL DEFAULT NULL,
	"trigger_value" VARCHAR(99) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "fk_scene_automation_id" FOREIGN KEY ("scene_automation_id") REFERENCES "scene_automations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.device_trigger_condition 的数据：-1 rows
/*!40000 ALTER TABLE "device_trigger_condition" DISABLE KEYS */;
/*!40000 ALTER TABLE "device_trigger_condition" ENABLE KEYS */;

-- 导出  表 public.device_user_logs 结构
CREATE TABLE IF NOT EXISTS "device_user_logs" (
	"id" VARCHAR(36) NOT NULL,
	"device_nums" INTEGER NOT NULL DEFAULT '0',
	"device_on" INTEGER NOT NULL DEFAULT '0',
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.device_user_logs 的数据：-1 rows
/*!40000 ALTER TABLE "device_user_logs" DISABLE KEYS */;
/*!40000 ALTER TABLE "device_user_logs" ENABLE KEYS */;

-- 导出  表 public.event_datas 结构
CREATE TABLE IF NOT EXISTS "event_datas" (
	"id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"identify" VARCHAR(255) NOT NULL,
	"ts" TIMESTAMPTZ NOT NULL,
	"data" JSON NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "event_datas_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.event_datas 的数据：-1 rows
/*!40000 ALTER TABLE "event_datas" DISABLE KEYS */;
REPLACE INTO "event_datas" ("id", "device_id", "identify", "ts", "data", "tenant_id") VALUES
	('ce9d21b7-d639-b4d5-0555-4af651b98748', '112233445566', 'Equipment_failure', '2025-02-13 14:27:45.904039+00', '{"Device_error":2}', 'd616bcbb'),
	('6bca88fc-2067-6024-f4c3-0b8fdc3a2209', '112233445566', 'Equipment_failure', '2025-02-13 14:27:47.659845+00', '{"Device_error":2}', 'd616bcbb'),
	('cba82f98-0eb5-f407-f7f9-11448bf0a16e', '112233445566', 'Equipment_failure', '2025-02-13 14:27:48.657687+00', '{"Device_error":2}', 'd616bcbb');
/*!40000 ALTER TABLE "event_datas" ENABLE KEYS */;

-- 导出  表 public.expected_datas 结构
CREATE TABLE IF NOT EXISTS "expected_datas" (
	"id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"send_type" VARCHAR(50) NOT NULL,
	"payload" JSONB NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"send_time" TIMESTAMPTZ NULL DEFAULT NULL,
	"status" VARCHAR(50) NOT NULL DEFAULT 'pending',
	"message" TEXT NULL DEFAULT NULL,
	"expiry_time" TIMESTAMPTZ NULL DEFAULT NULL,
	"label" VARCHAR(100) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "expected_datas_devices_fk" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

-- 正在导出表  public.expected_datas 的数据：-1 rows
/*!40000 ALTER TABLE "expected_datas" DISABLE KEYS */;
/*!40000 ALTER TABLE "expected_datas" ENABLE KEYS */;

-- 导出  表 public.groups 结构
CREATE TABLE IF NOT EXISTS "groups" (
	"id" VARCHAR(36) NOT NULL,
	"parent_id" VARCHAR(36) NULL DEFAULT '0',
	"tier" INTEGER NOT NULL DEFAULT '1',
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.groups 的数据：-1 rows
/*!40000 ALTER TABLE "groups" DISABLE KEYS */;
REPLACE INTO "groups" ("id", "parent_id", "tier", "name", "description", "created_at", "updated_at", "remark", "tenant_id") VALUES
	('59c539ac-78d0-7303-8dbf-22b86dfe1087', '0', -1, '北京', '', '2025-02-13 15:06:03.01649+00', '2025-02-13 15:06:03.01649+00', NULL, 'd616bcbb'),
	('7553a6f8-dcf1-8ccf-f0f0-4c3d1e6f728a', '0', -1, '上海', '', '2025-02-13 15:06:14.890003+00', '2025-02-13 15:06:14.890003+00', NULL, 'd616bcbb');
/*!40000 ALTER TABLE "groups" ENABLE KEYS */;

-- 导出  表 public.logo 结构
CREATE TABLE IF NOT EXISTS "logo" (
	"id" VARCHAR(36) NOT NULL,
	"system_name" VARCHAR(99) NOT NULL,
	"logo_cache" VARCHAR(255) NOT NULL,
	"logo_background" VARCHAR(255) NOT NULL,
	"logo_loading" VARCHAR(255) NOT NULL,
	"home_background" VARCHAR(255) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.logo 的数据：-1 rows
/*!40000 ALTER TABLE "logo" DISABLE KEYS */;
REPLACE INTO "logo" ("id", "system_name", "logo_cache", "logo_background", "logo_loading", "home_background", "remark") VALUES
	('a', 'ThingsPanel', '', '', '', '', NULL);
/*!40000 ALTER TABLE "logo" ENABLE KEYS */;

-- 导出  表 public.notification_groups 结构
CREATE TABLE IF NOT EXISTS "notification_groups" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(99) NOT NULL,
	"notification_type" VARCHAR(25) NOT NULL,
	"status" VARCHAR(10) NOT NULL,
	"notification_config" JSONB NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.notification_groups 的数据：-1 rows
/*!40000 ALTER TABLE "notification_groups" DISABLE KEYS */;
/*!40000 ALTER TABLE "notification_groups" ENABLE KEYS */;

-- 导出  表 public.notification_histories 结构
CREATE TABLE IF NOT EXISTS "notification_histories" (
	"id" VARCHAR(36) NOT NULL,
	"send_time" TIMESTAMPTZ NOT NULL,
	"send_content" TEXT NULL DEFAULT NULL,
	"send_target" VARCHAR(255) NOT NULL,
	"send_result" VARCHAR(25) NULL DEFAULT NULL,
	"notification_type" VARCHAR(25) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.notification_histories 的数据：-1 rows
/*!40000 ALTER TABLE "notification_histories" DISABLE KEYS */;
/*!40000 ALTER TABLE "notification_histories" ENABLE KEYS */;

-- 导出  表 public.notification_services_config 结构
CREATE TABLE IF NOT EXISTS "notification_services_config" (
	"id" VARCHAR(36) NOT NULL,
	"config" JSON NULL DEFAULT NULL,
	"notice_type" VARCHAR(36) NOT NULL,
	"status" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.notification_services_config 的数据：-1 rows
/*!40000 ALTER TABLE "notification_services_config" DISABLE KEYS */;
/*!40000 ALTER TABLE "notification_services_config" ENABLE KEYS */;

-- 导出  表 public.one_time_tasks 结构
CREATE TABLE IF NOT EXISTS "one_time_tasks" (
	"id" VARCHAR(36) NOT NULL,
	"scene_automation_id" VARCHAR(36) NOT NULL,
	"execution_time" TIMESTAMPTZ NOT NULL,
	"executing_state" VARCHAR(10) NOT NULL,
	"enabled" VARCHAR(10) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"expiration_time" BIGINT NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "fk_scene_automation_id" FOREIGN KEY ("scene_automation_id") REFERENCES "scene_automations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.one_time_tasks 的数据：-1 rows
/*!40000 ALTER TABLE "one_time_tasks" DISABLE KEYS */;
/*!40000 ALTER TABLE "one_time_tasks" ENABLE KEYS */;

-- 导出  表 public.operation_logs 结构
CREATE TABLE IF NOT EXISTS "operation_logs" (
	"id" VARCHAR(36) NOT NULL,
	"ip" VARCHAR(36) NOT NULL,
	"path" VARCHAR(2000) NULL DEFAULT NULL,
	"user_id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"latency" BIGINT NULL DEFAULT NULL,
	"request_message" TEXT NULL DEFAULT NULL,
	"response_message" TEXT NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.operation_logs 的数据：-1 rows
/*!40000 ALTER TABLE "operation_logs" DISABLE KEYS */;
REPLACE INTO "operation_logs" ("id", "ip", "path", "user_id", "name", "created_at", "latency", "request_message", "response_message", "tenant_id", "remark") VALUES
	('551f3fcc-dfa3-1edc-562d-919926d65d32', '127.0.0.1', '/api/v1/board/user/update/password', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:00:40.208322+00', 110, '{"old_password":"123456","password":"Qwer1234@"}', '', 'd616bcbb', NULL),
	('4ab90477-ca64-1cb4-beb9-e05fc8902dbd', '127.0.0.1', '/api/v1/device/template', '11111111-4fe9-b409-67c3-111111111111', 'PUT', '2025-02-13 14:01:05.46992+00', 7, '{"name":"ESP32物模型","templateTage":[],"version":"V1.0","author":"lxc","description":"ESP32设备物模型","path":"","label":"","id":"4c2f6999-630f-a358-bb7d-8eb5130a502c"}', '', 'd616bcbb', NULL),
	('c65a33f0-5742-0768-c4cb-8fb8dbefde88', '127.0.0.1', '/api/v1/device', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:01:21.126755+00', 10, '{"name":"ESP32","label":"","device_config_id":"315d9d82-5c76-3197-4eab-8c0a641ccdc9","access_way":"A"}', '', 'd616bcbb', NULL),
	('25f19cb0-38ac-0475-8b83-05f84fdb3f00', '127.0.0.1', '/api/v1/device/update/voucher', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:01:23.002635+00', 11, '{"device_id":"73738433-8322-495c-b9b7-364ffec1e56a","voucher":"{\"username\":\"dc2f9d45-cda6-526d-442\",\"password\":\"e08a2c1\"}"}', '', 'd616bcbb', NULL),
	('ff5dd5de-1302-8549-5f20-08fce34e99e9', '127.0.0.1', '/api/v1/device', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:01:37.926362+00', 8, '{"name":"气象站","label":"","device_config_id":"964d6220-ecbf-a043-1960-85b1a2758cea","access_way":"A"}', '', 'd616bcbb', NULL),
	('e929eb80-e8fd-640c-ef5d-65c4245de237', '127.0.0.1', '/api/v1/device/update/voucher', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:01:40.17577+00', 9, '{"device_id":"8a09d81c-3d13-f159-f968-ec69140da60d","voucher":"{\"username\":\"27ae0698-99ac-33b4-4ec\",\"password\":\"\"}"}', '', 'd616bcbb', NULL),
	('9c9190d5-f3c6-a233-d5e8-2e00656403f9', '127.0.0.1', '/api/v1/device/73738433-8322-495c-b9b7-364ffec1e56a', '11111111-4fe9-b409-67c3-111111111111', 'DELETE', '2025-02-13 14:10:02.202739+00', 51, '', '', 'd616bcbb', NULL),
	('c079b62f-0806-d2a3-143c-cb312105eede', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:16:38.745322+00', 2, '{"device_id":"112233445566","value":"{\r\n\"Reboot\":1,\r\n}","identify":"Device_restart"}', '', 'd616bcbb', NULL),
	('063c3360-080d-c34d-5311-2ee6ffa102ef', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:18:14.888239+00', 3, '{"device_id":"112233445566","value":"{\r\n\"Reboot\":1,\r\n}","identify":"Device_restart"}', '', 'd616bcbb', NULL),
	('7705ee35-7600-361b-1eb9-30f1351c46aa', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:19:04.276878+00', 3, '{"device_id":"112233445566","value":"{\r\n\"Reboot\":1,\r\n}","identify":"Device_restart"}', '', 'd616bcbb', NULL),
	('8a2e7589-ef07-5294-cfab-c54829b0b462', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:19:35.793509+00', 2, '{"device_id":"112233445566","value":"{\r\n\"Reboot\":1,\r\n}","identify":"Device_restart"}', '', 'd616bcbb', NULL),
	('cd39dee1-8e5c-b1d4-802f-3e87533cf380', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:19:59.796894+00', 13, '{"device_id":"112233445566","value":"{\r\n\"Reboot\":1\r\n}","identify":"Device_restart"}', '', 'd616bcbb', NULL),
	('2fb781b7-7257-e9d0-6818-51e9192d6c83', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:22:34.433284+00', 20, '{"device_id":"112233445566","value":"{\r\n\"Reboot\":1\r\n}","identify":"Device_restart"}', '', 'd616bcbb', NULL),
	('4e49b98e-72fd-458c-4cd0-4a808fe47c02', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:22:40.190012+00', 10, '{"device_id":"112233445566","value":"{\r\n\"Mode\":1\r\n}","identify":"Runs_automatically"}', '', 'd616bcbb', NULL),
	('0c44aa03-caa0-8f26-e743-9cfa5475155c', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:22:46.87592+00', 9, '{"device_id":"112233445566","value":"{\r\n\"sntp\":1\r\n}","identify":"Device_calibration_time"}', '', 'd616bcbb', NULL),
	('25d85207-b92f-e0ba-2dc4-c73e71df5bd5', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:22:52.084411+00', 10, '{"device_id":"112233445566","value":"{\r\n\"Mode\":0\r\n\r\n}","identify":"stop"}', '', 'd616bcbb', NULL),
	('e170a694-2668-ba22-7ab1-657d33c65542', '127.0.0.1', '/api/v1/command/datas/pub', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 14:23:22.756999+00', 10, '{"device_id":"112233445566","value":"{\"Mode\":1,\"Designated_azimuth\":52,\"Specifies_height\":86}","identify":"set_mode"}', '', 'd616bcbb', NULL),
	('5c91d3b5-a3c9-88d6-a089-b2b27b795630', '127.0.0.1', '/api/v1/device/template', '11111111-4fe9-b409-67c3-111111111111', 'PUT', '2025-02-13 14:25:59.90818+00', 7, '{"name":"ESP32物模型","templateTage":[],"version":"V1.0","author":"lxc","description":"ESP32设备物模型","path":"","label":"","id":"4c2f6999-630f-a358-bb7d-8eb5130a502c"}', '', 'd616bcbb', NULL),
	('01796c5d-c139-aea3-c858-aaa78d293174', '127.0.0.1', '/api/v1/device/template', '11111111-4fe9-b409-67c3-111111111111', 'PUT', '2025-02-13 14:28:40.568137+00', 7, '{"name":"ESP32物模型","templateTage":[],"version":"V1.0","author":"lxc","description":"ESP32设备物模型","path":"","label":"","id":"4c2f6999-630f-a358-bb7d-8eb5130a502c"}', '', 'd616bcbb', NULL),
	('4453b65d-a185-963b-8ce7-cd0261ea5ede', '127.0.0.1', '/api/v1/device_config', '11111111-4fe9-b409-67c3-111111111111', 'PUT', '2025-02-13 14:29:37.350261+00', 9, '{"id":"315d9d82-5c76-3197-4eab-8c0a641ccdc9","other_config":"{\"online_timeout\":10,\"heartbeat\":0}"}', '', 'd616bcbb', NULL),
	('5bb8fad2-ce14-3b70-ffe1-90022e624200', '127.0.0.1', '/api/v1/ui_elements', '00000000-4fe9-b409-67c3-000000000000', 'PUT', '2025-02-13 15:04:21.324508+00', 7, '{"parent_id":"0","element_code":"visualization","param1":"/visualization","multilingual":"route.visualization","param2":"icon-park-outline:data-server","param3":"self","orders":113,"description":"可视化","element_type":1,"authority":"[\"SYS_ADMIN\"]","route_path":"layout.base","remark":"","id":"95e2a961-382b-f4a6-87b3-1898123c95bc","children":[{"id":"a2654c98-3749-c88b-0472-b414049ca532","parent_id":"95e2a961-382b-f4a6-87b3-1898123c95bc","element_code":"route.visualization_kanban","element_type":3,"orders":1131,"param1":"/visualization/kanban","param2":"tabler:device-tv","param3":"self","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"看板","remark":"","multilingual":"route.visualization_kanban","route_path":"view.visualization_panel","children":[]},{"id":"ed4a5cfa-03e7-ccc0-6cc8-bcadccd25541","parent_id":"95e2a961-382b-f4a6-87b3-1898123c95bc","element_code":"visualization_kanban-details","element_type":3,"orders":1132,"param1":"/visualization/kanban-details","param2":"ic:baseline-credit-card","param3":"1","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"看板详情","remark":"","multilingual":"看板详情","route_path":"view.visualization_panel-details","children":[]},{"id":"502a0d6c-750e-92f6-a1a7-ffdd362dbbac","parent_id":"95e2a961-382b-f4a6-87b3-1898123c95bc","element_code":"visualization_panel-preview","element_type":3,"orders":1133,"param1":"/visualization/panel-preview","param2":"","param3":"1","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"看板预览","remark":"","multilingual":"route.visualization_panel-preview","route_path":"view.visualization_panel-preview","children":[]}]}', '', 'aaaaaa', NULL),
	('05bfe680-5eab-9518-2691-6ee728a80bdc', '127.0.0.1', '/api/v1/ui_elements', '00000000-4fe9-b409-67c3-000000000000', 'PUT', '2025-02-13 15:04:28.411147+00', 4, '{"parent_id":"0","element_code":"automation","param1":"/automation","multilingual":"route.automation","param2":"material-symbols:device-hub","param3":"self","orders":114,"description":"自动化","element_type":1,"authority":"[\"SYS_ADMIN\"]","route_path":"layout.base","remark":"","id":"676e8f33-875a-0473-e9ca-c82fd09fef57","children":[{"id":"975c9550-5db9-7b4c-5dea-7a4c326a37ff","parent_id":"676e8f33-875a-0473-e9ca-c82fd09fef57","element_code":"automation_scene-edit","element_type":3,"orders":1,"param1":"/automation/scene-edit","param2":"mdi:apps-box","param3":"1","authority":["TENANT_ADMIN"],"description":"新增场景","remark":"","multilingual":"route.automation_scene-edit","route_path":"view.automation_scene-edit","children":[]},{"id":"64f684f1-390c-b5f2-9994-36895025df8a","parent_id":"676e8f33-875a-0473-e9ca-c82fd09fef57","element_code":"automation_space-management","element_type":3,"orders":10,"param1":"automation/space-management","param2":"ic:baseline-security","param3":"1","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"场景管理","remark":"","multilingual":"default","route_path":"view.automation space-management","children":[]},{"id":"01dab674-9556-cdd7-b800-78bcb366adb4","parent_id":"676e8f33-875a-0473-e9ca-c82fd09fef57","element_code":"automation_scene-linkage","element_type":3,"orders":1141,"param1":"/automation/scene-linkage","param2":"mdi:airplane-edit","param3":"self","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"场景联动","remark":"","multilingual":"route.automation_scene-linkage","route_path":"view.automation_scene-linkage","children":[]},{"id":"51381989-1160-93cd-182e-d44a1c4ab89b","parent_id":"676e8f33-875a-0473-e9ca-c82fd09fef57","element_code":"automation_scene-manage","element_type":3,"orders":1142,"param1":"/automation/scene-manage","param2":"uil:brightness-plus","param3":"self","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"场景管理","remark":"","multilingual":"route.automation_scene-manage","route_path":"view.automation_scene-manage","children":[]},{"id":"96aa2fac-90b2-aca1-1ce0-51b5060f4081","parent_id":"676e8f33-875a-0473-e9ca-c82fd09fef57","element_code":"automation_linkage-edit","element_type":3,"orders":1143,"param1":"/automation/linkage-edit","param2":"","param3":"1","authority":["TENANT_ADMIN","SYS_ADMIN"],"description":"场景联动编辑","remark":"","multilingual":"route.automation_linkage-edit","route_path":"view.automation_linkage-edit","children":[]}]}', '', 'aaaaaa', NULL),
	('023cc59c-911a-f588-b61a-e2217b035808', '127.0.0.1', '/api/v1/ui_elements', '00000000-4fe9-b409-67c3-000000000000', 'PUT', '2025-02-13 15:04:33.852208+00', 4, '{"parent_id":"0","element_code":"alarm","param1":"/alarm","multilingual":"route.alarm","param2":"mdi:alert","param3":"self","orders":115,"description":"告警","element_type":1,"authority":"[\"SYS_ADMIN\"]","route_path":"layout.base","remark":"","id":"650bc444-7672-1123-1e41-7e37365b0186","children":[{"id":"c078182f-bf4b-b560-da97-02926fa98f78","parent_id":"650bc444-7672-1123-1e41-7e37365b0186","element_code":"alarm_notification-record","element_type":3,"orders":1,"param1":"/alarm/notification-record","param2":"icon-park-outline:editor","param3":"self","authority":["TENANT_ADMIN"],"description":"通知记录","remark":"","multilingual":"route.alarm_notification-record","route_path":"view.alarm_notification-record","children":[]},{"id":"82c46beb-9ec4-8a3d-c6e4-04ba426e525a","parent_id":"650bc444-7672-1123-1e41-7e37365b0186","element_code":"alarm_notification-group","element_type":3,"orders":1,"param1":"/alarm/notification-group","param2":"ic:round-supervisor-account","param3":"basic","authority":["TENANT_ADMIN"],"description":"通知组","remark":"","multilingual":"route.alarm_notification-group","route_path":"view.alarm_notification-group","children":[]},{"id":"485c2a20-ebc5-2216-4871-26453470d290","parent_id":"650bc444-7672-1123-1e41-7e37365b0186","element_code":"alarm_warning-message","element_type":3,"orders":999,"param1":"/alarm/warning-message","param2":"mdi:airballoon","param3":"self","authority":["TENANT_ADMIN"],"description":"警告信息","remark":"","multilingual":"route.alarm_warning-message","route_path":"view.alarm_warning-message","children":[]}]}', '', 'aaaaaa', NULL),
	('92e72237-444f-124d-8f7a-aa731300b924', '127.0.0.1', '/api/v1/ui_elements', '00000000-4fe9-b409-67c3-000000000000', 'PUT', '2025-02-13 15:05:27.2323+00', 6, '{"parent_id":"5373a6a2-1861-af35-eb4c-adfd5ca55ecd","element_code":"device_service-details","param1":"/device/service-details","multilingual":"route.device_service_details","param2":"ph:align-bottom","param3":"1","orders":1130,"description":"服务详情","element_type":3,"authority":"[\"SYS_ADMIN\"]","route_path":"","remark":"","id":"f960c45c-6d5b-e67a-c4ff-1f0e869c1625","children":[]}', '', 'aaaaaa', NULL),
	('cf5485d9-3d69-8a44-502e-af470270eb4f', '127.0.0.1', '/api/v1/ui_elements', '00000000-4fe9-b409-67c3-000000000000', 'PUT', '2025-02-13 15:05:32.422134+00', 4, '{"parent_id":"5373a6a2-1861-af35-eb4c-adfd5ca55ecd","element_code":"device_service-access","param1":"/device/service-access","multilingual":"route.device_service_access","param2":"mdi:ab-testing","param3":"0","orders":1129,"description":"服务接入点管理","element_type":3,"authority":"[\"SYS_ADMIN\"]","route_path":"","remark":"","id":"075d9f19-5618-bb9b-6ccd-f382bfd3292b","children":[]}', '', 'aaaaaa', NULL),
	('85462e87-4f32-bb49-e419-2e98ee6ef2fd', '127.0.0.1', '/api/v1/device/group', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 15:06:03.01649+00', 8, '{"id":"","parent_id":"0","name":"北京","description":""}', '', 'd616bcbb', NULL),
	('2462afc7-7b4d-9d4d-49eb-0bf2e222fe37', '127.0.0.1', '/api/v1/device/group', '11111111-4fe9-b409-67c3-111111111111', 'POST', '2025-02-13 15:06:14.890003+00', 5, '{"id":"","parent_id":"0","name":"上海","description":""}', '', 'd616bcbb', NULL);
/*!40000 ALTER TABLE "operation_logs" ENABLE KEYS */;

-- 导出  表 public.ota_upgrade_packages 结构
CREATE TABLE IF NOT EXISTS "ota_upgrade_packages" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(200) NOT NULL,
	"version" VARCHAR(36) NOT NULL,
	"target_version" VARCHAR(36) NULL DEFAULT NULL,
	"device_config_id" VARCHAR(36) NOT NULL,
	"module" VARCHAR(36) NULL DEFAULT NULL,
	"package_type" SMALLINT NOT NULL,
	"signature_type" VARCHAR(36) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT '{}',
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"package_url" VARCHAR(500) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"signature" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.ota_upgrade_packages 的数据：-1 rows
/*!40000 ALTER TABLE "ota_upgrade_packages" DISABLE KEYS */;
/*!40000 ALTER TABLE "ota_upgrade_packages" ENABLE KEYS */;

-- 导出  表 public.ota_upgrade_tasks 结构
CREATE TABLE IF NOT EXISTS "ota_upgrade_tasks" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(200) NOT NULL,
	"ota_upgrade_package_id" VARCHAR(36) NOT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.ota_upgrade_tasks 的数据：-1 rows
/*!40000 ALTER TABLE "ota_upgrade_tasks" DISABLE KEYS */;
/*!40000 ALTER TABLE "ota_upgrade_tasks" ENABLE KEYS */;

-- 导出  表 public.ota_upgrade_task_details 结构
CREATE TABLE IF NOT EXISTS "ota_upgrade_task_details" (
	"id" VARCHAR(36) NOT NULL,
	"ota_upgrade_task_id" VARCHAR(200) NOT NULL,
	"device_id" VARCHAR(200) NOT NULL,
	"steps" SMALLINT NULL DEFAULT NULL,
	"status" SMALLINT NOT NULL,
	"status_description" VARCHAR(500) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "fk_ota_upgrade_tasks" FOREIGN KEY ("ota_upgrade_task_id") REFERENCES "ota_upgrade_tasks" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT "ota_upgrade_task_details_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.ota_upgrade_task_details 的数据：-1 rows
/*!40000 ALTER TABLE "ota_upgrade_task_details" DISABLE KEYS */;
/*!40000 ALTER TABLE "ota_upgrade_task_details" ENABLE KEYS */;

-- 导出  表 public.periodic_tasks 结构
CREATE TABLE IF NOT EXISTS "periodic_tasks" (
	"id" VARCHAR(36) NOT NULL,
	"scene_automation_id" VARCHAR(36) NOT NULL,
	"task_type" VARCHAR(255) NOT NULL,
	"params" VARCHAR(50) NOT NULL,
	"execution_time" TIMESTAMPTZ NOT NULL,
	"enabled" VARCHAR(10) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"expiration_time" BIGINT NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "scene_automation_id_fkey" FOREIGN KEY ("scene_automation_id") REFERENCES "scene_automations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.periodic_tasks 的数据：-1 rows
/*!40000 ALTER TABLE "periodic_tasks" DISABLE KEYS */;
/*!40000 ALTER TABLE "periodic_tasks" ENABLE KEYS */;

-- 导出  表 public.products 结构
CREATE TABLE IF NOT EXISTS "products" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"product_type" VARCHAR(36) NULL DEFAULT NULL,
	"product_key" VARCHAR(255) NULL DEFAULT NULL,
	"product_model" VARCHAR(100) NULL DEFAULT NULL,
	"image_url" VARCHAR(500) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(500) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	"device_config_id" VARCHAR(36) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "products_device_configs_fk" FOREIGN KEY ("device_config_id") REFERENCES "device_configs" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.products 的数据：-1 rows
/*!40000 ALTER TABLE "products" DISABLE KEYS */;
/*!40000 ALTER TABLE "products" ENABLE KEYS */;

-- 导出  表 public.protocol_plugins 结构
CREATE TABLE IF NOT EXISTS "protocol_plugins" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(36) NOT NULL,
	"device_type" SMALLINT NOT NULL DEFAULT '1',
	"protocol_type" VARCHAR(50) NOT NULL,
	"access_address" VARCHAR(500) NULL DEFAULT NULL,
	"http_address" VARCHAR(500) NULL DEFAULT NULL,
	"sub_topic_prefix" VARCHAR(500) NULL DEFAULT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"additional_info" VARCHAR(1000) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"update_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.protocol_plugins 的数据：-1 rows
/*!40000 ALTER TABLE "protocol_plugins" DISABLE KEYS */;
/*!40000 ALTER TABLE "protocol_plugins" ENABLE KEYS */;

-- 导出  表 public.roles 结构
CREATE TABLE IF NOT EXISTS "roles" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(99) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"created_at" TIMESTAMP NULL DEFAULT NULL,
	"updated_at" TIMESTAMP NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.roles 的数据：-1 rows
/*!40000 ALTER TABLE "roles" DISABLE KEYS */;
/*!40000 ALTER TABLE "roles" ENABLE KEYS */;

-- 导出  表 public.r_group_device 结构
CREATE TABLE IF NOT EXISTS "r_group_device" (
	"group_id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	UNIQUE INDEX "r_group_device_group_id_device_id_key" ("group_id", "device_id"),
	CONSTRAINT "fk_group_device" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT "fk_group_device_2" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.r_group_device 的数据：-1 rows
/*!40000 ALTER TABLE "r_group_device" DISABLE KEYS */;
/*!40000 ALTER TABLE "r_group_device" ENABLE KEYS */;

-- 导出  表 public.scene_action_info 结构
CREATE TABLE IF NOT EXISTS "scene_action_info" (
	"id" VARCHAR(36) NOT NULL,
	"scene_id" VARCHAR(36) NOT NULL,
	"action_target" VARCHAR(36) NOT NULL,
	"action_type" VARCHAR(10) NOT NULL,
	"action_param_type" VARCHAR(10) NULL DEFAULT NULL,
	"action_param" VARCHAR(50) NULL DEFAULT NULL,
	"action_value" VARCHAR(255) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "scene_action_info_scene_id_fkey" FOREIGN KEY ("scene_id") REFERENCES "scene_info" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.scene_action_info 的数据：-1 rows
/*!40000 ALTER TABLE "scene_action_info" DISABLE KEYS */;
/*!40000 ALTER TABLE "scene_action_info" ENABLE KEYS */;

-- 导出  表 public.scene_automations 结构
CREATE TABLE IF NOT EXISTS "scene_automations" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"enabled" VARCHAR(10) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"creator" VARCHAR(36) NOT NULL,
	"updator" VARCHAR(36) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.scene_automations 的数据：-1 rows
/*!40000 ALTER TABLE "scene_automations" DISABLE KEYS */;
/*!40000 ALTER TABLE "scene_automations" ENABLE KEYS */;

-- 导出  表 public.scene_automation_log 结构
CREATE TABLE IF NOT EXISTS "scene_automation_log" (
	"scene_automation_id" VARCHAR(36) NOT NULL,
	"executed_at" TIMESTAMPTZ NOT NULL,
	"detail" TEXT NOT NULL,
	"execution_result" VARCHAR(10) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	CONSTRAINT "scene_automation_log_scene_automation_id_fkey" FOREIGN KEY ("scene_automation_id") REFERENCES "scene_automations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.scene_automation_log 的数据：-1 rows
/*!40000 ALTER TABLE "scene_automation_log" DISABLE KEYS */;
/*!40000 ALTER TABLE "scene_automation_log" ENABLE KEYS */;

-- 导出  表 public.scene_info 结构
CREATE TABLE IF NOT EXISTS "scene_info" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"creator" VARCHAR(36) NOT NULL,
	"updator" VARCHAR(36) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.scene_info 的数据：-1 rows
/*!40000 ALTER TABLE "scene_info" DISABLE KEYS */;
/*!40000 ALTER TABLE "scene_info" ENABLE KEYS */;

-- 导出  表 public.scene_log 结构
CREATE TABLE IF NOT EXISTS "scene_log" (
	"scene_id" VARCHAR(36) NOT NULL,
	"executed_at" TIMESTAMPTZ NOT NULL,
	"detail" TEXT NOT NULL,
	"execution_result" VARCHAR(10) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "scene_log_scene_id_fkey" FOREIGN KEY ("scene_id") REFERENCES "scene_info" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.scene_log 的数据：-1 rows
/*!40000 ALTER TABLE "scene_log" DISABLE KEYS */;
/*!40000 ALTER TABLE "scene_log" ENABLE KEYS */;

-- 导出  表 public.service_access 结构
CREATE TABLE IF NOT EXISTS "service_access" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(100) NOT NULL,
	"service_plugin_id" VARCHAR(36) NOT NULL,
	"voucher" VARCHAR(999) NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"service_access_config" JSON NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"create_at" TIMESTAMPTZ NOT NULL,
	"update_at" TIMESTAMPTZ NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "fk_service_plugin" FOREIGN KEY ("service_plugin_id") REFERENCES "service_plugins" ("id") ON UPDATE NO ACTION ON DELETE RESTRICT
);

-- 正在导出表  public.service_access 的数据：-1 rows
/*!40000 ALTER TABLE "service_access" DISABLE KEYS */;
/*!40000 ALTER TABLE "service_access" ENABLE KEYS */;

-- 导出  表 public.service_plugins 结构
CREATE TABLE IF NOT EXISTS "service_plugins" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"service_identifier" VARCHAR(100) NOT NULL,
	"service_type" INTEGER NOT NULL,
	"last_active_time" TIMESTAMPTZ NULL DEFAULT NULL,
	"version" VARCHAR(100) NULL DEFAULT NULL,
	"create_at" TIMESTAMPTZ NOT NULL,
	"update_at" TIMESTAMPTZ NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"service_config" JSON NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "unique_service_identifier" ("service_identifier"),
	UNIQUE INDEX "unique_name" ("name"),
	CONSTRAINT "service_plugins_service_type_check" CHECK (((service_type = ANY (ARRAY[1, 2]))))
);

-- 正在导出表  public.service_plugins 的数据：-1 rows
/*!40000 ALTER TABLE "service_plugins" DISABLE KEYS */;
REPLACE INTO "service_plugins" ("id", "name", "service_identifier", "service_type", "last_active_time", "version", "create_at", "update_at", "description", "service_config", "remark") VALUES
	('d073ba1d-445a-a07f-430b-cf6d154bc5e8', 'MODBUS-RTU', 'MODBUS_RTU', 1, '2024-12-25 02:30:43.678+00', 'v1.0.1', '2024-12-25 01:05:12.019+00', '2024-12-25 01:05:36.696+00', '', '{"http_address":"172.20.0.10:503","device_type":2,"sub_topic_prefix":"plugin/modbus/","access_address":":502"}', ''),
	('4bec425e-c7e3-476a-0303-ee8193ddf4ca', 'MODBUS-TCP	', 'MODBUS_TCP', 1, '2024-12-25 02:30:43.678+00', 'v1.0.1', '2024-12-25 01:03:25.401+00', '2024-12-25 01:10:02.208+00', '', '{"http_address":"172.20.0.10:503","device_type":2,"sub_topic_prefix":"plugin/modbus/","access_address":":502"}', '');
/*!40000 ALTER TABLE "service_plugins" ENABLE KEYS */;

-- 导出  表 public.sys_dict 结构
CREATE TABLE IF NOT EXISTS "sys_dict" (
	"id" VARCHAR(36) NOT NULL,
	"dict_code" VARCHAR(36) NOT NULL,
	"dict_value" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "sys_dict_dict_code_dict_value_key" ("dict_code", "dict_value")
);

-- 正在导出表  public.sys_dict 的数据：-1 rows
/*!40000 ALTER TABLE "sys_dict" DISABLE KEYS */;
REPLACE INTO "sys_dict" ("id", "dict_code", "dict_value", "created_at", "remark") VALUES
	('0013fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'GATEWAY_PROTOCOL', 'MQTT', '2024-01-18 07:39:38.469+00', NULL),
	('7162fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'DRIECT_ATTACHED_PROTOCOL', 'MQTT', '2024-01-18 07:39:38.469+00', NULL);
/*!40000 ALTER TABLE "sys_dict" ENABLE KEYS */;

-- 导出  表 public.sys_dict_language 结构
CREATE TABLE IF NOT EXISTS "sys_dict_language" (
	"id" VARCHAR(36) NOT NULL,
	"dict_id" VARCHAR(36) NOT NULL,
	"language_code" VARCHAR(36) NOT NULL,
	"translation" VARCHAR(255) NOT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "sys_dict_language_dict_id_language_code_key" ("dict_id", "language_code"),
	CONSTRAINT "sys_dict_language_dict_id_fkey" FOREIGN KEY ("dict_id") REFERENCES "sys_dict" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.sys_dict_language 的数据：-1 rows
/*!40000 ALTER TABLE "sys_dict_language" DISABLE KEYS */;
REPLACE INTO "sys_dict_language" ("id", "dict_id", "language_code", "translation") VALUES
	('001c3960-3067-536d-5c97-7645351a687c', '7162fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'zh_CN', 'MQTT协议'),
	('002c3960-3067-536d-5c97-7645351a687b', '0013fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'zh_CN', 'MQTT协议(网关)'),
	('7162fb9e-e3be-95d4-9c96-f18d1f9ddfss', '7162fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'en_US', 'MQTT Protocol'),
	('7162fb9e-e3be-95d4-9c96-f18d1f9ddfff', '0013fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'en_US', 'MQTT Protocol(Gateway)');
/*!40000 ALTER TABLE "sys_dict_language" ENABLE KEYS */;

-- 导出  表 public.sys_function 结构
CREATE TABLE IF NOT EXISTS "sys_function" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(50) NOT NULL,
	"enable_flag" VARCHAR(20) NOT NULL,
	"description" VARCHAR(500) NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.sys_function 的数据：-1 rows
/*!40000 ALTER TABLE "sys_function" DISABLE KEYS */;
REPLACE INTO "sys_function" ("id", "name", "enable_flag", "description", "remark") VALUES
	('function_1', 'use_captcha', 'disable', '验证码登陆', NULL),
	('function_2', 'enable_reg', 'disable', '租户注册', NULL),
	('function_3', 'frontend_res', 'disable', '前端RSA加密', NULL),
	('function_4', 'shared_account', 'enable', '共享账号', NULL);
/*!40000 ALTER TABLE "sys_function" ENABLE KEYS */;

-- 导出  表 public.sys_ui_elements 结构
CREATE TABLE IF NOT EXISTS "sys_ui_elements" (
	"id" VARCHAR(36) NOT NULL,
	"parent_id" VARCHAR(36) NOT NULL,
	"element_code" VARCHAR(100) NOT NULL,
	"element_type" SMALLINT NOT NULL,
	"orders" SMALLINT NULL DEFAULT NULL,
	"param1" VARCHAR(255) NULL DEFAULT NULL,
	"param2" VARCHAR(255) NULL DEFAULT NULL,
	"param3" VARCHAR(255) NULL DEFAULT NULL,
	"authority" JSON NOT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"multilingual" VARCHAR(100) NULL DEFAULT NULL,
	"route_path" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.sys_ui_elements 的数据：-1 rows
/*!40000 ALTER TABLE "sys_ui_elements" DISABLE KEYS */;
REPLACE INTO "sys_ui_elements" ("id", "parent_id", "element_code", "element_type", "orders", "param1", "param2", "param3", "authority", "description", "created_at", "remark", "multilingual", "route_path") VALUES
	('6e5e0963-46bf-bc27-d792-156e87a69f51', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm', 1, 115, '/alarm', 'simple-icons:antdesign', 'self', '["TENANT_ADMIN","SYS_ADMIN"]', '告警', '2024-03-07 13:46:40.055+00', '', 'route.alarm', 'layout.base'),
	('f9bd5f79-291e-26d2-1553-473c04b15ce4', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_setting', 3, 42, '/management/setting', 'uil:brightness-plus', 'self', '["SYS_ADMIN"]', '系统设置', '2024-02-18 09:52:08.236+00', '', 'route.management_setting', 'view.management_setting'),
	('51381989-1160-93cd-182e-d44a1c4ab89b', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_scene-manage', 3, 1142, '/automation/scene-manage', 'uil:brightness-plus', 'self', '["TENANT_ADMIN","SYS_ADMIN"]', '场景管理', '2024-03-07 13:44:11.106+00', '', 'route.automation_scene-manage', 'view.automation_scene-manage'),
	('8f4e9058-e30d-2fb5-ac6d-784613234883', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm-information', 3, 1151, '/alarm/alarm-information', 'ph:alarm', 'basic', '["TENANT_ADMIN"]', '告警信息', '2024-03-07 13:47:22.817+00', '', 'default', ''),
	('b6d57a4a-d37a-9d9d-6e4e-be33b955ff04', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm_notification-group', 3, 1152, '/alarm/notification-group', 'simple-icons:apacheecharts', 'basic', '["TENANT_ADMIN"]', '通知组', '2024-03-07 13:48:15.416+00', '', 'route.alarm_notification-group', 'view.alarm_notification-group'),
	('faf7e607-00ae-3483-40a1-b74f9245b100', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_auth', 3, 43, '/management/auth', 'ic:baseline-security', 'self', '["SYS_ADMIN"]', '菜单管理', '2024-02-18 09:49:31.209+00', '', 'route.management_auth', 'view.management_auth'),
	('e1ebd134-53df-3105-35f4-489fc674d173', '0', 'management', 1, 120, '/management', 'carbon:cloud-service-management', 'self', '["SYS_ADMIN","TENANT_ADMIN"]', '系统管理', '2024-02-18 09:48:45.265+00', '', 'route.management', 'layout.base'),
	('96aa2fac-90b2-aca1-1ce0-51b5060f4081', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_linkage-edit', 3, 1143, '/automation/linkage-edit', '', '1', '["TENANT_ADMIN","SYS_ADMIN"]', '场景联动编辑', '2024-03-14 17:36:03.938+00', '', 'route.automation_linkage-edit', 'view.automation_linkage-edit'),
	('e619f321-9823-b563-b24d-ecc16d7b23cc', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm_notification-record', 3, 1153, '/alarm/notification-record', 'mdi:monitor-dashboard', 'basic', '["TENANT_ADMIN"]', '通知记录', '2024-03-07 13:48:56.415+00', '', 'route.alarm_notification-record', 'view.alarm_notification-record'),
	('01dab674-9556-cdd7-b800-78bcb366adb4', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_scene-linkage', 3, 1141, '/automation/scene-linkage', 'mdi:airplane-edit', 'self', '["TENANT_ADMIN","SYS_ADMIN"]', '场景联动', '2024-03-07 13:43:33.92+00', '', 'route.automation_scene-linkage', 'view.automation_scene-linkage'),
	('c078182f-bf4b-b560-da97-02926fa98f78', '650bc444-7672-1123-1e41-7e37365b0186', 'alarm_notification-record', 3, 1, '/alarm/notification-record', 'icon-park-outline:editor', 'self', '["TENANT_ADMIN"]', '通知记录', '2024-03-20 02:04:34.927+00', '', 'route.alarm_notification-record', 'view.alarm_notification-record'),
	('485c2a20-ebc5-2216-4871-26453470d290', '650bc444-7672-1123-1e41-7e37365b0186', 'alarm_warning-message', 3, 999, '/alarm/warning-message', 'mdi:airballoon', 'self', '["TENANT_ADMIN"]', '警告信息', '2024-03-17 07:27:40.378+00', '', 'route.alarm_warning-message', 'view.alarm_warning-message'),
	('2f3ffd60-efec-aafb-a866-f1cb79f88390', 'e1ebd134-53df-3105-35f4-489fc674d173', 'system-management-user_system-log', 3, 1171, '/system-management-user/system-log', 'mdi:monitor-dashboard', 'basic', '["TENANT_ADMIN"]', '系统日志', '2024-03-07 14:23:08.576+00', '', 'route.system-management-user_system-log', 'view.system-management-user_system-log'),
	('e186a671-8e24-143a-5a2c-27a1f5f38bf3', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_config-edit', 3, 1128, '/device/config-edit', '', '1', '["TENANT_ADMIN"]', '设备配置编辑', '2024-03-11 13:49:34.952+00', '', 'route.device_config-edit', 'view.device_config-edit'),
	('a2c53126-029f-7138-4d7a-f45491f396da', '0', 'apply', 1, 3, '/apply', 'mdi:apps-box', '0', '["SYS_ADMIN"]', '应用管理', '2024-02-18 09:59:31.642+00', '', 'route.apply', 'layout.base'),
	('49857e46-2176-610e-98fc-892b4fde50f9', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_details', 3, 1124, '/device/details', 'mdi:monitor-dashboard', '1', '["TENANT_ADMIN"]', '设备详情', '2024-03-05 09:52:21.434+00', '', 'route.device_details', 'view.device_details'),
	('ed4a5cfa-03e7-ccc0-6cc8-bcadccd25541', '95e2a961-382b-f4a6-87b3-1898123c95bc', 'visualization_kanban-details', 3, 1132, '/visualization/kanban-details', 'ic:baseline-credit-card', '1', '["TENANT_ADMIN","SYS_ADMIN"]', '看板详情', '2024-03-12 02:14:50.152+00', '', '看板详情', 'view.visualization_panel-details'),
	('502a0d6c-750e-92f6-a1a7-ffdd362dbbac', '95e2a961-382b-f4a6-87b3-1898123c95bc', 'visualization_panel-preview', 3, 1133, '/visualization/panel-preview', '', '1', '["TENANT_ADMIN","SYS_ADMIN"]', '看板预览', '2024-03-12 02:16:29.336+00', '', 'route.visualization_panel-preview', 'view.visualization_panel-preview'),
	('75785418-a5af-d790-0783-e4ee4e42521e', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_grouping', 3, 1122, '/device/grouping', 'material-symbols:grid-on-outline-sharp', '0', '["TENANT_ADMIN"]', '设备分组', '2024-03-05 09:53:25.004+00', '', 'route.device_grouping', 'view.device_grouping'),
	('8de46003-170c-a24d-6baf-84d1c7298aa3', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_grouping-details', 3, 1123, '/device/grouping-details', '', '1', '["TENANT_ADMIN"]', '分组详情', '2024-03-05 09:54:23.158+00', '', 'route.device_grouping-details', 'view.device_grouping-details'),
	('5373a6a2-1861-af35-eb4c-adfd5ca55ecd', '0', 'device', 1, 112, '/device', 'icon-park-outline:workbench', '0', '["TENANT_ADMIN"]', '设备接入', '2024-03-05 09:51:19.298+00', '', 'route.device', 'layout.base'),
	('7419e37e-c167-f12b-7ace-76e479144181', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_template', 3, 1127, '/device/template', 'simple-icons:apacheecharts', 'self', '["TENANT_ADMIN"]', '功能模板', '2024-03-05 10:01:29.826+00', '定义物模型和显示图表', 'route.device_template', 'view.device_template'),
	('774a716d-9861-bac9-857f-acaa25e7659f', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_config', 3, 1126, '/device/config', 'clarity:plugin-line', 'self', '["TENANT_ADMIN"]', '配置模板', '2024-03-05 14:06:53.842+00', '设备的协议和其他参数等所有配置', 'route.device_config', 'view.device_config'),
	('36c4f5ce-3279-55f2-ede2-81b4a0bae24b', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_user', 3, 41, '/management/user', 'ic:round-manage-accounts', 'self', '["SYS_ADMIN"]', '租户管理 ', '2024-02-18 09:50:48.999+00', '', 'route.management_user', 'view.management_user'),
	('c4dff952-3bf4-8102-6882-e9d3f3cffbda', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_manage', 3, 1121, '/device/manage', 'icon-park-outline:analysis', '0', '["TENANT_ADMIN"]', '设备管理', '2024-03-05 09:55:08.17+00', '', 'route.device_manage', 'view.device_manage'),
	('fec91838-d30d-7d66-6715-0912f1b171d8', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_notification', 3, 44, '/management/notification', 'mdi:alert', 'self', '["SYS_ADMIN"]', '通知配置', '2024-03-15 11:50:07.495+00', '', 'route.management_notification', 'view.management_notification'),
	('82c46beb-9ec4-8a3d-c6e4-04ba426e525a', '650bc444-7672-1123-1e41-7e37365b0186', 'alarm_notification-group', 3, 1, '/alarm/notification-group', 'ic:round-supervisor-account', 'basic', '["TENANT_ADMIN"]', '通知组', '2024-03-20 02:03:19.955+00', '', 'route.alarm_notification-group', 'view.alarm_notification-group'),
	('64f684f1-390c-b5f2-9994-36895025df8a', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_space-management', 3, 10, 'automation/space-management', 'ic:baseline-security', '1', '["TENANT_ADMIN","SYS_ADMIN"]', '场景管理', '2024-03-22 05:25:38.82+00', '', 'default', 'view.automation space-management'),
	('676e8f33-875a-0473-e9ca-c82fd09fef57', '0', 'automation', 1, 114, '/automation', 'material-symbols:device-hub', 'self', '["SYS_ADMIN"]', '自动化', '2024-03-07 13:41:17.921+00', '', 'route.automation', 'layout.base'),
	('650bc444-7672-1123-1e41-7e37365b0186', '0', 'alarm', 1, 115, '/alarm', 'mdi:alert', 'self', '["SYS_ADMIN"]', '告警', '2024-03-17 01:01:52.183+00', '', 'route.alarm', 'layout.base'),
	('76bfc16e-ed22-bcc0-c688-d462666e8a8d', '0', 'personal-center', 3, 999, '/personal-center', 'carbon:user-role', '1', '["TENANT_ADMIN","SYS_ADMIN"]', '个人中心', '2024-03-17 01:27:01.048+00', '', 'route.personal_center', 'layout.base$view.personal-center'),
	('975c9550-5db9-7b4c-5dea-7a4c326a37ff', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_scene-edit', 3, 1, '/automation/scene-edit', 'mdi:apps-box', '1', '["TENANT_ADMIN"]', '新增场景', '2024-04-04 02:50:43.219+00', '', 'route.automation_scene-edit', 'view.automation_scene-edit'),
	('680cae76-6c50-90e6-c2f9-58d01389aa08', '9a11b3e4-9982-a0f0-996c-a9be6e738947', 'data-service_rule-engine', 3, 21, '/data-service/rule-engine', 'mdi:menu', '1', '["SYS_ADMIN"]', '规则引擎', '2024-03-07 09:06:02.804+00', '', 'route.data-service_rule-engine', 'view.data-service_rule-engine'),
	('a2654c98-3749-c88b-0472-b414049ca532', '95e2a961-382b-f4a6-87b3-1898123c95bc', 'route.visualization_kanban', 3, 1131, '/visualization/kanban', 'tabler:device-tv', 'self', '["TENANT_ADMIN","SYS_ADMIN"]', '看板', '2024-03-07 13:39:58.608+00', '', 'route.visualization_kanban', 'view.visualization_panel'),
	('86cb08fa-8b08-3d99-4b3a-d6132ee93a0f', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_config-detail', 3, 1127, '/device/config-detail', 'icon-park-outline:data-server', '1', '["TENANT_ADMIN"]', '设备配置详情', '2024-03-10 03:13:25.253+00', '', 'route.device_config-detail', 'view.device_config-detail'),
	('59612e2f-e297-acb7-fcf4-143bf6e66109', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_details-child', 3, 1124, '/device/details-child', '', '1', '["TENANT_ADMIN"]', '子设备详情', '2024-05-10 12:33:34.869+00', '', 'route.device_details-child', 'view.device_details-child'),
	('29a684f9-c2bb-1a6f-6045-314944bef580', 'a2c53126-029f-7138-4d7a-f45491f396da', 'plug_in', 3, 32, '/apply/plugin', 'mdi:emoticon', '0', '["SYS_ADMIN"]', '插件管理', '2024-06-28 17:04:51.301+00', '', 'route.apply_in', ''),
	('a190f7a5-1501-3814-9dd1-f3e1fbe7265e', '0', 'home', 3, 0, '/home', 'mdi:alpha-f-box-outline', 'self', '["SYS_ADMIN","TENANT_ADMIN"]', '首页', '2024-02-26 08:07:20.202+00', 'home', 'route.home', 'layout.base$view.home'),
	('9a11b3e4-9982-a0f0-996c-a9be6e738947', '0', 'data-service', 1, 2, '/data-service', 'mdi:monitor-dashboard', '1', '["SYS_ADMIN"]', '数据服务', '2024-03-07 09:05:04.101+00', '', 'route.data-service', 'layout.base'),
	('3aaca04b-2a2e-dfca-9fb4-0b2819362783', '2cc0c5ba-f086-91e5-0b8c-ad0546b1f2a9', 'test_kan-ban-test', 3, 1, '/test/kan-ban-test', '', '1', '["SYS_ADMIN","TENANT_ADMIN"]', '看板测试', '2024-05-20 17:17:16.911+00', '', 'route.test_kan-ban-test', 'view.test_kan-ban-test'),
	('2fe87d7c-627e-9ca3-94dd-6d0249853bd4', '990af72f-06ce-5f23-3af6-1694bd479c96', 'management_user', 3, 1, '/management/user', '', '0', '["SYS_ADMIN","TENANT_ADMIN"]', 'ment-user', '2024-09-04 02:04:34.658+00', '', 'default', ''),
	('18892c6e-ca04-f2b5-c243-f2c7230b3f33', '990af72f-06ce-5f23-3af6-1694bd479c96', 'manage_user', 3, 1, '/manage/user', '', '0', '["TENANT_ADMIN","SYS_ADMIN"]', 'user', '2024-09-04 02:05:06.377+00', '', 'default', ''),
	('95e2a961-382b-f4a6-87b3-1898123c95bc', '0', 'visualization', 1, 113, '/visualization', 'icon-park-outline:data-server', 'self', '["SYS_ADMIN"]', '可视化', '2024-03-07 13:37:16.042+00', '', 'route.visualization', 'layout.base'),
	('f960c45c-6d5b-e67a-c4ff-1f0e869c1625', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_service-details', 3, 1130, '/device/service-details', 'ph:align-bottom', '1', '["SYS_ADMIN"]', '服务详情', '2024-07-01 15:16:56.668+00', '', 'route.device_service_details', ''),
	('075d9f19-5618-bb9b-6ccd-f382bfd3292b', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_service-access', 3, 1129, '/device/service-access', 'mdi:ab-testing', '0', '["SYS_ADMIN"]', '服务接入点管理', '2024-07-01 13:52:09.402+00', '', 'route.device_service_access', '');
/*!40000 ALTER TABLE "sys_ui_elements" ENABLE KEYS */;

-- 导出  表 public.sys_version 结构
CREATE TABLE IF NOT EXISTS "sys_version" (
	"version_number" INTEGER NOT NULL DEFAULT '0',
	"version" VARCHAR(255) NOT NULL,
	PRIMARY KEY ("version_number")
);

-- 正在导出表  public.sys_version 的数据：-1 rows
/*!40000 ALTER TABLE "sys_version" DISABLE KEYS */;
REPLACE INTO "sys_version" ("version_number", "version") VALUES
	(4, '0.0.4');
/*!40000 ALTER TABLE "sys_version" ENABLE KEYS */;

-- 导出  表 public.telemetry_current_datas 结构
CREATE TABLE IF NOT EXISTS "telemetry_current_datas" (
	"device_id" VARCHAR(36) NOT NULL,
	"key" VARCHAR(255) NOT NULL,
	"ts" TIMESTAMPTZ NOT NULL,
	"bool_v" BOOLEAN NULL DEFAULT NULL,
	"number_v" DOUBLE PRECISION NULL DEFAULT NULL,
	"string_v" TEXT NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	UNIQUE INDEX "telemetry_current_datas_unique" ("device_id", "key"),
	INDEX "telemetry_datas_ts_idx_copy1" ("ts")
);

-- 正在导出表  public.telemetry_current_datas 的数据：-1 rows
/*!40000 ALTER TABLE "telemetry_current_datas" DISABLE KEYS */;
REPLACE INTO "telemetry_current_datas" ("device_id", "key", "ts", "bool_v", "number_v", "string_v", "tenant_id") VALUES
	('112233445566', 'Mode', '2025-02-13 14:35:11.189+00', NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'set_midpoint', '2025-02-13 14:35:11.189+00', NULL, 66, NULL, 'd616bcbb'),
	('112233445566', 'wind_speed', '2025-02-13 14:35:11.189+00', NULL, 120, NULL, 'd616bcbb'),
	('112233445566', 'Long', '2025-02-13 14:35:11.189+00', NULL, 55, NULL, 'd616bcbb'),
	('112233445566', 'Device_h', '2025-02-13 14:35:11.189+00', NULL, 12, NULL, 'd616bcbb'),
	('112233445566', 'Device_a', '2025-02-13 14:35:11.189+00', NULL, 33, NULL, 'd616bcbb');
/*!40000 ALTER TABLE "telemetry_current_datas" ENABLE KEYS */;

-- 导出  表 public.telemetry_datas 结构
CREATE TABLE IF NOT EXISTS "telemetry_datas" (
	"device_id" VARCHAR(36) NOT NULL,
	"key" VARCHAR(255) NOT NULL,
	"ts" BIGINT NOT NULL,
	"bool_v" BOOLEAN NULL DEFAULT NULL,
	"number_v" DOUBLE PRECISION NULL DEFAULT NULL,
	"string_v" TEXT NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	UNIQUE INDEX "telemetry_datas_device_id_key_ts_key" ("device_id", "key", "ts"),
	INDEX "telemetry_datas_ts_idx" ("ts")
);

-- 正在导出表  public.telemetry_datas 的数据：-1 rows
/*!40000 ALTER TABLE "telemetry_datas" DISABLE KEYS */;
REPLACE INTO "telemetry_datas" ("device_id", "key", "ts", "bool_v", "number_v", "string_v", "tenant_id") VALUES
	('112233445566', 'Mode', 1739455956451, NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'set_midpoint', 1739456955087, NULL, 88, NULL, 'd616bcbb'),
	('112233445566', 'Mode', 1739456955087, NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'Mode', 1739457128148, NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'set_midpoint', 1739457128148, NULL, 56, NULL, 'd616bcbb'),
	('112233445566', 'Mode', 1739457133526, NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'set_midpoint', 1739457133526, NULL, 66, NULL, 'd616bcbb'),
	('112233445566', 'Mode', 1739457280006, NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'set_midpoint', 1739457280006, NULL, 66, NULL, 'd616bcbb'),
	('112233445566', 'wind_speed', 1739457280006, NULL, 120, NULL, 'd616bcbb'),
	('112233445566', 'Long', 1739457280006, NULL, 55, NULL, 'd616bcbb'),
	('112233445566', 'Mode', 1739457311189, NULL, 1, NULL, 'd616bcbb'),
	('112233445566', 'set_midpoint', 1739457311189, NULL, 66, NULL, 'd616bcbb'),
	('112233445566', 'wind_speed', 1739457311189, NULL, 120, NULL, 'd616bcbb'),
	('112233445566', 'Long', 1739457311189, NULL, 55, NULL, 'd616bcbb'),
	('112233445566', 'Device_h', 1739457311189, NULL, 12, NULL, 'd616bcbb'),
	('112233445566', 'Device_a', 1739457311189, NULL, 33, NULL, 'd616bcbb');
/*!40000 ALTER TABLE "telemetry_datas" ENABLE KEYS */;

-- 导出  表 public.telemetry_set_logs 结构
CREATE TABLE IF NOT EXISTS "telemetry_set_logs" (
	"id" VARCHAR(36) NOT NULL,
	"device_id" VARCHAR(36) NOT NULL,
	"operation_type" VARCHAR(255) NULL DEFAULT NULL,
	"data" JSON NULL DEFAULT NULL,
	"status" VARCHAR(2) NULL DEFAULT NULL,
	"error_message" VARCHAR(500) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"user_id" VARCHAR(36) NULL DEFAULT NULL,
	"description" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "telemetry_set_logs_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- 正在导出表  public.telemetry_set_logs 的数据：-1 rows
/*!40000 ALTER TABLE "telemetry_set_logs" DISABLE KEYS */;
/*!40000 ALTER TABLE "telemetry_set_logs" ENABLE KEYS */;

-- 导出  表 public.users 结构
CREATE TABLE IF NOT EXISTS "users" (
	"id" VARCHAR(36) NOT NULL,
	"name" VARCHAR(255) NULL DEFAULT NULL,
	"phone_number" VARCHAR(50) NOT NULL,
	"email" VARCHAR(255) NOT NULL,
	"status" VARCHAR(2) NULL DEFAULT NULL,
	"authority" VARCHAR(50) NULL DEFAULT NULL,
	"password" VARCHAR(255) NOT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"additional_info" JSON NULL DEFAULT '{}',
	"created_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"password_last_updated" TIMESTAMPTZ NULL DEFAULT NULL,
	"last_visit_time" TIMESTAMPTZ NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "users_un" ("email")
);

-- 正在导出表  public.users 的数据：2 rows
/*!40000 ALTER TABLE "users" DISABLE KEYS */;
REPLACE INTO "users" ("id", "name", "phone_number", "email", "status", "authority", "password", "tenant_id", "remark", "additional_info", "created_at", "updated_at", "password_last_updated", "last_visit_time") VALUES
	('11111111-4fe9-b409-67c3-111111111111', 'Tenant', '17366666666', 'test@qq.com', 'N', 'TENANT_ADMIN', '$2a$10$HKe1P81Az0mTBFl8AmrRi.1iFAnG.Z6v5g7sfgtfhENW2KVlPAypy', 'd616bcbb', '', '{}', '2024-06-05 08:48:11.097+00', '2025-02-13 14:00:44.863691+00', '2025-02-13 14:00:40.263267+00', '2025-02-13 14:00:44.863691+00'),
	('00000000-4fe9-b409-67c3-000000000000', 'admin', '1231231321', 'super@super.cn', 'N', 'SYS_ADMIN', '$2a$10$dPDIqoOEt.rSDwEWsSHCqe9/PJEsnWvRK76DwXVZUFM/7J0D3ikfq', 'aaaaaa', 'dolor', '{}', NULL, '2025-02-13 15:04:06.647702+00', NULL, '2025-02-13 15:04:06.647702+00');
/*!40000 ALTER TABLE "users" ENABLE KEYS */;

-- 导出  表 public.vis_dashboard 结构
CREATE TABLE IF NOT EXISTS "vis_dashboard" (
	"id" VARCHAR(36) NOT NULL,
	"relation_id" VARCHAR(36) NULL DEFAULT NULL,
	"json_data" JSON NULL DEFAULT '{}',
	"dashboard_name" VARCHAR(99) NULL DEFAULT NULL,
	"create_at" TIMESTAMP NULL DEFAULT NULL,
	"sort" INTEGER NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	"tenant_id" VARCHAR(36) NULL DEFAULT NULL,
	"share_id" VARCHAR(36) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.vis_dashboard 的数据：0 rows
/*!40000 ALTER TABLE "vis_dashboard" DISABLE KEYS */;
/*!40000 ALTER TABLE "vis_dashboard" ENABLE KEYS */;

-- 导出  表 public.vis_files 结构
CREATE TABLE IF NOT EXISTS "vis_files" (
	"id" VARCHAR(36) NOT NULL,
	"vis_plugin_id" VARCHAR(36) NOT NULL,
	"file_name" VARCHAR(150) NULL DEFAULT NULL,
	"file_url" VARCHAR(150) NULL DEFAULT NULL,
	"file_size" VARCHAR(20) NULL DEFAULT NULL,
	"create_at" BIGINT NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.vis_files 的数据：-1 rows
/*!40000 ALTER TABLE "vis_files" DISABLE KEYS */;
/*!40000 ALTER TABLE "vis_files" ENABLE KEYS */;

-- 导出  表 public.vis_plugin 结构
CREATE TABLE IF NOT EXISTS "vis_plugin" (
	"id" VARCHAR(36) NOT NULL,
	"tenant_id" VARCHAR(36) NOT NULL,
	"plugin_name" VARCHAR(150) NOT NULL,
	"plugin_description" VARCHAR(150) NULL DEFAULT NULL,
	"create_at" BIGINT NULL DEFAULT NULL,
	"remark" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("id")
);

-- 正在导出表  public.vis_plugin 的数据：-1 rows
/*!40000 ALTER TABLE "vis_plugin" DISABLE KEYS */;
/*!40000 ALTER TABLE "vis_plugin" ENABLE KEYS */;

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
