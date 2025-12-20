<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div class="flex items-start justify-between mb-6 flex-wrap gap-3">
      <div>
        <p class="text-sm text-gray-500">Panel de agencia · Paquetes turísticos</p>
        <h1 class="text-3xl font-bold text-gray-900">{{ paquete?.nombre || 'Paquete' }}</h1>

        <div v-if="paquete" class="flex flex-wrap gap-2 mt-2">
          <Tag :value="getStatusLabel(paquete.status)" :severity="getStatusSeverity(paquete.status)" />
          <Tag :value="getFrecuenciaLabel(paquete.frecuencia)" severity="info" />
          <Tag
            v-if="paquete.frecuencia === 'salida_unica' && paquete.fecha_salida_fija"
            :value="formatFecha(paquete.fecha_salida_fija)"
            severity="secondary"
            icon="pi pi-calendar"
          />
          <Tag :value="paquete.visible_publico ? 'Público' : 'Solo yo'" :severity="paquete.visible_publico ? 'success' : 'warning'" />
        </div>
      </div>

      <div class="flex gap-2">
        <Button label="Volver" icon="pi pi-arrow-left" outlined @click="navigateTo('/agencia/paquetes')" />
        <Button
          label="Eliminar"
          icon="pi pi-trash"
          severity="danger"
          outlined
          :disabled="!paquete || deleting"
          @click="showDeleteDialog = true"
        />
      </div>
    </div>

    <Card class="mb-4">
      <template #content>
        <div class="flex flex-wrap gap-2">
          <Button
            v-for="tab in tabs"
            :key="tab.value"
            :label="tab.label"
            :icon="tab.icon"
            :severity="activeTab === tab.value ? 'primary' : 'secondary'"
            outlined
            @click="activeTab = tab.value"
          />
        </div>
      </template>
    </Card>

    <div v-if="loading" class="space-y-4">
      <Skeleton height="420px" />
    </div>

    <div v-else-if="!paquete" class="space-y-3">
      <Message severity="warn" :closable="false">No se pudo cargar el paquete.</Message>
      <Button label="Reintentar" icon="pi pi-refresh" outlined @click="loadAll" />
    </div>

    <template v-else>
      <Card v-if="activeTab === 'general'" class="surface-card">
        <template #title>Información del paquete</template>
        <template #content>
          <form class="space-y-6" @submit.prevent="saveGeneral">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Nombre</label>
                <InputText v-model="generalForm.nombre" class="w-full" />
              </div>

              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <Textarea v-model="generalForm.descripcion" class="w-full" rows="4" autoResize />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Frecuencia</label>
                <InputText :modelValue="getFrecuenciaLabel(paquete.frecuencia)" class="w-full" disabled />
                <small class="text-gray-500">La frecuencia no se puede cambiar.</small>
              </div>

              <div v-if="paquete.frecuencia === 'salida_unica'">
                <label class="block text-sm font-medium text-gray-700 mb-1">Fecha de salida fija</label>
                <DatePicker v-model="generalForm.fecha_salida_fija" class="w-full" dateFormat="yy-mm-dd" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Duración (días)</label>
                <InputNumber v-model="generalForm.duracion_dias" class="w-full" :min="1" :max="30" />
                <small class="text-gray-500" v-if="generalIsMultiDay">Noches: {{ generalDuracionNoches }}</small>
                <small class="text-gray-500" v-else>Paquete de un día.</small>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Días previos de compra</label>
                <InputNumber v-model="generalForm.dias_previos_compra" class="w-full" :min="1" :max="60" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nivel de dificultad</label>
                <Select
                  v-model="generalForm.nivel_dificultad"
                  :options="dificultadOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                  showClear
                  placeholder="(Opcional)"
                />
              </div>
            </div>

            <div v-if="!generalIsMultiDay" class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Horario</label>
                <Select
                  v-model="generalForm.horario"
                  :options="horarioOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                  showClear
                  placeholder="(Opcional)"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Hora de salida</label>
                <InputMask v-model="generalForm.hora_salida" mask="99:99" placeholder="08:30" class="w-full" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Duración (horas)</label>
                <InputNumber v-model="generalForm.duracion_horas_num" class="w-full" :min="1" :max="24" />
              </div>
            </div>

            <Divider />

            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Cupo mínimo</label>
                <InputNumber v-model="generalForm.cupo_minimo" class="w-full" :min="1" :max="500" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Cupo máximo</label>
                <InputNumber v-model="generalForm.cupo_maximo" class="w-full" :min="1" :max="500" />
              </div>
              <div class="md:col-span-2 flex items-center gap-3 pt-6">
                <Checkbox v-model="generalForm.permite_privado" binary inputId="permitePrivadoEdit" />
                <label for="permitePrivadoEdit" class="text-sm text-gray-700">Permite contratar como privado</label>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Precio base (nacionales)</label>
                <InputNumber v-model="generalForm.precio_base_nacionales" class="w-full" :min="0" :max="999999" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Adicional (extranjeros)</label>
                <InputNumber v-model="generalForm.precio_adicional_extranjeros" class="w-full" :min="0" :max="999999" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Estado</label>
                <Select
                  v-model="generalForm.status"
                  :options="statusOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full"
                />
              </div>
              <div class="flex items-center gap-3 pt-6">
                <InputSwitch v-model="generalForm.visible_publico" />
                <div>
                  <p class="text-sm font-medium text-gray-700">Visible al público</p>
                  <p class="text-xs text-gray-500">Si está desactivado, solo tú lo verás.</p>
                </div>
              </div>
            </div>

            <Divider />

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <StringListEditor v-model="generalForm.incluye" label="Incluye" placeholder="Ej: Transporte" emptyText="Sin elementos" />
              <StringListEditor v-model="generalForm.no_incluye" label="No incluye" placeholder="Ej: Almuerzo" emptyText="Sin elementos" />
              <StringListEditor v-model="generalForm.que_llevar" label="Qué llevar" placeholder="Ej: Bloqueador solar" emptyText="Sin elementos" />
            </div>

            <div class="flex justify-end gap-2 pt-2">
              <Button label="Guardar cambios" icon="pi pi-save" type="submit" :loading="savingGeneral" />
            </div>
          </form>
        </template>
      </Card>

      <Card v-else-if="activeTab === 'fotos'" class="surface-card">
        <template #title>Fotos</template>
        <template #content>
          <div class="space-y-4">
            <div class="flex flex-wrap items-center gap-3">
              <input
                ref="fileInput"
                type="file"
                accept="image/jpeg,image/png,image/webp"
                multiple
                class="hidden"
                :disabled="maxFotosReached || uploading"
                @change="handleFileSelect"
              />
              <Button
                label="Seleccionar Fotos"
                icon="pi pi-upload"
                outlined
                :disabled="maxFotosReached || uploading"
                @click="fileInput?.click()"
              />
              <Button
                :label="uploadButtonLabel"
                icon="pi pi-cloud-upload"
                :loading="uploading"
                :disabled="selectedFotos.length === 0 || uploading"
                @click="uploadSelectedFotos"
              />
              <span class="text-sm text-gray-600">{{ totalFotos }} / 6 en total ({{ selectedFotos.length }} nuevas)</span>
              <small class="text-gray-500">Max 5MB c/u. Max 6 fotos.</small>
            </div>

            <div class="flex items-center gap-2">
              <Checkbox
                v-model="setPrimeraComoPrincipal"
                binary
                inputId="setPrimeraComoPrincipalPaquete"
                :disabled="uploading"
              />
              <label for="setPrimeraComoPrincipalPaquete" class="text-sm text-gray-700">
                Marcar la primera como principal
              </label>
            </div>

            <div v-if="selectedFotos.length > 0">
              <p class="text-sm text-gray-600 mb-3">
                <i class="pi pi-info-circle"></i>
                Usa las flechas para reordenar.
              </p>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div
                  v-for="(foto, index) in selectedFotos"
                  :key="index"
                  class="relative group border-2 rounded-lg hover:border-blue-300 transition-all overflow-hidden"
                >
                  <img :src="foto.preview" class="w-full h-32 object-cover" :alt="foto.file.name" />

                  <div class="absolute top-2 left-2 flex flex-col gap-1">
                    <span class="px-2 py-1 bg-gray-900 bg-opacity-70 text-white rounded text-xs font-mono">
                      {{ index + 1 }}
                    </span>
                    <span
                      v-if="setPrimeraComoPrincipal && index === 0"
                      class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold"
                    >
                      PRINCIPAL
                    </span>
                  </div>

                  <div class="absolute top-2 right-2 flex flex-col gap-1">
                    <Button
                      icon="pi pi-arrow-up"
                      rounded
                      size="small"
                      text
                      class="bg-white"
                      :disabled="index === 0 || uploading"
                      @click="moveSelectedFotoUp(index)"
                    />
                    <Button
                      icon="pi pi-arrow-down"
                      rounded
                      size="small"
                      text
                      class="bg-white"
                      :disabled="index === selectedFotos.length - 1 || uploading"
                      @click="moveSelectedFotoDown(index)"
                    />
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      rounded
                      size="small"
                      :disabled="uploading"
                      class="bg-white"
                      @click="removeSelectedFoto(index)"
                    />
                  </div>

                  <div class="bg-gray-50 p-2">
                    <p class="text-xs text-gray-600 truncate" :title="foto.file.name">{{ foto.file.name }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="fotosOrdenadas.length > 0">
              <h4 class="font-semibold text-gray-900">Fotos actuales</h4>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div
                  v-for="foto in fotosOrdenadas"
                  :key="foto.id"
                  class="relative group border-2 rounded-lg overflow-hidden"
                >
                  <img
                    :src="resolveFotoUrl(foto.foto)"
                    class="w-full h-32 object-cover"
                    :alt="`Foto ${foto.id}`"
                  />
                  <div class="absolute top-2 left-2 flex flex-col gap-1">
                    <span v-if="foto.es_principal" class="px-2 py-1 bg-blue-500 text-white rounded text-xs font-semibold">
                      PRINCIPAL
                    </span>
                  </div>
                  <div class="absolute top-2 right-2">
                    <Button
                      icon="pi pi-trash"
                      severity="danger"
                      rounded
                      size="small"
                      :disabled="uploading"
                      class="bg-white"
                      @click="removeFoto(foto)"
                    />
                  </div>
                  <div class="bg-gray-50 p-2">
                    <p class="text-xs text-gray-600 truncate" :title="foto.foto">
                      Foto #{{ foto.id }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div
              v-if="selectedFotos.length === 0 && fotosOrdenadas.length === 0"
              class="text-center py-8 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300"
            >
              <i class="pi pi-images text-4xl text-gray-400 mb-2"></i>
              <p class="text-gray-600 text-sm">Haga clic en "Seleccionar Fotos" para agregar imágenes</p>
            </div>
          </div>
        </template>
      </Card>

      <Card v-else-if="activeTab === 'itinerario'" class="surface-card">
        <template #title>Itinerario</template>
        <template #content>
          <div v-if="!generalIsMultiDay" class="space-y-2">
            <Message severity="info" :closable="false">El itinerario aplica solo a paquetes de varios días.</Message>
          </div>

          <div v-else class="space-y-4">
            <div class="flex justify-end">
              <Button label="Agregar día" icon="pi pi-plus" @click="openItinerarioCreate" />
            </div>

            <DataTable :value="paquete.itinerario || []" responsiveLayout="scroll">
              <Column field="dia_numero" header="Día" style="width: 90px" />
              <Column field="titulo" header="Título" />
              <Column header="Actividades">
                <template #body="{ data }">
                  <span class="text-sm text-gray-700">{{ (data.actividades || []).length }}</span>
                </template>
              </Column>
              <Column header="Acciones" style="width: 160px">
                <template #body="{ data }">
                  <div class="flex gap-2 justify-end">
                    <Button icon="pi pi-pencil" severity="warning" text rounded @click="openItinerarioEdit(data)" />
                    <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmDeleteItinerario(data)" />
                  </div>
                </template>
              </Column>
            </DataTable>

            <div v-if="!(paquete.itinerario || []).length" class="text-sm text-gray-600">
              Aún no has agregado itinerario.
            </div>
          </div>
        </template>
      </Card>

      <Card v-else-if="activeTab === 'atracciones'" class="surface-card">
        <template #title>Atracciones incluidas</template>
        <template #content>
          <div class="space-y-4">
            <div class="flex justify-end">
              <Button label="Agregar atracción" icon="pi pi-plus" @click="openAddAtraccionDialog" />
            </div>

            <DataTable :value="paquete.atracciones || []" responsiveLayout="scroll">
              <Column v-if="generalIsMultiDay" field="dia_numero" header="Día" style="width: 90px" />
              <Column field="orden_visita" header="Orden" style="width: 110px" />
              <Column header="Atracción">
                <template #body="{ data }">
                  <div class="space-y-0.5">
                    <div class="font-medium text-gray-900">{{ data.atraccion?.nombre || `Atracción #${data.atraccion_id}` }}</div>
                    <div class="text-xs text-gray-500">{{ data.atraccion?.provincia?.nombre || '' }}</div>
                  </div>
                </template>
              </Column>
              <Column header="Duración (h)" style="width: 140px">
                <template #body="{ data }">
                  <span class="text-sm text-gray-700">{{ data.duracion_estimada_horas || '—' }}</span>
                </template>
              </Column>
              <Column header="Acciones" style="width: 160px">
                <template #body="{ data }">
                  <div class="flex gap-2 justify-end">
                    <Button icon="pi pi-pencil" severity="warning" text rounded @click="openEditAtraccion(data)" />
                    <Button icon="pi pi-trash" severity="danger" text rounded @click="confirmRemoveAtraccion(data)" />
                  </div>
                </template>
              </Column>
            </DataTable>

            <div v-if="!(paquete.atracciones || []).length" class="text-sm text-gray-600">
              Aún no has agregado atracciones a este paquete.
            </div>
          </div>
        </template>
      </Card>

      <Card v-else-if="activeTab === 'salidas'" class="surface-card">
        <template #title>Salidas habilitadas</template>
        <template #content>
          <div class="space-y-3">
            <Message v-if="paquete.frecuencia === 'salida_diaria'" severity="info" :closable="false">
              En paquetes de salida diaria, las salidas se crean cuando un turista compra para una fecha.
            </Message>

            <DataTable :value="paquete.salidas || []" responsiveLayout="scroll">
              <Column header="Fecha" style="width: 130px">
                <template #body="{ data }">
                  <span class="text-sm text-gray-700">{{ formatFecha(data.fecha_salida) }}</span>
                </template>
              </Column>
              <Column field="tipo_salida" header="Tipo" style="width: 130px" />
              <Column header="Cupos">
                <template #body="{ data }">
                  <span class="text-sm text-gray-700">
                    {{ data.cupos_confirmados || 0 }}/{{ data.cupo_maximo }} confirmados · {{ data.cupos_reservados || 0 }} reservados
                  </span>
                </template>
              </Column>
              <Column field="estado" header="Estado" style="width: 140px" />
              <Column header="Acciones" style="width: 120px">
                <template #body="{ data }">
                  <div class="flex justify-end">
                    <Button icon="pi pi-pencil" severity="warning" text rounded @click="openEditSalida(data)" />
                  </div>
                </template>
              </Column>
            </DataTable>

            <div v-if="!(paquete.salidas || []).length" class="text-sm text-gray-600">
              Aún no existen salidas para este paquete.
            </div>
          </div>
        </template>
      </Card>
    </template>

    <Dialog v-model:visible="showDeleteDialog" header="Eliminar paquete" modal :style="{ width: '420px' }">
      <div class="space-y-3">
        <p class="text-sm text-gray-700">
          Se eliminará el paquete (status: <strong>eliminado</strong>) y dejará de ser visible al público.
        </p>
        <Message severity="warn" :closable="false">Esta acción no borra físicamente registros en la base de datos.</Message>
      </div>
      <template #footer>
        <Button label="Cancelar" severity="secondary" outlined @click="showDeleteDialog = false" />
        <Button label="Eliminar" icon="pi pi-trash" severity="danger" :loading="deleting" @click="deleteThisPaquete" />
      </template>
    </Dialog>

    <Dialog
      v-model:visible="showItinerarioDialog"
      :header="itinerarioEditing ? 'Editar día' : 'Agregar día'"
      modal
      :style="{ width: '720px' }"
    >
      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Día</label>
            <InputNumber v-model="itForm.dia_numero" class="w-full" :min="1" :max="generalForm.duracion_dias" />
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Título</label>
            <InputText v-model="itForm.titulo" class="w-full" />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
          <Textarea v-model="itForm.descripcion" class="w-full" rows="3" autoResize />
        </div>

        <StringListEditor v-model="itForm.actividades" label="Actividades" placeholder="Ej: Visita a viñedos" />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Hospedaje (opcional)</label>
          <Textarea v-model="itForm.hospedaje_info" class="w-full" rows="2" autoResize />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" severity="secondary" outlined @click="showItinerarioDialog = false" />
        <Button
          :label="itinerarioEditing ? 'Guardar' : 'Agregar'"
          icon="pi pi-check"
          :loading="savingItinerario"
          @click="saveItinerario"
        />
      </template>
    </Dialog>

    <Dialog v-model:visible="showDeleteItinerarioDialog" header="Eliminar día" modal :style="{ width: '420px' }">
      <p class="text-sm text-gray-700">¿Seguro que deseas eliminar este elemento del itinerario?</p>
      <template #footer>
        <Button label="Cancelar" severity="secondary" outlined @click="showDeleteItinerarioDialog = false" />
        <Button
          label="Eliminar"
          icon="pi pi-trash"
          severity="danger"
          :loading="deletingItinerario"
          @click="deleteItinerario"
        />
      </template>
    </Dialog>

    <Dialog v-model:visible="showAddAtraccionDialog" header="Agregar atracción" modal :style="{ width: '920px' }">
      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Buscar</label>
            <IconField>
              <InputIcon class="pi pi-search" />
              <InputText v-model="atrSearch" class="w-full" placeholder="Buscar atracciones..." @keyup.enter="searchAtracciones" />
            </IconField>
          </div>
          <div class="flex gap-2 justify-end">
            <Button label="Buscar" icon="pi pi-search" outlined :loading="atrLoading" @click="searchAtracciones" />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div v-if="generalIsMultiDay">
            <label class="block text-sm font-medium text-gray-700 mb-1">Día</label>
            <InputNumber v-model="atrForm.dia_numero" class="w-full" :min="1" :max="generalForm.duracion_dias" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Orden</label>
            <InputNumber v-model="atrForm.orden_visita" class="w-full" :min="1" :max="200" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Duración (h)</label>
            <InputNumber v-model="atrForm.duracion_estimada_horas" class="w-full" :min="1" :max="24" />
          </div>
          <div class="text-xs text-gray-500 flex items-end">
            Elige una atracción de la lista para agregarla.
          </div>
        </div>

        <DataTable :value="atrResults" responsiveLayout="scroll" :loading="atrLoading">
          <Column field="id" header="#" style="width: 80px" />
          <Column field="nombre" header="Nombre" />
          <Column header="Ubicación">
            <template #body="{ data }">
              <span class="text-sm text-gray-700">{{ data.provincia?.nombre || '' }}</span>
            </template>
          </Column>
          <Column header="Acción" style="width: 140px">
            <template #body="{ data }">
              <Button label="Agregar" icon="pi pi-plus" @click="addAtraccionToPaquete(data)" />
            </template>
          </Column>
        </DataTable>
      </div>
      <template #footer>
        <Button label="Cerrar" severity="secondary" outlined @click="showAddAtraccionDialog = false" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showEditAtraccionDialog" header="Editar atracción" modal :style="{ width: '560px' }">
      <div class="space-y-4">
        <Message v-if="selectedAtraccionItem" severity="info" :closable="false">
          {{ selectedAtraccionItem.atraccion?.nombre || `Atracción #${selectedAtraccionItem.atraccion_id}` }}
        </Message>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div v-if="generalIsMultiDay">
            <label class="block text-sm font-medium text-gray-700 mb-1">Día</label>
            <InputNumber v-model="editAtrForm.dia_numero" class="w-full" :min="1" :max="generalForm.duracion_dias" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Orden</label>
            <InputNumber v-model="editAtrForm.orden_visita" class="w-full" :min="1" :max="200" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Duración (h)</label>
            <InputNumber v-model="editAtrForm.duracion_estimada_horas" class="w-full" :min="1" :max="24" />
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" severity="secondary" outlined @click="showEditAtraccionDialog = false" />
        <Button label="Guardar" icon="pi pi-check" :loading="savingAtraccion" @click="saveAtraccionEdit" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showRemoveAtraccionDialog" header="Quitar atracción" modal :style="{ width: '420px' }">
      <p class="text-sm text-gray-700">¿Seguro que deseas quitar esta atracción del paquete?</p>
      <template #footer>
        <Button label="Cancelar" severity="secondary" outlined @click="showRemoveAtraccionDialog = false" />
        <Button label="Quitar" icon="pi pi-trash" severity="danger" :loading="removingAtraccion" @click="removeAtraccion" />
      </template>
    </Dialog>

    <Dialog v-model:visible="showEditSalidaDialog" header="Editar salida" modal :style="{ width: '820px' }">
      <div class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Fecha</label>
            <InputText :modelValue="formatFecha(selectedSalida?.fecha_salida) || ''" class="w-full" disabled />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tipo</label>
            <InputText :modelValue="selectedSalida?.tipo_salida || ''" class="w-full" disabled />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Estado</label>
            <Select v-model="salidaForm.estado" :options="estadoSalidaOptions" optionLabel="label" optionValue="value" class="w-full" />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Punto de encuentro</label>
            <Textarea v-model="salidaForm.punto_encuentro" class="w-full" rows="2" autoResize />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Hora de encuentro</label>
            <InputMask v-model="salidaForm.hora_encuentro" mask="99:99" placeholder="08:00" class="w-full" />
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Guía (nombre)</label>
            <InputText v-model="salidaForm.guia_nombre" class="w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Guía (teléfono)</label>
            <InputText v-model="salidaForm.guia_telefono" class="w-full" />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Instrucciones para turistas</label>
          <Textarea v-model="salidaForm.instrucciones_turistas" class="w-full" rows="3" autoResize />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Notas internas</label>
          <Textarea v-model="salidaForm.notas_logistica" class="w-full" rows="3" autoResize />
        </div>

        <div v-if="salidaForm.estado === 'cancelada'">
          <label class="block text-sm font-medium text-gray-700 mb-1">Razón de cancelación</label>
          <Textarea v-model="salidaForm.razon_cancelacion" class="w-full" rows="2" autoResize />
        </div>
      </div>
      <template #footer>
        <Button label="Cancelar" severity="secondary" outlined @click="showEditSalidaDialog = false" />
        <Button label="Guardar" icon="pi pi-check" :loading="savingSalida" @click="saveSalida" />
      </template>
    </Dialog>

    <Toast />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import StringListEditor from '~/components/paquetes/StringListEditor.vue'

definePageMeta({
  middleware: 'encargado',
  layout: 'agencia'
})

const toast = useToast()
const route = useRoute()

const { getMiAgencia } = useAgencias()
const {
  getPaquete,
  updatePaquete,
  deletePaquete,
  uploadPaqueteFoto,
  removePaqueteFoto,
  createItinerarioItem,
  updateItinerarioItem,
  deleteItinerarioItem,
  addPaqueteAtraccion,
  updatePaqueteAtraccion,
  removePaqueteAtraccion,
  updatePaqueteSalida
} = usePaquetes()
const { getAtracciones } = useAtracciones()

const agencia = ref<any>(null)
const paquete = ref<any>(null)
const loading = ref(true)

const tabs = [
  { label: 'General', value: 'general', icon: 'pi pi-info-circle' },
  { label: 'Fotos', value: 'fotos', icon: 'pi pi-images' },
  { label: 'Itinerario', value: 'itinerario', icon: 'pi pi-calendar' },
  { label: 'Atracciones', value: 'atracciones', icon: 'pi pi-map' },
  { label: 'Salidas', value: 'salidas', icon: 'pi pi-compass' }
]
const activeTab = ref('general')

const statusOptions = [
  { label: 'Borrador', value: 'borrador' },
  { label: 'Activo', value: 'activo' },
  { label: 'Inactivo', value: 'inactivo' }
]
const dificultadOptions = [
  { label: 'Fácil', value: 'facil' },
  { label: 'Medio', value: 'medio' },
  { label: 'Difícil', value: 'dificil' },
  { label: 'Extremo', value: 'extremo' }
]
const horarioOptions = [
  { label: 'Mañana', value: 'mañana' },
  { label: 'Tarde', value: 'tarde' },
  { label: 'Todo el día', value: 'todo_dia' }
]

const agenciaId = computed(() => Number(agencia.value?.id || 0))
const paqueteId = computed(() => Number(route.params.id || 0))

const parseDate = (value?: string | null) => {
  if (!value) return null
  const raw = String(value).split('T')[0].trim()
  const date = new Date(`${raw}T00:00:00`)
  return Number.isNaN(date.getTime()) ? null : date
}

const formatDate = (date: Date) => {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const formatFecha = (value?: any) => {
  if (!value) return ''

  if (value instanceof Date) {
    const d = String(value.getUTCDate()).padStart(2, '0')
    const m = String(value.getUTCMonth() + 1).padStart(2, '0')
    const y = value.getUTCFullYear()
    return `${d}/${m}/${y}`
  }

  const raw = String(value)
  const datePart = raw.split('T')[0].split(' ')[0]
  const match = datePart.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (match) return `${match[3]}/${match[2]}/${match[1]}`

  const parsed = new Date(raw)
  if (!Number.isNaN(parsed.getTime())) {
    const d = String(parsed.getUTCDate()).padStart(2, '0')
    const m = String(parsed.getUTCMonth() + 1).padStart(2, '0')
    const y = parsed.getUTCFullYear()
    return `${d}/${m}/${y}`
  }

  return datePart || raw
}

const generalForm = ref({
  nombre: '',
  descripcion: '',
  fecha_salida_fija: null as Date | null,
  duracion_dias: 1,
  dias_previos_compra: 1,
  nivel_dificultad: null as string | null,
  horario: null as string | null,
  hora_salida: '',
  duracion_horas_num: 4,
  cupo_minimo: 1,
  cupo_maximo: 10,
  permite_privado: true,
  precio_base_nacionales: 0,
  precio_adicional_extranjeros: 0,
  incluye: [] as string[],
  no_incluye: [] as string[],
  que_llevar: [] as string[],
  status: 'borrador',
  visible_publico: true
})

const generalIsMultiDay = computed(() => Number(generalForm.value.duracion_dias || 1) > 1)
const generalDuracionNoches = computed(() => Math.max(0, Number(generalForm.value.duracion_dias || 1) - 1))

watch(
  () => generalForm.value.duracion_dias,
  (next) => {
    const safe = Number(next || 1)
    if (!Number.isFinite(safe) || safe < 1) generalForm.value.duracion_dias = 1
    if (safe > 1) {
      generalForm.value.horario = null
      generalForm.value.hora_salida = ''
    }
  }
)

const parseIntervalHours = (value?: string | null) => {
  if (!value) return null
  const m = String(value).match(/(\d+)/)
  return m ? Number(m[1]) : null
}

const hydrateGeneralForm = (p: any) => {
  generalForm.value = {
    nombre: p?.nombre || '',
    descripcion: p?.descripcion || '',
    fecha_salida_fija: parseDate(p?.fecha_salida_fija || null),
    duracion_dias: Number(p?.duracion_dias || 1),
    dias_previos_compra: Number(p?.dias_previos_compra || 1),
    nivel_dificultad: p?.nivel_dificultad || null,
    horario: p?.horario || null,
    hora_salida: p?.hora_salida ? String(p.hora_salida).slice(0, 5) : '',
    duracion_horas_num: parseIntervalHours(p?.duracion_horas) || 4,
    cupo_minimo: Number(p?.cupo_minimo || 1),
    cupo_maximo: Number(p?.cupo_maximo || 10),
    permite_privado: !!p?.permite_privado,
    precio_base_nacionales: Number(p?.precio_base_nacionales || 0),
    precio_adicional_extranjeros: Number(p?.precio_adicional_extranjeros || 0),
    incluye: (p?.incluye || []) as string[],
    no_incluye: (p?.no_incluye || []) as string[],
    que_llevar: (p?.que_llevar || []) as string[],
    status: p?.status || 'borrador',
    visible_publico: !!p?.visible_publico
  }
}

const loadAll = async () => {
  loading.value = true
  try {
    const a: any = await getMiAgencia()
    if (a.success) agencia.value = a.data

    const id = agenciaId.value
    const pid = paqueteId.value
    if (!id || !pid) {
      paquete.value = null
      return
    }

    const response: any = await getPaquete(id, pid)
    if (response.success) {
      paquete.value = response.data
      hydrateGeneralForm(response.data)
    } else {
      paquete.value = null
    }
  } catch (error: any) {
    paquete.value = null
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo cargar el paquete',
      life: 3500
    })
  } finally {
    loading.value = false
  }
}

const savingGeneral = ref(false)
const saveGeneral = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  if (!id || !pid) return

  const nombre = generalForm.value.nombre.trim()
  if (!nombre) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'El nombre es obligatorio', life: 3000 })
    return
  }

  const cupoMin = Number(generalForm.value.cupo_minimo || 0)
  const cupoMax = Number(generalForm.value.cupo_maximo || 0)
  if (cupoMin <= 0 || cupoMax <= 0 || cupoMax < cupoMin) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'Revisa cupo mínimo y máximo', life: 3000 })
    return
  }

  if (paquete.value?.frecuencia === 'salida_unica' && !generalForm.value.fecha_salida_fija) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'La fecha fija es obligatoria', life: 3000 })
    return
  }

  const duracionDias = Math.max(1, Number(generalForm.value.duracion_dias || 1))

  const payload: any = {
    nombre,
    descripcion: generalForm.value.descripcion?.trim() || '',
    duracion_dias: duracionDias,
    duracion_noches: duracionDias > 1 ? duracionDias - 1 : null,
    dias_previos_compra: Math.max(1, Number(generalForm.value.dias_previos_compra || 1)),
    nivel_dificultad: generalForm.value.nivel_dificultad || '',
    cupo_minimo: cupoMin,
    cupo_maximo: cupoMax,
    permite_privado: !!generalForm.value.permite_privado,
    precio_base_nacionales: Number(generalForm.value.precio_base_nacionales || 0),
    precio_adicional_extranjeros: Number(generalForm.value.precio_adicional_extranjeros || 0),
    incluye: generalForm.value.incluye || [],
    no_incluye: generalForm.value.no_incluye || [],
    que_llevar: generalForm.value.que_llevar || [],
    status: generalForm.value.status,
    visible_publico: !!generalForm.value.visible_publico
  }

  if (paquete.value?.frecuencia === 'salida_unica') {
    payload.fecha_salida_fija = generalForm.value.fecha_salida_fija ? formatDate(generalForm.value.fecha_salida_fija) : ''
  }

  if (duracionDias <= 1) {
    payload.horario = generalForm.value.horario || ''
    payload.hora_salida = generalForm.value.hora_salida?.trim() || ''
    payload.duracion_horas = generalForm.value.duracion_horas_num ? `${Number(generalForm.value.duracion_horas_num)} hours` : ''
  } else {
    payload.horario = null
    payload.hora_salida = null
    payload.duracion_horas = null
  }

  savingGeneral.value = true
  try {
    const response: any = await updatePaquete(id, pid, payload)
    if (response.success) {
      paquete.value = response.data
      hydrateGeneralForm(response.data)
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Paquete actualizado', life: 2500 })
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo guardar',
      life: 3500
    })
  } finally {
    savingGeneral.value = false
  }
}

const showDeleteDialog = ref(false)
const deleting = ref(false)
const deleteThisPaquete = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  if (!id || !pid) return

  deleting.value = true
  try {
    const response: any = await deletePaquete(id, pid)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Eliminado', detail: 'Paquete eliminado', life: 2500 })
      showDeleteDialog.value = false
      await navigateTo('/agencia/paquetes')
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo eliminar',
      life: 3500
    })
  } finally {
    deleting.value = false
  }
}

const getStatusLabel = (status?: string) => {
  const map: Record<string, string> = { activo: 'Activo', inactivo: 'Inactivo', borrador: 'Borrador', eliminado: 'Eliminado' }
  return map[status || ''] || (status || 'N/D')
}
const getStatusSeverity = (status?: string) => {
  const map: Record<string, string> = { activo: 'success', inactivo: 'warning', borrador: 'info', eliminado: 'danger' }
  return map[status || ''] || 'secondary'
}
const getFrecuenciaLabel = (frecuencia?: string) => {
  const map: Record<string, string> = { salida_diaria: 'Salida diaria', salida_unica: 'Salida única' }
  return map[frecuencia || ''] || (frecuencia || 'N/D')
}

// Fotos
const fileInput = ref<HTMLInputElement | null>(null)
const selectedFotos = ref<Array<{ file: File; preview: string }>>([])
const setPrimeraComoPrincipal = ref(true)
const maxFotos = 6

const fotosOrdenadas = computed(() => {
  const fotos = (paquete.value?.fotos || []) as any[]
  return [...fotos].sort((a, b) => {
    const aPrincipal = a?.es_principal ? 1 : 0
    const bPrincipal = b?.es_principal ? 1 : 0
    if (aPrincipal !== bPrincipal) return bPrincipal - aPrincipal

    const aOrden = Number(a?.orden ?? 0)
    const bOrden = Number(b?.orden ?? 0)
    if (aOrden !== bOrden) return aOrden - bOrden

    const aCreated = a?.created_at ? new Date(a.created_at).getTime() : 0
    const bCreated = b?.created_at ? new Date(b.created_at).getTime() : 0
    if (aCreated !== bCreated) return aCreated - bCreated

    return Number(a?.id ?? 0) - Number(b?.id ?? 0)
  })
})

const totalFotos = computed(() => fotosOrdenadas.value.length + selectedFotos.value.length)
const maxFotosReached = computed(() => totalFotos.value >= maxFotos)

const uploading = ref(false)
const uploadProgress = ref<{ current: number; total: number } | null>(null)

const uploadButtonLabel = computed(() => {
  if (!uploading.value) return 'Subir seleccionadas'
  if (!uploadProgress.value) return 'Subiendo...'
  return `Subiendo ${uploadProgress.value.current}/${uploadProgress.value.total}`
})

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const files = input.files
  if (!files) return

  const remainingSlots = maxFotos - fotosOrdenadas.value.length - selectedFotos.value.length
  if (remainingSlots <= 0) {
    toast.add({ severity: 'warn', summary: 'Límite', detail: 'Máximo 6 fotos', life: 2500 })
    input.value = ''
    return
  }

  const filesToAdd = Array.from(files).slice(0, remainingSlots)

  filesToAdd.forEach((file) => {
    if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
      toast.add({ severity: 'warn', summary: 'Formato no permitido', detail: `${file.name} no es válido`, life: 3000 })
      return
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.add({ severity: 'warn', summary: 'Archivo muy grande', detail: `${file.name} supera 5MB`, life: 3000 })
      return
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      selectedFotos.value.push({
        file,
        preview: e.target?.result as string
      })
    }
    reader.readAsDataURL(file)
  })

  input.value = ''
}

const removeSelectedFoto = (index: number) => {
  selectedFotos.value.splice(index, 1)
}

const moveSelectedFotoUp = (index: number) => {
  if (index === 0) return
  const temp = selectedFotos.value[index]
  selectedFotos.value[index] = selectedFotos.value[index - 1]
  selectedFotos.value[index - 1] = temp
}

const moveSelectedFotoDown = (index: number) => {
  if (index === selectedFotos.value.length - 1) return
  const temp = selectedFotos.value[index]
  selectedFotos.value[index] = selectedFotos.value[index + 1]
  selectedFotos.value[index + 1] = temp
}

const getNextOrdenBase = () => {
  const fotos = (paquete.value?.fotos || []) as any[]
  const maxOrden = fotos.reduce((max, f) => {
    const orden = Number(f?.orden)
    if (!Number.isFinite(orden)) return max
    return Math.max(max, orden)
  }, -1)
  return maxOrden + 1
}

const uploadSelectedFotos = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  if (!id || !pid) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Paquete no cargado', life: 2500 })
    return
  }

  if (selectedFotos.value.length === 0) {
    toast.add({ severity: 'warn', summary: 'Fotos', detail: 'Seleccione al menos una imagen', life: 2500 })
    return
  }

  if (totalFotos.value > maxFotos) {
    toast.add({ severity: 'warn', summary: 'Límite', detail: 'Máximo 6 fotos', life: 2500 })
    return
  }

  uploading.value = true
  uploadProgress.value = { current: 0, total: selectedFotos.value.length }

  const baseOrden = getNextOrdenBase()
  const total = selectedFotos.value.length
  const failed: Array<{ file: File; preview: string }> = []
  let okCount = 0

  for (let i = 0; i < total; i++) {
    const item = selectedFotos.value[i]
    const formData = new FormData()
    formData.append('foto', item.file)
    formData.append('orden', String(baseOrden + i))
    formData.append('es_principal', String(!!setPrimeraComoPrincipal.value && i === 0))

    try {
      const response: any = await uploadPaqueteFoto(id, pid, formData)
      if (response.success) okCount++
      else failed.push(item)
    } catch {
      failed.push(item)
    } finally {
      uploadProgress.value = { current: i + 1, total }
    }
  }

  selectedFotos.value = failed
  if (fileInput.value) fileInput.value.value = ''

  try {
    await loadAll()
  } catch {
    // ignore
  }

  if (okCount > 0 && failed.length === 0) {
    toast.add({ severity: 'success', summary: 'Subida', detail: `Se subieron ${okCount} fotos`, life: 2500 })
  } else if (okCount > 0) {
    toast.add({ severity: 'warn', summary: 'Subida parcial', detail: `Subidas: ${okCount}. Fallidas: ${failed.length}`, life: 3000 })
  } else {
    toast.add({ severity: 'error', summary: 'Error', detail: 'No se pudo subir ninguna foto', life: 3000 })
  }

  uploading.value = false
  uploadProgress.value = null
}

const removeFoto = async (foto: any) => {
  const id = agenciaId.value
  const pid = paqueteId.value
  if (!id || !pid || !foto?.id) return

  try {
    const response: any = await removePaqueteFoto(id, pid, Number(foto.id))
    if (response.success) {
      paquete.value.fotos = (paquete.value.fotos || []).filter((f: any) => f.id !== foto.id)
      toast.add({ severity: 'success', summary: 'Eliminada', detail: 'Foto eliminada', life: 2500 })
    }
  } catch (error: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: error.data?.error?.message || 'No se pudo eliminar', life: 3000 })
  }
}

const resolveFotoUrl = (path?: string) => {
  if (!path) return '/images/placeholder.svg'
  if (path.startsWith('http')) return path
  const clean = path.replace(/^\.?\//, '')

  const apiBase = useRuntimeConfig().public.apiBase as unknown as string
  let origin = ''
  if (typeof apiBase === 'string' && apiBase.startsWith('http')) {
    origin = new URL(apiBase).origin
  } else if (typeof window !== 'undefined') {
    origin = window.location.origin
  }

  return origin ? `${origin}/${clean}` : `/${clean}`
}

// Atracciones
const showAddAtraccionDialog = ref(false)
const atrSearch = ref('')
const atrLoading = ref(false)
const atrResults = ref<any[]>([])
const atrForm = ref({
  dia_numero: 1,
  orden_visita: 1,
  duracion_estimada_horas: 1
})

const nextOrdenForDay = (dia: number) => {
  const items = (paquete.value?.atracciones || []) as any[]
  const relevant = generalIsMultiDay.value ? items.filter((x) => Number(x.dia_numero || 0) === dia) : items
  const maxOrden = Math.max(0, ...relevant.map((x) => Number(x.orden_visita || 0)))
  return maxOrden + 1
}

watch(
  () => atrForm.value.dia_numero,
  (next) => {
    if (!generalIsMultiDay.value) return
    atrForm.value.orden_visita = nextOrdenForDay(Number(next || 1))
  }
)

const searchAtracciones = async () => {
  atrLoading.value = true
  try {
    const response: any = await getAtracciones({
      page: 1,
      limit: 10,
      search: atrSearch.value.trim(),
      status: 'activa'
    } as any)
    if (response.success) {
      atrResults.value = response.data?.atracciones || []
    } else {
      atrResults.value = []
    }
  } catch (error: any) {
    atrResults.value = []
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudieron cargar las atracciones',
      life: 3500
    })
  } finally {
    atrLoading.value = false
  }
}

const openAddAtraccionDialog = async () => {
  showAddAtraccionDialog.value = true
  atrSearch.value = ''
  atrResults.value = []
  atrForm.value = {
    dia_numero: 1,
    orden_visita: nextOrdenForDay(1),
    duracion_estimada_horas: 1
  }
  await searchAtracciones()
}

const addAtraccionToPaquete = async (atr: any) => {
  const id = agenciaId.value
  const pid = paqueteId.value
  if (!id || !pid || !atr?.id) return

  const payload: any = {
    atraccion_id: Number(atr.id),
    orden_visita: Math.max(1, Number(atrForm.value.orden_visita || 1)),
    duracion_estimada_horas: Math.max(1, Number(atrForm.value.duracion_estimada_horas || 1))
  }
  if (generalIsMultiDay.value) payload.dia_numero = Math.max(1, Number(atrForm.value.dia_numero || 1))

  try {
    const response: any = await addPaqueteAtraccion(id, pid, payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Agregado', detail: 'Atracción agregada', life: 2500 })
      await loadAll()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo agregar',
      life: 3500
    })
  }
}

const showEditAtraccionDialog = ref(false)
const savingAtraccion = ref(false)
const selectedAtraccionItem = ref<any>(null)
const editAtrForm = ref({
  dia_numero: 1,
  orden_visita: 1,
  duracion_estimada_horas: 1
})

const openEditAtraccion = (item: any) => {
  selectedAtraccionItem.value = item
  editAtrForm.value = {
    dia_numero: Number(item?.dia_numero || 1),
    orden_visita: Number(item?.orden_visita || 1),
    duracion_estimada_horas: Number(item?.duracion_estimada_horas || 1)
  }
  showEditAtraccionDialog.value = true
}

const saveAtraccionEdit = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  const item = selectedAtraccionItem.value
  if (!id || !pid || !item?.id) return

  const payload: any = {
    orden_visita: Math.max(1, Number(editAtrForm.value.orden_visita || 1)),
    duracion_estimada_horas: Math.max(1, Number(editAtrForm.value.duracion_estimada_horas || 1))
  }
  if (generalIsMultiDay.value) payload.dia_numero = Math.max(1, Number(editAtrForm.value.dia_numero || 1))

  savingAtraccion.value = true
  try {
    const response: any = await updatePaqueteAtraccion(id, pid, Number(item.id), payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Actualizado', life: 2500 })
      showEditAtraccionDialog.value = false
      selectedAtraccionItem.value = null
      await loadAll()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo actualizar',
      life: 3500
    })
  } finally {
    savingAtraccion.value = false
  }
}

const showRemoveAtraccionDialog = ref(false)
const removingAtraccion = ref(false)
const atrItemToRemove = ref<any>(null)
const confirmRemoveAtraccion = (item: any) => {
  atrItemToRemove.value = item
  showRemoveAtraccionDialog.value = true
}

const removeAtraccion = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  const item = atrItemToRemove.value
  if (!id || !pid || !item?.id) return

  removingAtraccion.value = true
  try {
    const response: any = await removePaqueteAtraccion(id, pid, Number(item.id))
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Quitada', detail: 'Atracción eliminada', life: 2500 })
      showRemoveAtraccionDialog.value = false
      atrItemToRemove.value = null
      await loadAll()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo quitar',
      life: 3500
    })
  } finally {
    removingAtraccion.value = false
  }
}

// Salidas
const estadoSalidaOptions = [
  { label: 'Pendiente', value: 'pendiente' },
  { label: 'Activa', value: 'activa' },
  { label: 'Completada', value: 'completada' },
  { label: 'Cancelada', value: 'cancelada' }
]

const showEditSalidaDialog = ref(false)
const savingSalida = ref(false)
const selectedSalida = ref<any>(null)
const salidaForm = ref({
  punto_encuentro: '',
  hora_encuentro: '',
  notas_logistica: '',
  instrucciones_turistas: '',
  guia_nombre: '',
  guia_telefono: '',
  estado: 'pendiente',
  razon_cancelacion: ''
})

const openEditSalida = (s: any) => {
  selectedSalida.value = s
  salidaForm.value = {
    punto_encuentro: s?.punto_encuentro || '',
    hora_encuentro: s?.hora_encuentro ? String(s.hora_encuentro).slice(0, 5) : '',
    notas_logistica: s?.notas_logistica || '',
    instrucciones_turistas: s?.instrucciones_turistas || '',
    guia_nombre: s?.guia_nombre || '',
    guia_telefono: s?.guia_telefono || '',
    estado: s?.estado || 'pendiente',
    razon_cancelacion: s?.razon_cancelacion || ''
  }
  showEditSalidaDialog.value = true
}

const saveSalida = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  const s = selectedSalida.value
  if (!id || !pid || !s?.id) return

  if (salidaForm.value.estado === 'cancelada' && !salidaForm.value.razon_cancelacion.trim()) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'La razón de cancelación es obligatoria', life: 3000 })
    return
  }

  const payload: any = {
    punto_encuentro: salidaForm.value.punto_encuentro?.trim() || '',
    hora_encuentro: salidaForm.value.hora_encuentro?.trim() || '',
    notas_logistica: salidaForm.value.notas_logistica?.trim() || '',
    instrucciones_turistas: salidaForm.value.instrucciones_turistas?.trim() || '',
    guia_nombre: salidaForm.value.guia_nombre?.trim() || '',
    guia_telefono: salidaForm.value.guia_telefono?.trim() || '',
    estado: salidaForm.value.estado,
    razon_cancelacion: salidaForm.value.razon_cancelacion?.trim() || ''
  }

  savingSalida.value = true
  try {
    const response: any = await updatePaqueteSalida(id, pid, Number(s.id), payload)
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Guardado', detail: 'Salida actualizada', life: 2500 })
      showEditSalidaDialog.value = false
      selectedSalida.value = null
      await loadAll()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo actualizar la salida',
      life: 3500
    })
  } finally {
    savingSalida.value = false
  }
}

// Itinerario
const showItinerarioDialog = ref(false)
const itinerarioEditing = ref(false)
const itItem = ref<any>(null)
const itForm = ref({
  dia_numero: 1,
  titulo: '',
  descripcion: '',
  actividades: [] as string[],
  hospedaje_info: ''
})
const savingItinerario = ref(false)

const openItinerarioCreate = () => {
  itinerarioEditing.value = false
  itItem.value = null
  itForm.value = { dia_numero: 1, titulo: '', descripcion: '', actividades: [], hospedaje_info: '' }
  showItinerarioDialog.value = true
}

const openItinerarioEdit = (item: any) => {
  itinerarioEditing.value = true
  itItem.value = item
  itForm.value = {
    dia_numero: Number(item?.dia_numero || 1),
    titulo: item?.titulo || '',
    descripcion: item?.descripcion || '',
    actividades: (item?.actividades || []) as string[],
    hospedaje_info: item?.hospedaje_info || ''
  }
  showItinerarioDialog.value = true
}

const saveItinerario = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  if (!id || !pid) return

  const payload: any = {
    dia_numero: Number(itForm.value.dia_numero || 1),
    titulo: itForm.value.titulo.trim(),
    descripcion: itForm.value.descripcion?.trim() || '',
    actividades: itForm.value.actividades || [],
    hospedaje_info: itForm.value.hospedaje_info?.trim() || ''
  }

  if (!payload.titulo) {
    toast.add({ severity: 'warn', summary: 'Validación', detail: 'El título es obligatorio', life: 3000 })
    return
  }

  savingItinerario.value = true
  try {
    if (itinerarioEditing.value && itItem.value?.id) {
      const response: any = await updateItinerarioItem(id, pid, Number(itItem.value.id), payload)
      if (response.success) {
        toast.add({ severity: 'success', summary: 'Guardado', detail: 'Itinerario actualizado', life: 2500 })
      }
    } else {
      const response: any = await createItinerarioItem(id, pid, payload)
      if (response.success) {
        toast.add({ severity: 'success', summary: 'Agregado', detail: 'Itinerario agregado', life: 2500 })
      }
    }
    showItinerarioDialog.value = false
    await loadAll()
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo guardar el itinerario',
      life: 3500
    })
  } finally {
    savingItinerario.value = false
  }
}

const showDeleteItinerarioDialog = ref(false)
const deletingItinerario = ref(false)
const itToDelete = ref<any>(null)
const confirmDeleteItinerario = (item: any) => {
  itToDelete.value = item
  showDeleteItinerarioDialog.value = true
}
const deleteItinerario = async () => {
  const id = agenciaId.value
  const pid = paqueteId.value
  const item = itToDelete.value
  if (!id || !pid || !item?.id) return

  deletingItinerario.value = true
  try {
    const response: any = await deleteItinerarioItem(id, pid, Number(item.id))
    if (response.success) {
      toast.add({ severity: 'success', summary: 'Eliminado', detail: 'Itinerario eliminado', life: 2500 })
      showDeleteItinerarioDialog.value = false
      itToDelete.value = null
      await loadAll()
    }
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.data?.error?.message || 'No se pudo eliminar',
      life: 3500
    })
  } finally {
    deletingItinerario.value = false
  }
}

onMounted(loadAll)
</script>
