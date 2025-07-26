-- 场景管理动作值最大支持到2048
ALTER TABLE public.scene_action_info 
ALTER COLUMN action_value TYPE VARCHAR(2048);

-- 增加太阳总辐射的单位
UPDATE device_model_telemetry
SET 
  unit = 'W/m²',
  updated_at = CURRENT_TIMESTAMP
WHERE 
  id = '241de086-8339-5488-91de-d27bf49ec43d';


INSERT INTO "device_model_attributes" ("id", "device_template_id", "data_name", "data_identifier", "read_write_flag", "data_type", "unit", "description", "additional_info", "created_at", "updated_at", "remark", "tenant_id") VALUES
	('c262bba2-702a-a3c7-5b3b-663096342462', 'd391b336-4273-d101-27f8-d1fbedfde866', '经度', 'longitude', 'R', 'String', '', '', '[]', '2025-04-19 23:58:51.738758+08', '2025-04-19 23:58:51.738758+08', NULL, 'd616bcbb'),
	('289653c1-462c-7f28-695f-a3f9440670fa', 'd391b336-4273-d101-27f8-d1fbedfde866', '纬度', 'latitude', 'R', 'String', '', '', '[]', '2025-04-19 23:59:14.780432+08', '2025-04-19 23:59:14.780432+08', NULL, 'd616bcbb'),
	('4dcbe6e2-ff16-1067-172c-b0564cd84525', 'd391b336-4273-d101-27f8-d1fbedfde866', 'GPS可用卫星数', 'starNum', 'R', 'Number', '', '', '[]', '2025-04-20 00:00:03.757069+08', '2025-04-20 00:00:03.757069+08', NULL, 'd616bcbb'),
	('ac091bba-ed2e-f85e-7df1-8ce9aea4dca0', 'd391b336-4273-d101-27f8-d1fbedfde866', '移动网络信号', 'signal', 'R', 'String', '', '', '[]', '2025-04-20 01:35:45.328602+08', '2025-04-20 01:35:45.328602+08', NULL, 'd616bcbb'),
	('7e95b6cf-a254-b3ef-99d7-6b8fbe321d4e', 'd391b336-4273-d101-27f8-d1fbedfde866', '移动网络类型', 'network', 'RW', 'String', '', '', '[]', '2025-04-19 23:58:18.08621+08', '2025-04-20 01:35:55.214328+08', NULL, 'd616bcbb'),
	('ce4b57e1-f557-fad7-213d-b4c64949c03d', 'd391b336-4273-d101-27f8-d1fbedfde866', '4G模组固件', 'modelVersion', 'R', 'String', '', '', '[]', '2025-04-20 01:36:26.643714+08', '2025-04-20 01:36:26.643714+08', NULL, 'd616bcbb');