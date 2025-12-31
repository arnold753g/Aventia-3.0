import * as XLSX from 'xlsx'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'

export const exportToExcel = (data: any[], filename: string = 'usuarios') => {
  const worksheet = XLSX.utils.json_to_sheet(data)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, 'Usuarios')
  XLSX.writeFile(workbook, `${filename}_${Date.now()}.xlsx`)
}

export const exportToPDF = (data: any[], filename: string = 'usuarios') => {
  const doc = new jsPDF()

  doc.text('ANDARIA - Lista de Usuarios', 14, 15)
  doc.setFontSize(10)
  doc.text(`Generado: ${new Date().toLocaleString('es-BO')}`, 14, 22)

  autoTable(doc, {
    startY: 30,
    head: [['ID', 'Nombre', 'Email', 'CI', 'Rol', 'Estado']],
    body: data.map(u => [
      u.id,
      `${u.nombre} ${u.apellido_paterno}`,
      u.email,
      u.ci,
      u.rol,
      u.status
    ]),
    styles: { fontSize: 8 },
    headStyles: { fillColor: [107, 123, 78] }
  })

  doc.save(`${filename}_${Date.now()}.pdf`)
}
