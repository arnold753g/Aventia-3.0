package database

import (
	"fmt"

	"gorm.io/gorm"
)

func ApplySQLBootstrap(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	if err := ensureProcesarCompraPaquete(db); err != nil {
		return err
	}

	if err := ensureAgenciaMetricsFunction(db); err != nil {
		return err
	}

	if err := ensureTriggerPagoConfirmado(db); err != nil {
		return err
	}

	if err := ensureTriggerPagoRechazado(db); err != nil {
		return err
	}

	if err := ensureNotificacionesTable(db); err != nil {
		return err
	}

	if err := ensureTriggerNuevoPago(db); err != nil {
		return err
	}

	return nil
}

func ensureProcesarCompraPaquete(db *gorm.DB) error {
	// Always replace to keep the function compatible with restored schemas.
	if err := db.Exec(sqlProcesarCompraPaquete).Error; err != nil {
		return fmt.Errorf("procesar_compra_paquete bootstrap failed: %w", err)
	}

	return nil
}

const sqlProcesarCompraPaquete = `
DROP FUNCTION IF EXISTS public.procesar_compra_paquete(
    INTEGER,
    INTEGER,
    DATE,
    TEXT,
    BOOLEAN,
    INTEGER,
    INTEGER,
    INTEGER,
    BOOLEAN,
    TEXT,
    TEXT
);

CREATE OR REPLACE FUNCTION public.procesar_compra_paquete(
    p_turista_id INTEGER,
    p_paquete_id INTEGER,
    p_fecha_seleccionada DATE,
    p_tipo_compra TEXT,
    p_extranjero BOOLEAN,
    p_cantidad_adultos INTEGER,
    p_cantidad_ninos_pagan INTEGER,
    p_cantidad_ninos_gratis INTEGER,
    p_tiene_discapacidad BOOLEAN,
    p_descripcion_discapacidad TEXT,
    p_notas_turista TEXT
)
RETURNS TABLE (
    compra_id INTEGER,
    salida_id INTEGER,
    precio_total NUMERIC,
    mensaje TEXT,
    success BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_paquete RECORD;
    v_max_salidas_por_dia INTEGER := 5;
    v_max_salidas_por_horario INTEGER := 3;
    v_recargo_privado_porcentaje NUMERIC := 0;

    v_total_participantes INTEGER := 0;
    v_personas_pagan INTEGER := 0;

    v_horario_capacidad TEXT := 'todo_dia';
    v_horario_seleccionado TEXT := NULL;

    v_salidas_dia INTEGER := 0;
    v_salidas_horario INTEGER := 0;

    v_salida_id INTEGER := 0;
    v_compra_id INTEGER := 0;

    v_precio_base NUMERIC := 0;
    v_recargo_extranjero_total NUMERIC := 0;
    v_subtotal_base NUMERIC := 0;
    v_subtotal_full NUMERIC := 0;
    v_recargo_privado NUMERIC := 0;
    v_total_recargo NUMERIC := 0;
    v_precio_total NUMERIC := 0;
BEGIN
    -- Validar turista
    IF NOT EXISTS (SELECT 1 FROM usuarios WHERE id = p_turista_id) THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Turista no encontrado';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM usuarios WHERE id = p_turista_id AND rol = 'turista') THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Solo usuarios con rol "turista" pueden realizar compras';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    -- Validaciones básicas de participantes
    IF p_cantidad_adultos IS NULL OR p_cantidad_adultos < 1 THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Debe haber al menos 1 adulto';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    IF COALESCE(p_cantidad_ninos_pagan, 0) < 0 OR COALESCE(p_cantidad_ninos_gratis, 0) < 0 THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Las cantidades de niños no pueden ser negativas';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    v_total_participantes := p_cantidad_adultos
        + COALESCE(p_cantidad_ninos_pagan, 0)
        + COALESCE(p_cantidad_ninos_gratis, 0);
    v_personas_pagan := p_cantidad_adultos + COALESCE(p_cantidad_ninos_pagan, 0);

    IF v_total_participantes < 1 THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Debe registrar al menos 1 participante';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    -- Cargar paquete y validar estado/visibilidad de agencia y paquete
    SELECT
        p.id,
        p.agencia_id,
        p.frecuencia,
        p.fecha_salida_fija,
        p.duracion_dias,
        p.horario,
        p.permite_privado,
        p.dias_previos_compra,
        p.cupo_minimo,
        p.cupo_maximo,
        p.precio_base_nacionales,
        p.precio_adicional_extranjeros,
        p.status,
        p.visible_publico,
        a.status AS agencia_status,
        a.visible_publico AS agencia_visible
    INTO v_paquete
    FROM paquetes_turisticos p
    JOIN agencias_turismo a ON a.id = p.agencia_id
    WHERE p.id = p_paquete_id;

    IF NOT FOUND THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Paquete no encontrado';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    IF v_paquete.status <> 'activo' OR v_paquete.visible_publico <> TRUE THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'El paquete no está disponible';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    IF v_paquete.agencia_status <> 'activa' OR v_paquete.agencia_visible <> TRUE THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'La agencia no está disponible';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    -- Validar tipo de compra
    IF p_tipo_compra NOT IN ('compartido', 'privado') THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'Tipo de compra inválido';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    IF p_tipo_compra = 'privado' THEN
        IF v_paquete.frecuencia <> 'salida_diaria' THEN
            compra_id := 0;
            salida_id := 0;
            precio_total := 0;
            mensaje := 'El tipo privado solo está disponible para paquetes de salida diaria';
            success := FALSE;
            RETURN NEXT;
            RETURN;
        END IF;
        IF v_paquete.permite_privado IS DISTINCT FROM TRUE THEN
            compra_id := 0;
            salida_id := 0;
            precio_total := 0;
            mensaje := 'Este paquete no permite compras privadas';
            success := FALSE;
            RETURN NEXT;
            RETURN;
        END IF;
    END IF;

    -- Validar fecha seleccionada
    IF p_fecha_seleccionada < (CURRENT_DATE + COALESCE(v_paquete.dias_previos_compra, 1)::int) THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'La fecha seleccionada no cumple los días previos de compra';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    IF v_paquete.frecuencia = 'salida_unica' THEN
        IF v_paquete.fecha_salida_fija IS NULL THEN
            compra_id := 0;
            salida_id := 0;
            precio_total := 0;
            mensaje := 'El paquete no tiene fecha de salida fija configurada';
            success := FALSE;
            RETURN NEXT;
            RETURN;
        END IF;
        IF p_fecha_seleccionada <> v_paquete.fecha_salida_fija::date THEN
            compra_id := 0;
            salida_id := 0;
            precio_total := 0;
            mensaje := 'La compra debe realizarse en la fecha de salida fija del paquete';
            success := FALSE;
            RETURN NEXT;
            RETURN;
        END IF;
    END IF;

    -- Horario (para capacidad y para la compra)
    IF COALESCE(v_paquete.duracion_dias, 1) > 1 THEN
        v_horario_seleccionado := NULL;
        v_horario_capacidad := 'todo_dia';
    ELSE
        v_horario_seleccionado := v_paquete.horario;
        v_horario_capacidad := COALESCE(v_paquete.horario, 'todo_dia');
    END IF;

    -- Políticas (recargo privado)
    SELECT recargo_privado_porcentaje
    INTO v_recargo_privado_porcentaje
    FROM paquete_politicas
    WHERE agencia_id = v_paquete.agencia_id;

    IF NOT FOUND THEN
        v_recargo_privado_porcentaje := 0;
    END IF;

    -- Capacidad (si no existe, defaults)
    SELECT max_salidas_por_dia, max_salidas_por_horario
    INTO v_max_salidas_por_dia, v_max_salidas_por_horario
    FROM agencia_capacidad
    WHERE agencia_id = v_paquete.agencia_id;

    IF NOT FOUND THEN
        v_max_salidas_por_dia := 5;
        v_max_salidas_por_horario := 3;
    END IF;

    -- Validar cupo máximo del paquete
    IF v_total_participantes > v_paquete.cupo_maximo THEN
        compra_id := 0;
        salida_id := 0;
        precio_total := 0;
        mensaje := 'La cantidad de participantes excede el cupo máximo del paquete';
        success := FALSE;
        RETURN NEXT;
        RETURN;
    END IF;

    -- Buscar/crear salida
    IF p_tipo_compra = 'compartido' THEN
        -- Buscar una salida compartida existente con cupo disponible
        SELECT s.id
        INTO v_salida_id
        FROM paquete_salidas_habilitadas s
        WHERE s.paquete_id = p_paquete_id
          AND s.fecha_salida = p_fecha_seleccionada
          AND s.tipo_salida = 'compartido'
          AND s.estado IN ('pendiente', 'activa')
          AND (s.cupo_maximo - s.cupos_reservados - s.cupos_confirmados) >= v_total_participantes
        ORDER BY s.id
        FOR UPDATE
        LIMIT 1;

        IF NOT FOUND THEN
            -- No existe salida, crear una nueva salida compartida
            -- Primero verificar capacidad máxima de salidas
            SELECT COUNT(*)
            INTO v_salidas_dia
            FROM paquete_salidas_habilitadas s
            JOIN paquetes_turisticos p ON p.id = s.paquete_id
            WHERE p.agencia_id = v_paquete.agencia_id
              AND s.fecha_salida = p_fecha_seleccionada
              AND s.estado IN ('pendiente', 'activa');

            IF v_salidas_dia >= v_max_salidas_por_dia THEN
                compra_id := 0;
                salida_id := 0;
                precio_total := 0;
                mensaje := 'La agencia alcanzó su máximo de salidas para ese día';
                success := FALSE;
                RETURN NEXT;
                RETURN;
            END IF;

            SELECT COUNT(*)
            INTO v_salidas_horario
            FROM paquete_salidas_habilitadas s
            JOIN paquetes_turisticos p ON p.id = s.paquete_id
            WHERE p.agencia_id = v_paquete.agencia_id
              AND s.fecha_salida = p_fecha_seleccionada
              AND s.estado IN ('pendiente', 'activa')
              AND COALESCE(p.horario, 'todo_dia') = v_horario_capacidad;

            IF v_salidas_horario >= v_max_salidas_por_horario THEN
                compra_id := 0;
                salida_id := 0;
                precio_total := 0;
                mensaje := 'La agencia alcanzó su máximo de salidas simultáneas para ese horario';
                success := FALSE;
                RETURN NEXT;
                RETURN;
            END IF;

            -- Crear nueva salida compartida
            INSERT INTO paquete_salidas_habilitadas (
                paquete_id,
                fecha_salida,
                tipo_salida,
                cupo_minimo,
                cupo_maximo,
                cupos_reservados,
                cupos_confirmados,
                estado,
                created_at,
                updated_at
            ) VALUES (
                p_paquete_id,
                p_fecha_seleccionada,
                'compartido',
                v_paquete.cupo_minimo,
                v_paquete.cupo_maximo,
                v_total_participantes,
                0,
                'pendiente',
                CURRENT_TIMESTAMP,
                CURRENT_TIMESTAMP
            ) RETURNING id INTO v_salida_id;
        ELSE
            -- Salida existente encontrada, actualizar cupos reservados
            UPDATE paquete_salidas_habilitadas
            SET cupos_reservados = cupos_reservados + v_total_participantes,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = v_salida_id;
        END IF;
    ELSE
        -- Privado: siempre crea una salida exclusiva
        SELECT COUNT(*)
        INTO v_salidas_dia
        FROM paquete_salidas_habilitadas s
        JOIN paquetes_turisticos p ON p.id = s.paquete_id
        WHERE p.agencia_id = v_paquete.agencia_id
          AND s.fecha_salida = p_fecha_seleccionada
          AND s.estado IN ('pendiente', 'activa');

        IF v_salidas_dia >= v_max_salidas_por_dia THEN
            compra_id := 0;
            salida_id := 0;
            precio_total := 0;
            mensaje := 'La agencia alcanzó su máximo de salidas para ese día';
            success := FALSE;
            RETURN NEXT;
            RETURN;
        END IF;

        SELECT COUNT(*)
        INTO v_salidas_horario
        FROM paquete_salidas_habilitadas s
        JOIN paquetes_turisticos p ON p.id = s.paquete_id
        WHERE p.agencia_id = v_paquete.agencia_id
          AND s.fecha_salida = p_fecha_seleccionada
          AND s.estado IN ('pendiente', 'activa')
          AND COALESCE(p.horario, 'todo_dia') = v_horario_capacidad;

        IF v_salidas_horario >= v_max_salidas_por_horario THEN
            compra_id := 0;
            salida_id := 0;
            precio_total := 0;
            mensaje := 'La agencia alcanzó su máximo de salidas simultáneas para ese horario';
            success := FALSE;
            RETURN NEXT;
            RETURN;
        END IF;

        INSERT INTO paquete_salidas_habilitadas (
            paquete_id,
            fecha_salida,
            tipo_salida,
            cupo_minimo,
            cupo_maximo,
            cupos_reservados,
            cupos_confirmados,
            estado,
            created_at,
            updated_at
        ) VALUES (
            p_paquete_id,
            p_fecha_seleccionada,
            'privado',
            v_total_participantes,
            v_total_participantes,
            v_total_participantes,
            0,
            'pendiente',
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        ) RETURNING id INTO v_salida_id;
    END IF;

    -- Calcular precios
    v_precio_base := COALESCE(v_paquete.precio_base_nacionales, 0);

    IF p_extranjero THEN
        v_recargo_extranjero_total := COALESCE(v_paquete.precio_adicional_extranjeros, 0) * v_personas_pagan;
    ELSE
        v_recargo_extranjero_total := 0;
    END IF;

    v_subtotal_base := v_precio_base * v_personas_pagan;
    v_subtotal_full := v_subtotal_base + v_recargo_extranjero_total;

    IF p_tipo_compra = 'privado' AND v_recargo_privado_porcentaje > 0 THEN
        v_recargo_privado := v_subtotal_full * (v_recargo_privado_porcentaje / 100);
    ELSE
        v_recargo_privado := 0;
    END IF;

    v_total_recargo := v_recargo_extranjero_total + v_recargo_privado;
    v_precio_total := v_subtotal_base + v_total_recargo;

    -- Crear compra
    INSERT INTO compras_paquetes (
        turista_id,
        paquete_id,
        salida_id,
        fecha_compra,
        fecha_seleccionada,
        horario_seleccionado,
        tipo_compra,
        extranjero,
        cantidad_adultos,
        cantidad_ninos_pagan,
        cantidad_ninos_gratis,
        total_participantes,
        precio_unitario,
        recargo_privado_porcentaje,
        recargo_extranjero,
        subtotal,
        total_recargo,
        precio_total,
        tiene_discapacidad,
        descripcion_discapacidad,
        notas_turista,
        status,
        created_at,
        updated_at
    ) VALUES (
        p_turista_id,
        p_paquete_id,
        v_salida_id,
        CURRENT_TIMESTAMP,
        p_fecha_seleccionada,
        v_horario_seleccionado,
        p_tipo_compra,
        COALESCE(p_extranjero, FALSE),
        p_cantidad_adultos,
        COALESCE(p_cantidad_ninos_pagan, 0),
        COALESCE(p_cantidad_ninos_gratis, 0),
        v_total_participantes,
        v_precio_base,
        CASE WHEN p_tipo_compra = 'privado' THEN v_recargo_privado_porcentaje ELSE 0 END,
        v_recargo_extranjero_total,
        v_subtotal_base,
        v_total_recargo,
        v_precio_total,
        COALESCE(p_tiene_discapacidad, FALSE),
        p_descripcion_discapacidad,
        p_notas_turista,
        'pendiente_confirmacion',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ) RETURNING id INTO v_compra_id;

    compra_id := v_compra_id;
    salida_id := v_salida_id;
    precio_total := v_precio_total;
    mensaje := 'Compra registrada. Esperando confirmación de pago.';
    success := TRUE;
    RETURN NEXT;
    RETURN;
END;
$$;
`


func ensureAgenciaMetricsFunction(db *gorm.DB) error {
	if err := db.Exec(sqlAgenciaMetrics).Error; err != nil {
		return fmt.Errorf("agencia metrics bootstrap failed: %w", err)
	}

	return nil
}

const sqlAgenciaMetrics = `
CREATE OR REPLACE FUNCTION public.get_agencia_metrics(p_agencia_id INT, p_mes INT, p_anio INT)
RETURNS JSON AS $$
DECLARE
    v_result JSON;
BEGIN
    SELECT json_build_object(
        'ingresos_mes', (
            SELECT COALESCE(SUM(precio_total), 0)
            FROM compras_paquetes cp
            JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
            WHERE pt.agencia_id = p_agencia_id
              AND cp.status = 'confirmada'
              AND EXTRACT(MONTH FROM cp.fecha_confirmacion) = p_mes
              AND EXTRACT(YEAR FROM cp.fecha_confirmacion) = p_anio
        ),
        'ventas_confirmadas', (
            SELECT COUNT(*)
            FROM compras_paquetes cp
            JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
            WHERE pt.agencia_id = p_agencia_id
              AND cp.status = 'confirmada'
              AND EXTRACT(MONTH FROM cp.fecha_confirmacion) = p_mes
              AND EXTRACT(YEAR FROM cp.fecha_confirmacion) = p_anio
        ),
        'turistas_atendidos', (
            SELECT COALESCE(SUM(total_participantes), 0)
            FROM compras_paquetes cp
            JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
            WHERE pt.agencia_id = p_agencia_id
              AND cp.status = 'confirmada'
              AND EXTRACT(MONTH FROM cp.fecha_confirmacion) = p_mes
              AND EXTRACT(YEAR FROM cp.fecha_confirmacion) = p_anio
        ),
        'pendientes_pago', (
            SELECT COUNT(*)
            FROM pagos_compras pc
            JOIN compras_paquetes cp ON pc.compra_id = cp.id
            JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
            WHERE pt.agencia_id = p_agencia_id
              AND pc.estado = 'pendiente'
        ),
        'top_paquetes', (
            SELECT COALESCE(json_agg(row ORDER BY row.ventas DESC), '[]'::json)
            FROM (
                SELECT
                    pt.id AS paquete_id,
                    pt.nombre AS nombre,
                    COUNT(*) AS ventas,
                    SUM(cp.precio_total) AS ingresos
                FROM compras_paquetes cp
                JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
                WHERE pt.agencia_id = p_agencia_id
                  AND cp.status = 'confirmada'
                  AND EXTRACT(MONTH FROM cp.fecha_confirmacion) = p_mes
                  AND EXTRACT(YEAR FROM cp.fecha_confirmacion) = p_anio
                GROUP BY pt.id, pt.nombre
                ORDER BY ventas DESC
                LIMIT 5
            ) row
        )
    ) INTO v_result;

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;
`

func ensureTriggerPagoConfirmado(db *gorm.DB) error {
	if err := db.Exec(sqlTriggerPagoConfirmado).Error; err != nil {
		return fmt.Errorf("trigger pago confirmado bootstrap failed: %w", err)
	}

	return nil
}

const sqlTriggerPagoConfirmado = `
-- Función trigger para cuando un pago es confirmado
CREATE OR REPLACE FUNCTION public.fn_on_pago_confirmado()
RETURNS TRIGGER AS $$
DECLARE
    v_compra RECORD;
    v_paquete_nombre VARCHAR(255);
    v_fecha_salida DATE;
    v_confirmado_por_nombre TEXT;
    v_notif_id INTEGER;
BEGIN
    -- Solo ejecutar cuando el estado cambia a 'confirmado'
    IF NEW.estado = 'confirmado' AND OLD.estado = 'pendiente' THEN
        -- Obtener datos de la compra
        SELECT id, salida_id, total_participantes, status, turista_id, paquete_id
        INTO v_compra
        FROM compras_paquetes
        WHERE id = NEW.compra_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Compra no encontrada: %', NEW.compra_id;
        END IF;

        -- Solo procesar si la compra está pendiente de confirmación
        IF v_compra.status = 'pendiente_confirmacion' THEN
            -- Actualizar estado de la compra
            UPDATE compras_paquetes
            SET status = 'confirmada',
                fecha_confirmacion = NOW(),
                codigo_confirmacion = COALESCE(codigo_confirmacion, 'CONF-' || LPAD(NEW.compra_id::text, 6, '0')),
                updated_at = NOW()
            WHERE id = NEW.compra_id;

            -- Mover cupos de reservados a confirmados
            IF v_compra.salida_id IS NOT NULL THEN
                UPDATE paquete_salidas_habilitadas
                SET cupos_reservados = GREATEST(0, cupos_reservados - v_compra.total_participantes),
                    cupos_confirmados = cupos_confirmados + v_compra.total_participantes,
                    updated_at = NOW()
                WHERE id = v_compra.salida_id;

                -- Obtener fecha de salida
                SELECT fecha_salida INTO v_fecha_salida
                FROM paquete_salidas_habilitadas
                WHERE id = v_compra.salida_id;
            END IF;

            -- Obtener nombre del paquete
            SELECT nombre INTO v_paquete_nombre
            FROM paquetes_turisticos
            WHERE id = v_compra.paquete_id;

            -- Obtener nombre de quien confirmó
            IF NEW.confirmado_por IS NOT NULL THEN
                SELECT nombre || ' ' || apellido_paterno INTO v_confirmado_por_nombre
                FROM usuarios
                WHERE id = NEW.confirmado_por;
            END IF;

            -- Crear notificación para el turista
            INSERT INTO notificaciones (
                usuario_id,
                tipo,
                titulo,
                mensaje,
                datos_json
            ) VALUES (
                v_compra.turista_id,
                'pago_confirmado',
                '¡Tu pago fue confirmado!',
                'Tu pago de Bs ' || NEW.monto || ' para "' || v_paquete_nombre || '" fue confirmado exitosamente',
                jsonb_build_object(
                    'pago_id', NEW.id,
                    'compra_id', NEW.compra_id,
                    'paquete_id', v_compra.paquete_id,
                    'paquete_nombre', v_paquete_nombre,
                    'confirmado_por', v_confirmado_por_nombre,
                    'fecha_salida', v_fecha_salida,
                    'monto', NEW.monto
                )
            ) RETURNING id INTO v_notif_id;

            -- Notificar vía PostgreSQL NOTIFY
            PERFORM pg_notify('notificaciones', json_build_object(
                'usuario_id', v_compra.turista_id,
                'notificacion_id', v_notif_id
            )::text);
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear el trigger si no existe
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'trg_pago_confirmado'
    ) THEN
        CREATE TRIGGER trg_pago_confirmado
            AFTER UPDATE ON pagos_compras
            FOR EACH ROW
            EXECUTE FUNCTION fn_on_pago_confirmado();
    END IF;
END $$;
`

// Trigger para rechazo de pago
const triggerPagoRechazadoSignature = "public.fn_on_pago_rechazado()"

func ensureTriggerPagoRechazado(db *gorm.DB) error {
	var exists bool
	if err := db.Raw(`SELECT to_regprocedure(?::text) IS NOT NULL`, triggerPagoRechazadoSignature).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check fn_on_pago_rechazado existence: %w", err)
	}
	if exists {
		return nil
	}

	if err := db.Exec(sqlTriggerPagoRechazado).Error; err != nil {
		return fmt.Errorf("trigger pago rechazado bootstrap failed: %w", err)
	}

	return nil
}

const sqlTriggerPagoRechazado = `
-- Función trigger para cuando un pago es rechazado
CREATE OR REPLACE FUNCTION public.fn_on_pago_rechazado()
RETURNS TRIGGER AS $$
DECLARE
    v_compra RECORD;
    v_paquete_nombre VARCHAR(255);
    v_notif_id INTEGER;
BEGIN
    -- Solo ejecutar cuando el estado cambia a 'rechazado'
    IF NEW.estado = 'rechazado' AND OLD.estado = 'pendiente' THEN
        -- Obtener datos de la compra
        SELECT id, salida_id, total_participantes, status, turista_id, paquete_id
        INTO v_compra
        FROM compras_paquetes
        WHERE id = NEW.compra_id;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'Compra no encontrada: %', NEW.compra_id;
        END IF;

        -- Solo procesar si la compra está pendiente
        IF v_compra.status = 'pendiente_confirmacion' THEN
            -- Marcar compra como rechazada
            UPDATE compras_paquetes
            SET status = 'rechazada',
                razon_rechazo = NEW.razon_rechazo,
                fecha_rechazo = NOW(),
                updated_at = NOW()
            WHERE id = NEW.compra_id;

            -- Liberar cupos reservados
            IF v_compra.salida_id IS NOT NULL THEN
                UPDATE paquete_salidas_habilitadas
                SET cupos_reservados = GREATEST(0, cupos_reservados - v_compra.total_participantes),
                    updated_at = NOW()
                WHERE id = v_compra.salida_id;
            END IF;

            -- Obtener nombre del paquete
            SELECT nombre INTO v_paquete_nombre
            FROM paquetes_turisticos
            WHERE id = v_compra.paquete_id;

            -- Crear notificación para el turista
            INSERT INTO notificaciones (
                usuario_id,
                tipo,
                titulo,
                mensaje,
                datos_json
            ) VALUES (
                v_compra.turista_id,
                'pago_rechazado',
                'Tu pago fue rechazado',
                'Tu pago de Bs ' || NEW.monto || ' para "' || v_paquete_nombre || '" fue rechazado. Razón: ' || COALESCE(NEW.razon_rechazo, 'No especificada'),
                jsonb_build_object(
                    'pago_id', NEW.id,
                    'compra_id', NEW.compra_id,
                    'paquete_id', v_compra.paquete_id,
                    'paquete_nombre', v_paquete_nombre,
                    'razon_rechazo', NEW.razon_rechazo,
                    'puede_reintentar', false,
                    'monto', NEW.monto
                )
            ) RETURNING id INTO v_notif_id;

            -- Notificar vía PostgreSQL NOTIFY
            PERFORM pg_notify('notificaciones', json_build_object(
                'usuario_id', v_compra.turista_id,
                'notificacion_id', v_notif_id
            )::text);
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear el trigger si no existe
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'trg_pago_rechazado'
    ) THEN
        CREATE TRIGGER trg_pago_rechazado
            AFTER UPDATE ON pagos_compras
            FOR EACH ROW
            EXECUTE FUNCTION fn_on_pago_rechazado();
    END IF;
END $$;
`

// ensureNotificacionesTable crea la tabla de notificaciones si no existe
func ensureNotificacionesTable(db *gorm.DB) error {
	const sqlCreateNotificaciones = `
CREATE TABLE IF NOT EXISTS notificaciones (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    tipo VARCHAR(50) NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    mensaje TEXT NOT NULL,
    datos_json JSONB,
    leida BOOLEAN DEFAULT FALSE,
    fecha_leida TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notificaciones_usuario_leida
    ON notificaciones(usuario_id, leida);

CREATE INDEX IF NOT EXISTS idx_notificaciones_created_at
    ON notificaciones(created_at DESC);
`
	return db.Exec(sqlCreateNotificaciones).Error
}

// ensureTriggerNuevoPago crea el trigger para notificar cuando se registra un nuevo pago
func ensureTriggerNuevoPago(db *gorm.DB) error {
	const sqlFunctionNuevoPago = `
CREATE OR REPLACE FUNCTION fn_notificar_nuevo_pago()
RETURNS TRIGGER AS $$
DECLARE
    v_encargado_id INTEGER;
    v_paquete_id INTEGER;
    v_paquete_nombre VARCHAR(255);
    v_turista_nombre TEXT;
    v_notif_id INTEGER;
BEGIN
    -- Obtener datos de la compra, paquete y turista
    SELECT
        pt.id,
        pt.nombre,
        pt.agencia_id,
        u.nombre || ' ' || u.apellido_paterno AS turista_nombre
    INTO
        v_paquete_id,
        v_paquete_nombre,
        v_encargado_id,
        v_turista_nombre
    FROM compras_paquetes cp
    JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
    JOIN agencias_turismo a ON pt.agencia_id = a.id
    JOIN usuarios u ON cp.turista_id = u.id
    WHERE cp.id = NEW.compra_id;

    -- Obtener encargado principal de la agencia
    SELECT encargado_principal_id INTO v_encargado_id
    FROM agencias_turismo a
    JOIN paquetes_turisticos pt ON a.id = pt.agencia_id
    JOIN compras_paquetes cp ON pt.id = cp.paquete_id
    WHERE cp.id = NEW.compra_id;

    -- Solo crear notificación si hay encargado asignado
    IF v_encargado_id IS NOT NULL THEN
        -- Insertar notificación en la base de datos
        INSERT INTO notificaciones (
            usuario_id,
            tipo,
            titulo,
            mensaje,
            datos_json
        ) VALUES (
            v_encargado_id,
            'nuevo_pago_pendiente',
            'Nuevo pago pendiente de confirmación',
            v_turista_nombre || ' registró un pago de Bs ' || NEW.monto || ' para "' || v_paquete_nombre || '"',
            jsonb_build_object(
                'pago_id', NEW.id,
                'compra_id', NEW.compra_id,
                'paquete_id', v_paquete_id,
                'turista_nombre', v_turista_nombre,
                'monto', NEW.monto,
                'metodo_pago', NEW.metodo_pago,
                'comprobante_foto', NEW.comprobante_foto
            )
        ) RETURNING id INTO v_notif_id;

        -- Notificar vía PostgreSQL NOTIFY para WebSocket
        PERFORM pg_notify('notificaciones', json_build_object(
            'usuario_id', v_encargado_id,
            'notificacion_id', v_notif_id
        )::text);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`

	const sqlTriggerNuevoPago = `
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'trg_notificar_nuevo_pago'
    ) THEN
        CREATE TRIGGER trg_notificar_nuevo_pago
            AFTER INSERT ON pagos_compras
            FOR EACH ROW
            EXECUTE FUNCTION fn_notificar_nuevo_pago();
    END IF;
END $$;
`

	if err := db.Exec(sqlFunctionNuevoPago).Error; err != nil {
		return err
	}

	return db.Exec(sqlTriggerNuevoPago).Error
}
