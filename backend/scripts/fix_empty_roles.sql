-- Script para corregir usuarios sin rol asignado
-- Ejecutar este script si encuentras usuarios con campos de rol vacíos en la tabla

-- 1. Ver cuántos usuarios tienen rol vacío o NULL
SELECT COUNT(*) as usuarios_sin_rol
FROM usuarios
WHERE rol IS NULL OR rol = '';

-- 2. Ver todos los usuarios sin rol
SELECT id, nombre, apellido_paterno, email, rol, status, created_at
FROM usuarios
WHERE rol IS NULL OR rol = '';

-- 3. Asignar rol 'turista' por defecto a usuarios sin rol
-- (Descomentar la siguiente línea para ejecutar la actualización)
-- UPDATE usuarios SET rol = 'turista' WHERE rol IS NULL OR rol = '';

-- 4. Verificar que todos los usuarios ahora tienen un rol válido
SELECT rol, COUNT(*) as cantidad
FROM usuarios
GROUP BY rol
ORDER BY cantidad DESC;

-- 5. Si necesitas asignar roles específicos basándote en algún criterio:
-- Por ejemplo, asignar 'admin' al primer usuario creado:
-- UPDATE usuarios SET rol = 'admin' WHERE id = 1;

-- 6. Verificación final - todos los usuarios deben tener uno de estos roles:
SELECT
    COUNT(*) as total_usuarios,
    SUM(CASE WHEN rol = 'admin' THEN 1 ELSE 0 END) as admins,
    SUM(CASE WHEN rol = 'turista' THEN 1 ELSE 0 END) as turistas,
    SUM(CASE WHEN rol = 'encargado_agencia' THEN 1 ELSE 0 END) as encargados,
    SUM(CASE WHEN rol NOT IN ('admin', 'turista', 'encargado_agencia') OR rol IS NULL OR rol = '' THEN 1 ELSE 0 END) as sin_rol
FROM usuarios;
