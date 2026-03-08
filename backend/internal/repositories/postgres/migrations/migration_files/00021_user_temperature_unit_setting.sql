-- Add user_temperature_unit service setting (celsius vs fahrenheit, default fahrenheit)
INSERT INTO service_settings (id, name, type, description, default_value, enumeration, admins_only)
VALUES (
  'd6me6i4n9qd3gcf5j1p0',
  'user_temperature_unit',
  'user',
  'Preferred unit for displaying temperatures (e.g. oven, storage)',
  'fahrenheit',
  'celsius|fahrenheit',
  false
);
