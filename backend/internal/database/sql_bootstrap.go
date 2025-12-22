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

	if err := ensureTriggerPagoConfirmado(db); err != nil {
		return err
	}

	if err := ensureTriggerPagoRechazado(db); err != nil {
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
        RETURN QUERY SELECT 0, 0, 0, 'Turista no encontrado', FALSE;
        RETURN;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM usuarios WHERE id = p_turista_id AND rol = 'turista') THEN
        RETURN QUERY SELECT 0, 0, 0, 'Solo usuarios con rol "turista" pueden realizar compras', FALSE;
        RETURN;
    END IF;

    -- Validaciones básicas de participantes
    IF p_cantidad_adultos IS NULL OR p_cantidad_adultos < 1 THEN
        RETURN QUERY SELECT 0, 0, 0, 'Debe haber al menos 1 adulto', FALSE;
        RETURN;
    END IF;

    IF COALESCE(p_cantidad_ninos_pagan, 0) < 0 OR COALESCE(p_cantidad_ninos_gratis, 0) < 0 THEN
        RETURN QUERY SELECT 0, 0, 0, 'Las cantidades de niños no pueden ser negativas', FALSE;
        RETURN;
    END IF;

    v_total_participantes := p_cantidad_adultos
        + COALESCE(p_cantidad_ninos_pagan, 0)
        + COALESCE(p_cantidad_ninos_gratis, 0);
    v_personas_pagan := p_cantidad_adultos + COALESCE(p_cantidad_ninos_pagan, 0);

    IF v_total_participantes < 1 THEN
        RETURN QUERY SELECT 0, 0, 0, 'Debe registrar al menos 1 participante', FALSE;
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
        RETURN QUERY SELECT 0, 0, 0, 'Paquete no encontrado', FALSE;
        RETURN;
    END IF;

    IF v_paquete.status <> 'activo' OR v_paquete.visible_publico <> TRUE THEN
        RETURN QUERY SELECT 0, 0, 0, 'El paquete no está disponible', FALSE;
        RETURN;
    END IF;

    IF v_paquete.agencia_status <> 'activa' OR v_paquete.agencia_visible <> TRUE THEN
        RETURN QUERY SELECT 0, 0, 0, 'La agencia no está disponible', FALSE;
        RETURN;
    END IF;

    -- Validar tipo de compra
    IF p_tipo_compra NOT IN ('compartido', 'privado') THEN
        RETURN QUERY SELECT 0, 0, 0, 'Tipo de compra inválido', FALSE;
        RETURN;
    END IF;

    IF p_tipo_compra = 'privado' THEN
        IF v_paquete.frecuencia <> 'salida_diaria' THEN
            RETURN QUERY SELECT 0, 0, 0, 'El tipo privado solo está disponible para paquetes de salida diaria', FALSE;
            RETURN;
        END IF;
        IF v_paquete.permite_privado IS DISTINCT FROM TRUE THEN
            RETURN QUERY SELECT 0, 0, 0, 'Este paquete no permite compras privadas', FALSE;
            RETURN;
        END IF;
    END IF;

    -- Validar fecha seleccionada
    IF p_fecha_seleccionada < (CURRENT_DATE + COALESCE(v_paquete.dias_previos_compra, 1)::int) THEN
        RETURN QUERY SELECT 0, 0, 0, 'La fecha seleccionada no cumple los días previos de compra', FALSE;
        RETURN;
    END IF;

    IF v_paquete.frecuencia = 'salida_unica' THEN
        IF v_paquete.fecha_salida_fija IS NULL THEN
            RETURN QUERY SELECT 0, 0, 0, 'El paquete no tiene fecha de salida fija configurada', FALSE;
            RETURN;
        END IF;
        IF p_fecha_seleccionada <> v_paquete.fecha_salida_fija::date THEN
            RETURN QUERY SELECT 0, 0, 0, 'La compra debe realizarse en la fecha de salida fija del paquete', FALSE;
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
        RETURN QUERY SELECT 0, 0, 0, 'La cantidad de participantes excede el cupo máximo del paquete', FALSE;
        RETURN;
    END IF;

    -- Buscar/crear salida
    IF p_tipo_compra = 'compartido' THEN
        -- Si no existe una salida previa para esa fecha, exigir cupo mínimo para habilitar la primera compra compartida.
        IF NOT EXISTS (
            SELECT 1
            FROM paquete_salidas_habilitadas s0
            WHERE s0.paquete_id = p_paquete_id
              AND s0.fecha_salida = p_fecha_seleccionada
              AND s0.tipo_salida = 'compartido'
              AND s0.estado IN ('pendiente', 'activa')
        ) THEN
            IF v_total_participantes < v_paquete.cupo_minimo THEN
                RETURN QUERY SELECT 0, 0, 0, format('Para habilitar la primera salida en esta fecha debe registrar al menos %s participantes', v_paquete.cupo_minimo), FALSE;
                RETURN;
            END IF;
        END IF;

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
            -- Validar capacidad de la agencia solo si se creará una salida nueva
            SELECT COUNT(*)
            INTO v_salidas_dia
            FROM paquete_salidas_habilitadas s
            JOIN paquetes_turisticos p ON p.id = s.paquete_id
            WHERE p.agencia_id = v_paquete.agencia_id
              AND s.fecha_salida = p_fecha_seleccionada
              AND s.estado IN ('pendiente', 'activa');

            IF v_salidas_dia >= v_max_salidas_por_dia THEN
                RETURN QUERY SELECT 0, 0, 0, 'La agencia alcanzó su máximo de salidas para ese día', FALSE;
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
                RETURN QUERY SELECT 0, 0, 0, 'La agencia alcanzó su máximo de salidas simultáneas para ese horario', FALSE;
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
            RETURN QUERY SELECT 0, 0, 0, 'La agencia alcanzó su máximo de salidas para ese día', FALSE;
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
            RETURN QUERY SELECT 0, 0, 0, 'La agencia alcanzó su máximo de salidas simultáneas para ese horario', FALSE;
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

    RETURN QUERY SELECT v_compra_id, v_salida_id, v_precio_total, 'Compra registrada. Esperando confirmación de pago.', TRUE;
    RETURN;
END;
$$;
`

// Trigger para confirmación de pago
const triggerPagoConfirmadoSignature = "public.fn_on_pago_confirmado()"

func ensureTriggerPagoConfirmado(db *gorm.DB) error {
	var exists bool
	if err := db.Raw(`SELECT to_regprocedure(?::text) IS NOT NULL`, triggerPagoConfirmadoSignature).Scan(&exists).Error; err != nil {
		return fmt.Errorf("failed to check fn_on_pago_confirmado existence: %w", err)
	}
	if exists {
		return nil
	}

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
BEGIN
    -- Solo ejecutar cuando el estado cambia a 'confirmado'
    IF NEW.estado = 'confirmado' AND OLD.estado = 'pendiente' THEN
        -- Obtener datos de la compra
        SELECT id, salida_id, total_participantes, status
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
                updated_at = NOW()
            WHERE id = NEW.compra_id;

            -- Mover cupos de reservados a confirmados
            IF v_compra.salida_id IS NOT NULL THEN
                UPDATE paquete_salidas_habilitadas
                SET cupos_reservados = GREATEST(0, cupos_reservados - v_compra.total_participantes),
                    cupos_confirmados = cupos_confirmados + v_compra.total_participantes,
                    updated_at = NOW()
                WHERE id = v_compra.salida_id;
            END IF;
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
BEGIN
    -- Solo ejecutar cuando el estado cambia a 'rechazado'
    IF NEW.estado = 'rechazado' AND OLD.estado = 'pendiente' THEN
        -- Obtener datos de la compra
        SELECT id, salida_id, total_participantes, status
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
