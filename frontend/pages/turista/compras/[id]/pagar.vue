<template>
  <div class="page-shell">
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">{{ pageTitle }}</h1>
          <p class="muted mt-1">{{ pageSubtitle }}</p>
        </div>
        <div class="flex gap-2">
          <Button
            label="Mis compras"
            icon="pi pi-list"
            severity="secondary"
            outlined
            :disabled="loading"
            @click="navigateTo('/turista/mis-compras')"
          />
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 space-y-6">
      <div v-if="loading">
        <Skeleton height="220px" class="mb-4" />
        <Skeleton height="320px" />
      </div>

      <div v-else-if="error" class="text-center space-y-4 py-12">
        <i class="pi pi-exclamation-triangle text-5xl text-orange-500"></i>
        <h2 class="text-2xl font-bold text-gray-900">No se pudo cargar la compra</h2>
        <p class="muted">{{ error }}</p>
        <div class="flex justify-center gap-2">
          <Button label="Volver" icon="pi pi-arrow-left" severity="secondary" outlined @click="navigateTo('/turista/mis-compras')" />
          <Button label="Reintentar" icon="pi pi-refresh" @click="loadAll" />
        </div>
      </div>

      <div v-else-if="compra" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <Card v-if="isConfirmada" class="surface-card border border-emerald-100 bg-emerald-50/60">
            <template #title>
              <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div class="space-y-1">
                  <div class="flex items-center gap-2">
                    <i class="pi pi-check-circle text-emerald-600"></i>
                    <span>Confirmacion de compra</span>
                  </div>
                  <div class="text-xs text-gray-600">
                    <span class="font-semibold text-emerald-700">Estado confirmado</span>
                  </div>
                </div>
                <ExportarConfirmacionPDF
                  :codigo="codigoConfirmacion"
                  :rows="pdfRows"
                  :file-name="pdfFileName"
                  :disabled="loading"
                />
              </div>
            </template>
            <template #content>
              <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <p class="text-xs text-gray-500">Codigo de confirmacion</p>
                  <p class="text-lg font-semibold text-gray-900">{{ codigoConfirmacion }}</p>
                </div>
                <p class="text-sm text-gray-600">
                  Guarda esta confirmacion para presentarla el dia de la salida.
                </p>
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-calendar text-emerald-600"></i>
                <span>Informacion de la salida</span>
              </div>
            </template>
            <template #content>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <div>
                  <p class="text-xs text-gray-500">Nombre del paquete</p>
                  <p class="font-semibold text-gray-900">{{ compra.paquete?.nombre || 'Paquete' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha y hora de salida</p>
                  <p class="font-semibold text-gray-900">{{ fechaHoraSalidaLabel }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Punto de encuentro</p>
                  <p class="font-semibold text-gray-900">{{ displayValue(puntoEncuentroLabel) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Duracion</p>
                  <p class="font-semibold text-gray-900">{{ displayValue(duracionLabel) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Horario</p>
                  <p class="font-semibold text-gray-900">{{ displayValue(horarioLabel) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Dificultad</p>
                  <p class="font-semibold text-gray-900">{{ displayValue(dificultadLabel) }}</p>
                </div>
              </div>

              <div v-if="instruccionesTuristas" class="mt-4 rounded-lg border border-gray-200 bg-gray-50 px-4 py-3">
                <p class="text-xs text-gray-500">Instrucciones</p>
                <p class="text-sm text-gray-700 whitespace-pre-line">{{ instruccionesTuristas }}</p>
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-ticket text-emerald-600"></i>
                <span>Tu reserva</span>
              </div>
            </template>
            <template #content>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-gray-500">Tipo</p>
                  <Tag :value="tipoCompraLabel" severity="secondary" />
                </div>
                <div>
                  <p class="text-xs text-gray-500">Total pagado</p>
                  <p class="text-2xl font-bold text-emerald-700">Bs. {{ formatMoney(compra.precio_total) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha de compra</p>
                  <p class="font-semibold text-gray-900">{{ formatFechaHora(compra.fecha_compra) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha de confirmacion</p>
                  <p class="font-semibold text-gray-900">{{ fechaConfirmacionLabel }}</p>
                </div>
              </div>

              <div class="mt-4 rounded-lg border border-gray-200 bg-white px-4 py-3">
                <p class="text-xs text-gray-500 mb-2">Participantes</p>
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
                  <div>
                    <p class="text-xs text-gray-500">Adultos</p>
                    <p class="font-semibold text-gray-900">{{ adultosLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Ninos pagan</p>
                    <p class="font-semibold text-gray-900">{{ ninosPaganLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Ninos gratis</p>
                    <p class="font-semibold text-gray-900">{{ ninosGratisLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Total</p>
                    <p class="font-semibold text-gray-900">{{ totalParticipantesLabel }}</p>
                  </div>
                </div>
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada && paquete" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-sparkles text-blue-600"></i>
                <span>Detalle del paquete</span>
              </div>
            </template>
            <template #content>
              <div class="space-y-6">
                <div>
                  <p class="text-sm text-gray-700 whitespace-pre-line">
                    {{ paquete.descripcion || 'No hay descripcion disponible.' }}
                  </p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div>
                    <p class="text-sm font-semibold text-gray-900 mb-2">Incluye</p>
                    <ul v-if="(paquete.incluye || []).length" class="space-y-2">
                      <li v-for="(item, idx) in paquete.incluye" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                        <i class="pi pi-check-circle text-emerald-600 mt-0.5"></i>
                        <span>{{ item }}</span>
                      </li>
                    </ul>
                    <p v-else class="text-sm text-gray-500">Sin especificar.</p>
                  </div>

                  <div>
                    <p class="text-sm font-semibold text-gray-900 mb-2">No incluye</p>
                    <ul v-if="(paquete.no_incluye || []).length" class="space-y-2">
                      <li v-for="(item, idx) in paquete.no_incluye" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                        <i class="pi pi-times-circle text-rose-500 mt-0.5"></i>
                        <span>{{ item }}</span>
                      </li>
                    </ul>
                    <p v-else class="text-sm text-gray-500">Sin especificar.</p>
                  </div>

                  <div>
                    <p class="text-sm font-semibold text-gray-900 mb-2">Que llevar</p>
                    <ul v-if="(paquete.que_llevar || []).length" class="space-y-2">
                      <li v-for="(item, idx) in paquete.que_llevar" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                        <i class="pi pi-briefcase text-slate-600 mt-0.5"></i>
                        <span>{{ item }}</span>
                      </li>
                    </ul>
                    <p v-else class="text-sm text-gray-500">Sin especificar.</p>
                  </div>
                </div>

                <div class="border-t border-gray-200 pt-4">
                  <p class="text-sm font-semibold text-gray-900 mb-2">Requisitos</p>
                  <ul v-if="requisitosItems.length" class="space-y-2">
                    <li v-for="(item, idx) in requisitosItems" :key="idx" class="flex items-start gap-2 text-sm text-gray-700">
                      <i class="pi pi-check text-emerald-600 mt-0.5"></i>
                      <span>{{ item }}</span>
                    </li>
                  </ul>
                  <p v-else class="text-sm text-gray-500">Sin requisitos registrados.</p>
                </div>

                <div class="border-t border-gray-200 pt-4">
                  <div class="flex items-center gap-2 mb-3">
                    <i class="pi pi-route text-emerald-600"></i>
                    <span class="text-sm font-semibold text-gray-900">Itinerario</span>
                  </div>
                  <div v-if="isMultiDay" class="space-y-4">
                    <div v-if="itinerarioOrdenado.length" class="space-y-3">
                      <div
                        v-for="item in itinerarioOrdenado"
                        :key="item.id"
                        class="rounded-xl border border-gray-200 bg-white p-4"
                      >
                        <div class="flex items-start justify-between gap-3">
                          <div>
                            <p class="text-xs uppercase tracking-wider text-emerald-700/80 font-semibold">Dia {{ item.dia_numero }}</p>
                            <h3 class="text-lg font-semibold text-gray-900 mt-1">{{ item.titulo }}</h3>
                          </div>
                          <Tag :value="`Dia ${item.dia_numero}`" severity="secondary" />
                        </div>

                        <p v-if="item.descripcion" class="text-sm text-gray-700 mt-2 whitespace-pre-line">
                          {{ item.descripcion }}
                        </p>

                        <div v-if="(item.actividades || []).length" class="mt-3">
                          <p class="text-sm font-semibold text-gray-900 mb-2">Actividades</p>
                          <ul class="grid grid-cols-1 md:grid-cols-2 gap-2">
                            <li v-for="(act, idx) in item.actividades" :key="idx" class="text-sm text-gray-700 flex items-start gap-2">
                              <i class="pi pi-check text-emerald-600 mt-0.5"></i>
                              <span>{{ act }}</span>
                            </li>
                          </ul>
                        </div>

                        <div v-if="item.hospedaje_info" class="mt-3 p-3 rounded-lg bg-slate-50 border border-slate-200">
                          <p class="text-xs text-slate-500 mb-1">Hospedaje</p>
                          <p class="text-sm text-slate-700 whitespace-pre-line">{{ item.hospedaje_info }}</p>
                        </div>
                      </div>
                    </div>

                    <div v-else class="text-sm text-gray-600">
                      Este paquete no tiene itinerario detallado aun.
                    </div>
                  </div>
                  <div v-else class="text-sm text-gray-700">
                    Este paquete es de un dia. Revisa los detalles de salida para el horario y punto de encuentro.
                  </div>
                </div>

                <div class="border-t border-gray-200 pt-4">
                  <div class="flex items-center gap-2 mb-3">
                    <i class="pi pi-map-marker text-blue-600"></i>
                    <span class="text-sm font-semibold text-gray-900">Orden de atracciones</span>
                  </div>
                  <div v-if="atraccionesOrdenadas.length" class="space-y-4">
                    <div v-if="isMultiDay" class="space-y-4">
                      <div v-for="dia in diasAtracciones" :key="dia" class="space-y-2">
                        <div class="flex items-center justify-between">
                          <p class="text-xs uppercase tracking-wider text-gray-500 font-semibold">Dia {{ dia }}</p>
                          <span class="text-xs text-gray-500">{{ (atraccionesPorDia[dia] || []).length }} atracciones</span>
                        </div>
                        <ul class="space-y-2">
                          <li
                            v-for="(item, idx) in atraccionesPorDia[dia]"
                            :key="item.id || `${dia}-${idx}`"
                            class="flex items-start gap-3 text-sm text-gray-700"
                          >
                            <span
                              class="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full bg-emerald-100 text-emerald-700 text-xs font-semibold"
                            >
                              {{ getOrdenVisitaLabel(atraccionesPorDia[dia], idx) }}
                            </span>
                            <div class="min-w-0">
                              <p class="font-semibold text-gray-900 truncate">
                                {{ item.atraccion?.nombre || `Atraccion #${item.atraccion_id}` }}
                              </p>
                              <p class="text-xs text-gray-500">
                                {{ item.atraccion?.provincia?.nombre || '' }}
                                <span v-if="item.atraccion?.provincia?.departamento?.nombre">
                                  , {{ item.atraccion.provincia.departamento.nombre }}
                                </span>
                              </p>
                            </div>
                          </li>
                        </ul>
                      </div>
                    </div>
                    <ul v-else class="space-y-2">
                      <li v-for="(item, idx) in atraccionesOrdenadas" :key="item.id || idx" class="flex items-start gap-3 text-sm text-gray-700">
                        <span
                          class="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full bg-emerald-100 text-emerald-700 text-xs font-semibold"
                        >
                          {{ getOrdenVisitaLabel(atraccionesOrdenadas, idx) }}
                        </span>
                        <div class="min-w-0">
                          <p class="font-semibold text-gray-900 truncate">
                            {{ item.atraccion?.nombre || `Atraccion #${item.atraccion_id}` }}
                          </p>
                          <p class="text-xs text-gray-500">
                            {{ item.atraccion?.provincia?.nombre || '' }}
                            <span v-if="item.atraccion?.provincia?.departamento?.nombre">
                              , {{ item.atraccion.provincia.departamento.nombre }}
                            </span>
                          </p>
                        </div>
                      </li>
                    </ul>
                  </div>
                  <p v-else class="text-sm text-gray-500">Sin atracciones registradas.</p>
                </div>

                <div class="border-t border-gray-200 pt-4">
                  <div class="flex items-center gap-2 mb-3">
                    <i class="pi pi-images text-blue-600"></i>
                    <span class="text-sm font-semibold text-gray-900">Galeria de fotos</span>
                  </div>
                  <div v-if="fotosOrdenadas.length" class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                    <a
                      v-for="(foto, idx) in fotosOrdenadas.slice(0, 9)"
                      :key="foto.id || idx"
                      :href="resolveAssetUrl(foto.foto)"
                      target="_blank"
                      rel="noopener"
                      class="relative group rounded-xl overflow-hidden border border-gray-200"
                    >
                      <img
                        :src="resolveAssetUrl(foto.foto)"
                        alt="Foto del paquete"
                        class="w-full h-28 object-cover group-hover:scale-105 transition-transform duration-300"
                        loading="lazy"
                      />
                      <div class="absolute inset-0 bg-black/0 group-hover:bg-black/15 transition-colors" />
                    </a>
                  </div>
                  <p v-else class="text-sm text-gray-500">Sin fotos registradas.</p>
                </div>
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada && tienePoliticas" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-file text-blue-600"></i>
                <span>Politicas</span>
              </div>
            </template>
            <template #content>
              <div class="text-sm text-gray-700 space-y-2">
                <div v-if="edadMinimaPago > 0" class="flex items-start gap-2">
                  <i class="pi pi-users text-emerald-600 mt-0.5"></i>
                  <span>Edad minima de pago: <strong>{{ edadMinimaPago }} anos</strong></span>
                </div>
                <div v-if="permitePrivado && recargoPrivado > 0" class="flex items-start gap-2">
                  <i class="pi pi-lock text-amber-600 mt-0.5"></i>
                  <span>Recargo privado: <strong>{{ formatMoney(recargoPrivado) }}%</strong></span>
                </div>
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada && compra.notas_turista" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-comment text-slate-600"></i>
                <span>Notas de compra</span>
              </div>
            </template>
            <template #content>
              <p class="text-sm text-gray-700 whitespace-pre-line">{{ compra.notas_turista }}</p>
            </template>
          </Card>

          <Card v-if="!isConfirmada" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-receipt text-emerald-600"></i>
                <span>Resumen</span>
              </div>
            </template>
            <template #content>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-gray-500">Paquete</p>
                  <p class="font-semibold text-gray-900">{{ compra.paquete?.nombre || `#${compra.paquete?.id}` }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Estado</p>
                  <Tag :value="statusLabel" :severity="statusSeverity" />
                </div>
                <div>
                  <p class="text-xs text-gray-500">Tipo de compra</p>
                  <p class="font-semibold text-gray-900">{{ tipoCompraLabel }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha de compra</p>
                  <p class="font-semibold text-gray-900">{{ formatFechaHora(compra.fecha_compra) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha seleccionada</p>
                  <p class="font-semibold text-gray-900">{{ formatFecha(compra.fecha_seleccionada) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Participantes</p>
                  <p class="font-semibold text-gray-900">{{ compra.total_participantes }}</p>
                </div>
                <div class="md:col-span-2">
                  <p class="text-xs text-gray-500">Total</p>
                  <p class="text-2xl font-bold text-emerald-700">Bs. {{ formatMoney(compra.precio_total) }}</p>
                </div>
              </div>
            </template>
          </Card>

          <Card v-if="paquete && !isConfirmada" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-map-marker text-sky-600"></i>
                <span>Detalle del paquete</span>
              </div>
            </template>
            <template #content>
              <div class="flex flex-col md:flex-row gap-4">
                <img
                  v-if="paqueteFotoUrl"
                  :src="paqueteFotoUrl"
                  alt="Foto del paquete"
                  class="w-full md:w-56 h-40 object-cover rounded-lg border border-gray-200"
                  loading="lazy"
                />
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 flex-1">
                  <div>
                    <p class="text-xs text-gray-500">Agencia</p>
                    <p class="font-semibold text-gray-900">{{ agenciaNombre }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Ubicacion</p>
                    <p class="font-semibold text-gray-900">{{ agenciaUbicacion }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Duracion</p>
                    <p class="font-semibold text-gray-900">{{ duracionLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Horario</p>
                    <p class="font-semibold text-gray-900">{{ horarioLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Hora salida</p>
                    <p class="font-semibold text-gray-900">{{ horaSalidaLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Dificultad</p>
                    <p class="font-semibold text-gray-900">{{ dificultadLabel }}</p>
                  </div>
                </div>
              </div>
            </template>
          </Card>

          <Message v-if="compra.status !== 'pendiente_confirmacion' && compra.status !== 'confirmada'" severity="info" :closable="false">
            Esta compra ya no admite registrar pagos. Estado: {{ statusLabel }}.
          </Message>

          <Message v-else-if="compra.ultimo_pago?.estado === 'pendiente'" severity="success" :closable="false">
            Ya registraste un pago pendiente para esta compra. Espera la confirmación del encargado.
          </Message>

          <SubirComprobante
            v-else-if="compra.status === 'pendiente_confirmacion'"
            :compra-id="compraId"
            :monto="Number(compra.precio_total || 0)"
            @pago-registrado="loadAll"
          />
        </div>

        <div class="space-y-6">
          <Card v-if="paquete && !isConfirmada" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-blue-600"></i>
                <span>Datos de pago</span>
              </div>
            </template>
            <template #content>
              <div v-if="paquete.agencia_datos_pago" class="space-y-3">
                <div v-if="paquete.agencia_datos_pago.nombre_banco" class="text-sm">
                  <span class="text-gray-500">Banco:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ paquete.agencia_datos_pago.nombre_banco }}</span>
                </div>
                <div v-if="paquete.agencia_datos_pago.numero_cuenta" class="text-sm">
                  <span class="text-gray-500">Cuenta:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ paquete.agencia_datos_pago.numero_cuenta }}</span>
                </div>
                <div v-if="paquete.agencia_datos_pago.nombre_titular" class="text-sm">
                  <span class="text-gray-500">Titular:</span>
                  <span class="font-semibold text-gray-900 ml-2">{{ paquete.agencia_datos_pago.nombre_titular }}</span>
                </div>

                <div v-if="paquete.agencia_datos_pago.qr_pago_foto" class="pt-3">
                  <p class="text-xs text-gray-500 mb-2">QR</p>
                  <img
                    :src="resolveAssetUrl(paquete.agencia_datos_pago.qr_pago_foto)"
                    alt="QR de pago"
                    class="w-full rounded-lg border border-gray-200"
                    loading="lazy"
                  />
                </div>
              </div>

              <div v-else class="text-sm text-gray-600">
                La agencia no tiene datos de pago configurados.
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada && paquete?.agencia" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-blue-600"></i>
                <span>Informacion de la agencia</span>
              </div>
            </template>
            <template #content>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                <div>
                  <p class="text-xs text-gray-500">Nombre comercial</p>
                  <p class="font-semibold text-gray-900">{{ displayValue(agenciaNombre) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Ubicacion / Direccion</p>
                  <p class="font-semibold text-gray-900">{{ displayValue(agenciaDireccionLabel) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Telefono</p>
                  <a
                    v-if="agenciaTelefonoLink"
                    :href="agenciaTelefonoLink"
                    class="font-semibold text-emerald-700 hover:underline"
                  >
                    {{ agenciaTelefono }}
                  </a>
                  <p v-else class="font-semibold text-gray-900">Por definir</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Email</p>
                  <a v-if="agenciaEmail" :href="`mailto:${agenciaEmail}`" class="font-semibold text-emerald-700 hover:underline">{{ agenciaEmail }}</a>
                  <p v-else class="font-semibold text-gray-900">Por definir</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">WhatsApp</p>
                  <a
                    v-if="agenciaWhatsappLink"
                    :href="agenciaWhatsappLink"
                    class="font-semibold text-emerald-700 hover:underline"
                    target="_blank"
                    rel="noopener"
                  >
                    {{ agenciaWhatsapp }}
                  </a>
                  <p v-else class="font-semibold text-gray-900">Por definir</p>
                </div>
              </div>
            </template>
          </Card>

          <Card v-if="compra?.status === 'pendiente_confirmacion'" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-clock text-emerald-600"></i>
                <span>Tiempo para pagar</span>
              </div>
            </template>
            <template #content>
              <div v-if="isPagoPendiente" class="text-sm text-gray-600">
                Ya registraste un pago pendiente. El contador esta detenido.
              </div>
              <div v-else>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-500">Tiempo restante</span>
                  <span class="text-xl font-semibold text-emerald-700">{{ countdownLabel }}</span>
                </div>
                <div class="mt-3 h-2 w-full rounded-full bg-gray-100 overflow-hidden">
                  <div class="h-full bg-emerald-500 transition-all" :style="{ width: `${countdownPercent}%` }"></div>
                </div>
                <p class="mt-3 text-xs text-gray-500">
                  Si no registras el pago en 5 minutos, la compra se cancelara y los cupos se liberaran.
                </p>
              </div>
              <p v-if="autoCancelling" class="mt-3 text-xs text-gray-500">
                Cancelando compra por tiempo agotado...
              </p>
            </template>
          </Card>

          <Card v-if="compra?.ultimo_pago" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-file text-emerald-600"></i>
                <span>Comprobante de pago</span>
              </div>
            </template>
            <template #content>
              <div class="space-y-3">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
                  <div>
                    <p class="text-xs text-gray-500">Metodo de pago</p>
                    <p class="font-semibold text-gray-900">{{ metodoPagoLabel }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Monto</p>
                    <p class="font-semibold text-gray-900">Bs. {{ formatMoney(compra.ultimo_pago?.monto) }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">Estado</p>
                    <Tag :value="pagoStatusLabel" :severity="pagoStatusSeverity" />
                  </div>
                  <div v-if="compra.ultimo_pago?.fecha_confirmacion">
                    <p class="text-xs text-gray-500">Confirmado</p>
                    <p class="font-semibold text-gray-900">{{ formatFechaHora(compra.ultimo_pago?.fecha_confirmacion) }}</p>
                  </div>
                  <div v-if="compra.ultimo_pago?.fecha_registro">
                    <p class="text-xs text-gray-500">Registrado</p>
                    <p class="font-semibold text-gray-900">{{ formatFechaHora(compra.ultimo_pago?.fecha_registro) }}</p>
                  </div>
                </div>

                <div v-if="compra.ultimo_pago?.comprobante_foto" class="pt-3">
                  <p class="text-xs text-gray-500 mb-2">Imagen del comprobante</p>
                  <img
                    :src="resolveAssetUrl(compra.ultimo_pago?.comprobante_foto)"
                    alt="Comprobante"
                    class="w-full rounded-lg border border-gray-200"
                    loading="lazy"
                  />
                </div>
                <p v-else class="text-sm text-gray-600">No hay comprobante adjunto.</p>
              </div>
            </template>
          </Card>

          <Card v-if="isConfirmada" class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-info-circle text-emerald-600"></i>
                <span>Instrucciones importantes</span>
              </div>
            </template>
            <template #content>
              <div class="space-y-4 text-sm text-gray-700">
                <div class="flex items-start gap-3">
                  <i class="pi pi-clock text-emerald-600 mt-0.5"></i>
                  <div>
                    <p class="font-semibold text-gray-900">Llegar 15 min antes</p>
                    <p class="text-gray-600">Te ayuda a coordinar el check-in y el abordaje.</p>
                  </div>
                </div>
                <div class="flex items-start gap-3">
                  <i class="pi pi-file text-blue-600 mt-0.5"></i>
                  <div>
                    <p class="font-semibold text-gray-900">Politicas de cancelacion</p>
                    <p class="text-gray-600 whitespace-pre-line">{{ politicaCancelacion || 'Por definir' }}</p>
                  </div>
                </div>
                <div class="flex items-start gap-3">
                  <i class="pi pi-phone text-emerald-600 mt-0.5"></i>
                  <div>
                    <p class="font-semibold text-gray-900">Contacto de emergencia</p>
                    <p class="text-gray-600">{{ displayValue(contactoEmergenciaLabel) }}</p>
                  </div>
                </div>
                <div class="flex items-start gap-3">
                  <i class="pi pi-cloud text-slate-600 mt-0.5"></i>
                  <div>
                    <p class="font-semibold text-gray-900">Recomendaciones climaticas</p>
                    <p class="text-gray-600 whitespace-pre-line">{{ displayValue(recomendacionesClimaticasLabel) }}</p>
                  </div>
                </div>
              </div>
            </template>
          </Card>

          <Card v-else class="surface-card">
            <template #title>
              <div class="flex items-center gap-2">
                <i class="pi pi-info-circle text-gray-600"></i>
                <span>Recordatorio</span>
              </div>
            </template>
            <template #content>
              <ul class="text-sm text-gray-700 space-y-2">
                <li>El monto debe coincidir exactamente con el total.</li>
                <li>El encargado confirmará o rechazará tu pago.</li>
                <li>Si se rechaza, se liberan cupos reservados.</li>
              </ul>
            </template>
          </Card>
        </div>
      </div>
    </div>
    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import ExportarConfirmacionPDF from '~/components/compras/ExportarConfirmacionPDF.vue'
import SubirComprobante from '~/components/paquetes/SubirComprobante.vue'

definePageMeta({
  middleware: 'turista',
  layout: 'turista'
})

const toast = useToast()
const route = useRoute()
const config = useRuntimeConfig()
const assetsBase = String(config.public.apiBase || '').replace(/\/api\/v1\/?$/, '')

const { obtenerDetalleCompra, cancelarCompra } = useCompra()
const { getPaquete } = usePaquetesTuristicos()

const compraId = Number(route.params.id)

const compra = ref<any>(null)
const paquete = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const COUNTDOWN_SECONDS = 5 * 60
const countdownStorageKey = `compra_timer_${compraId}`
const countdownDeadline = ref<number | null>(null)
const remainingSeconds = ref<number>(COUNTDOWN_SECONDS)
const countdownTimer = ref<ReturnType<typeof setInterval> | null>(null)
const autoCancelled = ref(false)
const autoCancelling = ref(false)

const statusLabel = computed(() => {
  const map: Record<string, string> = {
    pendiente_confirmacion: 'Pendiente de confirmación',
    confirmada: 'Confirmada',
    rechazada: 'Rechazada',
    cancelada: 'Cancelada',
    completada: 'Completada'
  }
  return map[compra.value?.status] || compra.value?.status || 'N/D'
})

const statusSeverity = computed(() => {
  const map: Record<string, any> = {
    pendiente_confirmacion: 'warning',
    confirmada: 'success',
    rechazada: 'danger',
    cancelada: 'secondary',
    completada: 'info'
  }
  return map[compra.value?.status] || 'secondary'
})

const isConfirmada = computed(() => compra.value?.status === 'confirmada')
const codigoConfirmacion = computed(() => {
  const raw = compra.value?.codigo_confirmacion
  if (raw) return String(raw)
  if (compra.value?.id) return `#${compra.value.id}`
  return 'N/D'
})

const pageTitle = computed(() => {
  if (isConfirmada.value) return 'Compra confirmada'
  if (compra.value?.status && compra.value?.status !== 'pendiente_confirmacion') return 'Detalle de compra'
  return 'Pagar compra'
})

const pageSubtitle = computed(() => {
  if (isConfirmada.value) return 'Tu pago fue confirmado. Aqui tienes la guia de tu salida.'
  if (compra.value?.status && compra.value?.status !== 'pendiente_confirmacion') {
    return 'Revisa el estado y los detalles de tu compra.'
  }
  return 'Registra tu pago para que el encargado pueda confirmarlo.'
})

const tipoCompraLabel = computed(() => {
  const map: Record<string, string> = {
    compartido: 'Compartido',
    privado: 'Privado'
  }
  return map[compra.value?.tipo_compra] || compra.value?.tipo_compra || 'N/D'
})

const metodoPagoLabel = computed(() => {
  const map: Record<string, string> = {
    efectivo: 'Efectivo',
    qr: 'QR',
    transferencia: 'Transferencia'
  }
  return map[compra.value?.ultimo_pago?.metodo_pago] || compra.value?.ultimo_pago?.metodo_pago || 'N/D'
})

const pagoStatusLabel = computed(() => {
  const map: Record<string, string> = {
    pendiente: 'Pendiente',
    confirmado: 'Confirmado',
    rechazado: 'Rechazado'
  }
  return map[compra.value?.ultimo_pago?.estado] || compra.value?.ultimo_pago?.estado || 'N/D'
})

const pagoStatusSeverity = computed(() => {
  const map: Record<string, any> = {
    pendiente: 'warning',
    confirmado: 'success',
    rechazado: 'danger'
  }
  return map[compra.value?.ultimo_pago?.estado] || 'secondary'
})

const isPagoPendiente = computed(() => compra.value?.ultimo_pago?.estado === 'pendiente')

const shouldCountdown = computed(() => {
  if (!compra.value) return false
  if (compra.value.status !== 'pendiente_confirmacion') return false
  if (isPagoPendiente.value) return false
  if (autoCancelled.value) return false
  return true
})

const countdownLabel = computed(() => {
  const total = Math.max(0, remainingSeconds.value)
  const minutes = Math.floor(total / 60)
  const seconds = total % 60
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})

const countdownPercent = computed(() => {
  if (!COUNTDOWN_SECONDS) return 0
  return Math.max(0, Math.min(100, (remainingSeconds.value / COUNTDOWN_SECONDS) * 100))
})

const duracionLabel = computed(() => {
  const value = Number(paquete.value?.duracion_dias || 0)
  if (!value) return 'N/D'
  return value === 1 ? '1 dia' : `${value} dias`
})

const horarioLabel = computed(() => {
  const raw = paquete.value?.horario
  if (!raw) return 'N/D'
  const map: Record<string, string> = {
    manana: 'Mañana',
    mañana: 'Mañana',
    tarde: 'Tarde',
    todo_dia: 'Todo el dia'
  }
  return map[String(raw)] || String(raw)
})

const dificultadLabel = computed(() => {
  const map: Record<string, string> = {
    facil: 'Facil',
    medio: 'Medio',
    dificil: 'Dificil',
    extremo: 'Extremo'
  }
  return map[paquete.value?.nivel_dificultad] || paquete.value?.nivel_dificultad || 'N/D'
})

const formatHora = (raw?: any) => {
  if (!raw) return 'N/D'
  const value = String(raw).trim()
  const match = value.match(/(\d{1,2}):(\d{2})/)
  if (match) return `${match[1].padStart(2, '0')}:${match[2]}`
  const parsed = new Date(value)
  if (!Number.isNaN(parsed.getTime())) {
    const hh = String(parsed.getHours()).padStart(2, '0')
    const mm = String(parsed.getMinutes()).padStart(2, '0')
    return `${hh}:${mm}`
  }
  return value
}

const horaSalidaLabel = computed(() => formatHora(paquete.value?.hora_salida))
const horaSalidaConfirmada = computed(() => {
  const encuentro = formatHora(compra.value?.salida?.hora_encuentro)
  if (encuentro !== 'N/D') return encuentro
  return horaSalidaLabel.value
})
const fechaSalidaLabel = computed(() => {
  const value = compra.value?.salida?.fecha_salida || compra.value?.fecha_seleccionada
  return value ? formatFecha(value) : 'Por definir'
})
const fechaHoraSalidaLabel = computed(() => {
  const fecha = fechaSalidaLabel.value
  if (!fecha || fecha === 'Por definir') return 'Por definir'
  const hora = horaSalidaConfirmada.value
  if (!hora || hora === 'N/D') return fecha
  return `${fecha} ${hora}`
})
const fechaConfirmacionLabel = computed(() => {
  const value = compra.value?.fecha_confirmacion || compra.value?.ultimo_pago?.fecha_confirmacion
  if (!value) return 'Por definir'
  const formatted = formatFechaHora(value)
  return formatted || 'Por definir'
})
const puntoEncuentroLabel = computed(() => compra.value?.salida?.punto_encuentro || paquete.value?.punto_encuentro || '')
const instruccionesTuristas = computed(() => compra.value?.salida?.instrucciones_turistas || '')

const agenciaNombre = computed(() => paquete.value?.agencia?.nombre_comercial || 'N/D')
const agenciaUbicacion = computed(() => paquete.value?.agencia?.departamento?.nombre || paquete.value?.agencia?.direccion || 'N/D')
const agenciaDireccionLabel = computed(() => {
  const dept = paquete.value?.agencia?.departamento?.nombre
  const dir = paquete.value?.agencia?.direccion
  if (dept && dir) return `${dept} - ${dir}`
  return dept || dir || 'N/D'
})
const agenciaTelefono = computed(() => paquete.value?.agencia?.telefono || '')
const agenciaEmail = computed(() => paquete.value?.agencia?.email || '')
const agenciaWhatsapp = computed(() => {
  return (
    paquete.value?.agencia?.whatsapp ||
    paquete.value?.agencia?.telefono_whatsapp ||
    paquete.value?.agencia?.telefonoWhatsapp ||
    paquete.value?.agencia?.telefono ||
    ''
  )
})
const agenciaTelefonoLink = computed(() => {
  const digits = normalizePhone(agenciaTelefono.value)
  return digits ? `tel:${digits}` : ''
})
const agenciaWhatsappLink = computed(() => {
  const digits = normalizePhone(agenciaWhatsapp.value)
  return digits ? `https://wa.me/${digits}` : ''
})

const paqueteFotoUrl = computed(() => {
  const fotos = (paquete.value?.fotos || []).slice()
  fotos.sort((a: any, b: any) => {
    if (!!a.es_principal !== !!b.es_principal) return a.es_principal ? -1 : 1
    return (a.orden || 0) - (b.orden || 0)
  })
  const principal = fotos[0]
  return principal?.foto ? resolveAssetUrl(principal.foto) : ''
})

const fotosOrdenadas = computed(() => {
  const fotos = (paquete.value?.fotos || []).slice()
  fotos.sort((a: any, b: any) => {
    if (!!a.es_principal !== !!b.es_principal) return a.es_principal ? -1 : 1
    return (a.orden || 0) - (b.orden || 0)
  })
  return fotos
})

const isMultiDay = computed(() => Number(paquete.value?.duracion_dias || 0) > 1)

const itinerarioOrdenado = computed(() => {
  const items = (paquete.value?.itinerario || []).slice()
  items.sort((a: any, b: any) => (a.dia_numero || 0) - (b.dia_numero || 0))
  return items
})

const atraccionesOrdenadas = computed(() => {
  const items = (paquete.value?.atracciones || []).slice()
  items.sort((a: any, b: any) => {
    const diaA = Number(a?.dia_numero || 0)
    const diaB = Number(b?.dia_numero || 0)
    if (diaA !== diaB) return diaA - diaB
    return Number(a?.orden_visita || 0) - Number(b?.orden_visita || 0)
  })
  return items
})

const atraccionesPorDia = computed(() => {
  const map: Record<number, any[]> = {}
  atraccionesOrdenadas.value.forEach((item: any) => {
    const dia = Math.max(1, Number(item?.dia_numero || 1))
    map[dia] ||= []
    map[dia].push(item)
  })
  return map
})

const diasAtracciones = computed(() => Object.keys(atraccionesPorDia.value).map((k) => Number(k)).sort((a, b) => a - b))

const getOrdenVisitaLabel = (items: any[] | undefined, idx: number) => {
  if (!items || items.length === 0) return idx + 1
  const ordenes = items.map((item) => Number(item?.orden_visita || 0))
  const allValid = ordenes.every((orden) => Number.isFinite(orden) && orden > 0)
  if (!allValid) return idx + 1
  if (new Set(ordenes).size !== ordenes.length) return idx + 1
  return ordenes[idx] || idx + 1
}

const requisitosItems = computed(() => {
  const raw = paquete.value?.requisitos ?? paquete.value?.requisitos_texto ?? paquete.value?.requisitos_lista
  if (Array.isArray(raw)) {
    return raw.map((item: any) => String(item || '').trim()).filter(Boolean)
  }
  if (typeof raw === 'string') {
    const lines = raw.split('\n').map((item) => item.trim()).filter(Boolean)
    if (lines.length > 1) return lines
    const commas = raw.split(',').map((item) => item.trim()).filter(Boolean)
    if (commas.length > 1) return commas
    const trimmed = raw.trim()
    return trimmed ? [trimmed] : []
  }
  return []
})

const permitePrivado = computed(() => !!paquete.value?.permite_privado)
const edadMinimaPago = computed(() => Number(paquete.value?.politicas?.edad_minima_pago || 0))
const recargoPrivado = computed(() => Number(paquete.value?.politicas?.recargo_privado_porcentaje || 0))
const politicaCancelacion = computed(() => paquete.value?.politicas?.politica_cancelacion || '')
const tienePoliticas = computed(() => {
  if (edadMinimaPago.value > 0) return true
  if (permitePrivado.value && recargoPrivado.value > 0) return true
  return false
})

const contactoEmergenciaLabel = computed(() => {
  const raw =
    paquete.value?.contacto_emergencia ||
    paquete.value?.contactoEmergencia ||
    paquete.value?.emergencia_contacto ||
    paquete.value?.emergenciaContacto ||
    ''
  if (raw) return String(raw).trim()
  if (agenciaTelefono.value) return `Agencia: ${agenciaTelefono.value}`
  return ''
})

const recomendacionesClimaticasLabel = computed(() => {
  const raw =
    paquete.value?.recomendaciones_climaticas ||
    paquete.value?.recomendacionesClimaticas ||
    paquete.value?.recomendaciones_clima ||
    paquete.value?.recomendacionesClima ||
    ''
  return raw ? String(raw).trim() : ''
})

const formatMoney = (value: any) => {
  const n = Number(value || 0)
  return n.toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const displayValue = (value?: string) => {
  if (!value) return 'Por definir'
  const raw = String(value).trim()
  if (!raw || raw === 'N/D') return 'Por definir'
  return raw
}

const formatCount = (value?: any) => {
  if (value === null || value === undefined) return 'Por definir'
  const num = Number(value)
  if (!Number.isFinite(num)) return 'Por definir'
  return String(num)
}

const adultosLabel = computed(() => formatCount(compra.value?.cantidad_adultos))
const ninosPaganLabel = computed(() => formatCount(compra.value?.cantidad_ninos_pagan))
const ninosGratisLabel = computed(() => formatCount(compra.value?.cantidad_ninos_gratis))
const totalParticipantesLabel = computed(() => formatCount(compra.value?.total_participantes))

const formatFechaHora = (value?: any) => {
  if (!value) return ''
  const raw = String(value)
  const parsed = new Date(raw)
  if (!Number.isNaN(parsed.getTime())) {
    const d = String(parsed.getDate()).padStart(2, '0')
    const m = String(parsed.getMonth() + 1).padStart(2, '0')
    const y = parsed.getFullYear()
    const hh = String(parsed.getHours()).padStart(2, '0')
    const mm = String(parsed.getMinutes()).padStart(2, '0')
    return `${d}/${m}/${y} ${hh}:${mm}`
  }
  return formatFecha(raw)
}

const formatFecha = (value?: any) => {
  if (!value) return ''
  const raw = String(value)
  const datePart = raw.split('T').shift() ?? raw
  const clean = datePart.split(' ').shift() ?? datePart
  const match = clean.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`
  return clean || raw
}

const pdfRows = computed(() => {
  if (!compra.value) return []
  return [
    ['Codigo de confirmacion', codigoConfirmacion.value],
    ['Estado', statusLabel.value],
    ['Paquete', compra.value?.paquete?.nombre || 'Paquete'],
    ['Tipo de compra', tipoCompraLabel.value],
    ['Fecha de compra', displayValue(formatFechaHora(compra.value.fecha_compra))],
    ['Fecha de confirmacion', fechaConfirmacionLabel.value],
    ['Fecha y hora de salida', fechaHoraSalidaLabel.value],
    ['Punto de encuentro', displayValue(puntoEncuentroLabel.value)],
    ['Duracion', displayValue(duracionLabel.value)],
    ['Horario', displayValue(horarioLabel.value)],
    ['Dificultad', displayValue(dificultadLabel.value)],
    ['Politicas de cancelacion', politicaCancelacion.value || 'Por definir'],
    ['Contacto de emergencia', displayValue(contactoEmergenciaLabel.value)],
    ['Recomendaciones climaticas', displayValue(recomendacionesClimaticasLabel.value)],
    ['Participantes (total)', totalParticipantesLabel.value],
    ['Adultos', adultosLabel.value],
    ['Ninos pagan', ninosPaganLabel.value],
    ['Ninos gratis', ninosGratisLabel.value],
    ['Total pagado', `Bs. ${formatMoney(compra.value.precio_total)}`],
    ['Nombre comercial', displayValue(agenciaNombre.value)],
    ['Ubicacion / Direccion', displayValue(agenciaDireccionLabel.value)],
    ['Telefono', displayValue(agenciaTelefono.value)],
    ['Email', displayValue(agenciaEmail.value)],
    ['WhatsApp', displayValue(agenciaWhatsapp.value)]
  ]
})

const pdfFileName = computed(() => {
  const id = compra.value?.id
  return id ? `confirmacion_compra_${id}.pdf` : 'confirmacion_compra.pdf'
})

const resolveAssetUrl = (path?: string) => {
  if (!path) return ''
  let normalized = String(path).replace(/\\/g, '/')
  if (/^https?:\/\//i.test(normalized)) return normalized
  const uploadsIdx = normalized.indexOf('uploads/')
  if (uploadsIdx > -1) normalized = normalized.slice(uploadsIdx)
  normalized = normalized.replace(/^\.\//, '')
  const clean = normalized.startsWith('/') ? normalized.slice(1) : normalized
  return `${assetsBase}/${clean}`
}

const normalizePhone = (value?: string) => {
  if (!value) return ''
  const digits = String(value).replace(/\D/g, '')
  return digits
}

const stopCountdown = () => {
  if (countdownTimer.value) {
    clearInterval(countdownTimer.value)
    countdownTimer.value = null
  }
}

const getStoredDeadline = () => {
  if (typeof window === 'undefined') return null
  const raw = window.localStorage.getItem(countdownStorageKey)
  if (!raw) return null
  const parsed = Number(raw)
  if (!Number.isFinite(parsed)) return null
  return parsed
}

const setStoredDeadline = (deadline: number) => {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(countdownStorageKey, String(deadline))
}

const clearStoredDeadline = () => {
  if (typeof window === 'undefined') return
  window.localStorage.removeItem(countdownStorageKey)
}

const ensureCountdownDeadline = () => {
  const stored = getStoredDeadline()
  if (stored && stored > Date.now()) {
    countdownDeadline.value = stored
    return
  }
  if (stored && stored <= Date.now()) {
    countdownDeadline.value = stored
    remainingSeconds.value = 0
    stopCountdown()
    handleAutoCancel()
    return
  }
  const deadline = Date.now() + COUNTDOWN_SECONDS * 1000
  countdownDeadline.value = deadline
  setStoredDeadline(deadline)
}

const updateCountdown = () => {
  if (!countdownDeadline.value) return
  const remaining = Math.max(0, Math.ceil((countdownDeadline.value - Date.now()) / 1000))
  remainingSeconds.value = remaining
  if (remaining <= 0) {
    stopCountdown()
    handleAutoCancel()
  }
}

const startCountdown = () => {
  if (countdownTimer.value) return
  updateCountdown()
  countdownTimer.value = setInterval(updateCountdown, 1000)
}

const resetCountdownState = () => {
  remainingSeconds.value = COUNTDOWN_SECONDS
  countdownDeadline.value = null
  clearStoredDeadline()
}

const handleAutoCancel = async () => {
  if (autoCancelled.value || autoCancelling.value) return
  if (!compra.value || compra.value.status !== 'pendiente_confirmacion') return
  if (isPagoPendiente.value) return

  autoCancelled.value = true
  autoCancelling.value = true
  try {
    const response: any = await cancelarCompra(compraId, 'Compra expirada por tiempo en pantalla de pago')
    if (!response?.success) {
      autoCancelled.value = false
      toast.add({
        severity: 'warn',
        summary: 'No se pudo cancelar',
        detail: response?.error?.message || 'No se pudo cancelar la compra automaticamente',
        life: 4000
      })
      return
    }
    toast.add({
      severity: 'info',
      summary: 'Tiempo agotado',
      detail: 'La compra fue cancelada y los cupos se liberaron',
      life: 4000
    })
  } catch (err: any) {
    autoCancelled.value = false
    toast.add({
      severity: 'warn',
      summary: 'No se pudo cancelar',
      detail: err?.data?.error?.message || err?.message || 'No se pudo cancelar la compra automaticamente',
      life: 4000
    })
  } finally {
    autoCancelling.value = false
    await loadAll()
  }
}

const loadAll = async () => {
  loading.value = true
  error.value = null
  try {
    const response: any = await obtenerDetalleCompra(compraId)
    if (!response.success) {
      error.value = response?.error?.message || 'No se pudo cargar la compra'
      return
    }
    compra.value = response.data
    const paqueteFromCompra = compra.value?.paquete ?? null
    paquete.value = paqueteFromCompra

    const needsPaqueteFetch =
      paqueteFromCompra &&
      !Object.prototype.hasOwnProperty.call(paqueteFromCompra, 'agencia') &&
      !Object.prototype.hasOwnProperty.call(paqueteFromCompra, 'incluye') &&
      !Object.prototype.hasOwnProperty.call(paqueteFromCompra, 'descripcion')

    const paqueteId = Number(paqueteFromCompra?.id)
    if (needsPaqueteFetch && paqueteId) {
      const paqueteResp: any = await getPaquete(paqueteId)
      if (paqueteResp.success) {
        paquete.value = paqueteResp.data
      }
    }
  } catch (err: any) {
    error.value = err?.data?.error?.message || err?.message || 'No se pudo cargar la compra'
    toast.add({ severity: 'error', summary: 'Error', detail: error.value, life: 3000 })
  } finally {
    loading.value = false
  }
}

watch(
  () => [compra.value?.status, isPagoPendiente.value, autoCancelled.value],
  () => {
    if (shouldCountdown.value) {
      ensureCountdownDeadline()
      startCountdown()
      return
    }
    stopCountdown()
    if (compra.value && (compra.value.status !== 'pendiente_confirmacion' || isPagoPendiente.value)) {
      resetCountdownState()
    }
  }
)

onMounted(() => {
  loadAll()
})

onBeforeUnmount(() => {
  stopCountdown()
})
</script>
