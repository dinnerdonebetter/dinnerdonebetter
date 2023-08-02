-- name: CreateServiceSettingConfiguration :exec

INSERT INTO service_setting_configurations (id,value,notes,service_setting_id,belongs_to_user,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6);
