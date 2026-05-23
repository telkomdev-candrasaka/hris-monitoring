package services

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
)

type PayslipPDFService struct {
}

func NewPayslipPDFService() *PayslipPDFService {
	return &PayslipPDFService{}
}

func (s *PayslipPDFService) GeneratePayslipPDF(payroll *models.Payroll, user *models.User) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "SLIP GAJI", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	pdf.Ln(5)

	pdf.Cell(50, 8, "Nama Karyawan")
	pdf.CellFormat(100, 8, ": "+user.Name, "", 1, "", false, 0, "")

	pdf.Cell(50, 8, "Email")
	pdf.CellFormat(100, 8, ": "+user.Email, "", 1, "", false, 0, "")

	pdf.Cell(50, 8, "Periode")
	periodStr := fmt.Sprintf("%d-%02d", payroll.Year, payroll.Month)
	pdf.CellFormat(100, 8, ": "+periodStr, "", 1, "", false, 0, "")

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, "RINCIAN GAJI", "", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(80, 8, "Gaji Pokok")
	pdf.CellFormat(60, 8, fmt.Sprintf("Rp %.2f", payroll.BaseSalary), "", 1, "R", false, 0, "")

	pdf.Cell(80, 8, "Tunjangan Cold Storage")
	pdf.CellFormat(60, 8, fmt.Sprintf("Rp %.2f", payroll.ColdStorageAllowance), "", 1, "R", false, 0, "")

	pdf.Cell(80, 8, "Uang Lembur")
	pdf.CellFormat(60, 8, fmt.Sprintf("Rp %.2f", payroll.OvertimePay), "", 1, "R", false, 0, "")

	pdf.Cell(80, 8, "Gaji Kotor")
	pdf.CellFormat(60, 8, fmt.Sprintf("Rp %.2f", payroll.GrossSalary), "", 1, "R", false, 0, "")

	pdf.Cell(80, 8, "Kehadiran")
	pdf.CellFormat(60, 8, fmt.Sprintf("%d hari", payroll.TotalAttendance), "", 1, "R", false, 0, "")

	pdf.Cell(80, 8, "Absensi")
	pdf.CellFormat(60, 8, fmt.Sprintf("-%d hari", payroll.AbsentCount), "", 1, "R", false, 0, "")

	pdf.Cell(80, 8, "Keterlambatan")
	pdf.CellFormat(60, 8, fmt.Sprintf("-%d kali", payroll.LateCount), "", 1, "R", false, 0, "")

	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(80, 8, "Total Potongan")
	pdf.CellFormat(60, 8, fmt.Sprintf("Rp %.2f", payroll.TotalDeduction), "", 1, "R", false, 0, "")

	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(80, 10, "GAJI BERSIH")
	pdf.CellFormat(60, 10, fmt.Sprintf("Rp %.2f", payroll.NetSalary), "", 1, "R", false, 0, "")

	pdf.Ln(10)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("Digenerate pada: %s", payroll.GeneratedAt.Format("02-01-2006 15:04:05")), "", 1, "C", false, 0, "")

	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
