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