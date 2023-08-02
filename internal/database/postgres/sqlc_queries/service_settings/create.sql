-- name: CreateServiceSetting :exec

INSERT INTO service_settings (id,name,type,description,default_value,admins_only,enumeration) VALUES
($1,$2,$3,$4,$5,$6,$7);
