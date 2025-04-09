-- 场景管理动作值最大支持到2048
ALTER TABLE public.scene_action_info 
ALTER COLUMN action_value TYPE VARCHAR(2048);