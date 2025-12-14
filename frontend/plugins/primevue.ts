import { defineNuxtPlugin } from '#app'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Toast from 'primevue/toast'
import ToastService from 'primevue/toastservice'
import Calendar from 'primevue/calendar'
import Select from 'primevue/select'
import Dropdown from 'primevue/dropdown'
import MultiSelect from 'primevue/multiselect'
import Password from 'primevue/password'
import Message from 'primevue/message'
import Skeleton from 'primevue/skeleton'
import Chip from 'primevue/chip'
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import Tag from 'primevue/tag'
import Tooltip from 'primevue/tooltip'
import DatePicker from 'primevue/datepicker'
import Form from '@primevue/forms/form'
import FormField from '@primevue/forms/formfield'
import Checkbox from 'primevue/checkbox'
import SelectButton from 'primevue/selectbutton'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import ProgressSpinner from 'primevue/progressspinner'
import Divider from 'primevue/divider'
import InputMask from 'primevue/inputmask'
import InputSwitch from 'primevue/inputswitch'

const localeEs = {
  firstDayOfWeek: 1,
  dayNames: ['Domingo', 'Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado'],
  dayNamesShort: ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb'],
  dayNamesMin: ['D', 'L', 'M', 'M', 'J', 'V', 'S'],
  monthNames: [
    'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
    'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'
  ],
  monthNamesShort: ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'],
  today: 'Hoy',
  clear: 'Borrar'
}

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(PrimeVue, {
    theme: {
      preset: Aura,
      options: {
        prefix: 'p',
        darkModeSelector: '.dark',
        cssLayer: false
      }
    },
    locale: localeEs,
    ripple: true
  })
  nuxtApp.vueApp.use(ToastService)

  // Registrar directivas
  nuxtApp.vueApp.directive('tooltip', Tooltip)

  // Registrar componentes globalmente
  nuxtApp.vueApp.component('Button', Button)
  nuxtApp.vueApp.component('InputText', InputText)
  nuxtApp.vueApp.component('Card', Card)
  nuxtApp.vueApp.component('DataTable', DataTable)
  nuxtApp.vueApp.component('Column', Column)
  nuxtApp.vueApp.component('Dialog', Dialog)
  nuxtApp.vueApp.component('Toast', Toast)
  nuxtApp.vueApp.component('Calendar', Calendar)
  nuxtApp.vueApp.component('Select', Select)
  nuxtApp.vueApp.component('Dropdown', Dropdown)
  nuxtApp.vueApp.component('MultiSelect', MultiSelect)
  nuxtApp.vueApp.component('Password', Password)
  nuxtApp.vueApp.component('Message', Message)
  nuxtApp.vueApp.component('Skeleton', Skeleton)
  nuxtApp.vueApp.component('Chip', Chip)
  nuxtApp.vueApp.component('IconField', IconField)
  nuxtApp.vueApp.component('InputIcon', InputIcon)
  nuxtApp.vueApp.component('Tag', Tag)
  nuxtApp.vueApp.component('DatePicker', DatePicker)
  nuxtApp.vueApp.component('Form', Form)
  nuxtApp.vueApp.component('FormField', FormField)
  nuxtApp.vueApp.component('Checkbox', Checkbox)
  nuxtApp.vueApp.component('SelectButton', SelectButton)
  nuxtApp.vueApp.component('Textarea', Textarea)
  nuxtApp.vueApp.component('InputNumber', InputNumber)
  nuxtApp.vueApp.component('ProgressSpinner', ProgressSpinner)
  nuxtApp.vueApp.component('Divider', Divider)
  nuxtApp.vueApp.component('InputMask', InputMask)
  nuxtApp.vueApp.component('InputSwitch', InputSwitch)
})
