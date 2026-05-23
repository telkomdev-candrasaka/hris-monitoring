// @title HRIS Monitoring API
// @version 1.0
// @description Backend HRIS dan operasional frozen-food retail dengan JWT auth, geofenced attendance + selfie, leave approval with staffing risk, payroll, PPE/asset compliance, dan executive reports.
// @host localhost:8080
// @basePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: Bearer <token>

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/config"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/handlers"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/middlewares"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/models"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/repositories"
	"github.com/telkomdev-candrasaka/hris-monitoring.git/internal/services"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/telkomdev-candrasaka/hris-monitoring.git/docs"
)

func main() {
	fmt.Println("Mulai menjalankan aplikasi HRIS Monitoring...")

	// Panggil fungsi untuk menyambungkan database
	config.ConnectDatabase()

	// AutoMigrate untuk membuat tabel sesuai model
	err := config.DB.AutoMigrate(&models.Location{}, &models.Shift{}, &models.User{}, &models.Attendance{}, &models.Leave{}, &models.Payroll{}, &models.Asset{}, &models.MandatoryEquipment{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	// Inisialisasi layer repository, service, dan handler untuk Location
	locationRepo := repositories.NewLocationRepository()
	locationService := services.NewLocationService(locationRepo)
	locationHandler := handlers.NewLocationHandler(locationService)

	// Inisialisasi layer repository, service, dan handler untuk User
	userRepo := repositories.NewUserRepository()
	assetRepo := repositories.NewAssetRepository()
	userService := services.NewUserServiceWithAsset(userRepo, assetRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Inisialisasi auth service dan login route
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	authHandler.RegisterRoutes()

	// Protected endpoints
	http.Handle("/locations", middlewares.JWTAuth(http.HandlerFunc(locationHandler.LocationsHandler)))
	http.Handle("/locations/", middlewares.JWTAuth(http.HandlerFunc(locationHandler.LocationByIDHandler)))
	http.Handle("/users", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr", "manager_outlet")(http.HandlerFunc(userHandler.UsersHandler))))
	http.Handle("/users/", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr", "manager_outlet")(http.HandlerFunc(userHandler.UserByIDHandler))))

	attendanceRepo := repositories.NewAttendanceRepository()
	mandatoryEquipmentRepo := repositories.NewMandatoryEquipmentRepository()
	attendanceService := services.NewAttendanceService(attendanceRepo, locationRepo, userRepo, assetRepo, mandatoryEquipmentRepo)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	http.Handle("/attendance/checkin", middlewares.JWTAuth(http.HandlerFunc(attendanceHandler.CheckInHandler)))
	http.Handle("/attendance/checkout", middlewares.JWTAuth(http.HandlerFunc(attendanceHandler.CheckOutHandler)))
	http.Handle("/attendance", middlewares.JWTAuth(http.HandlerFunc(attendanceHandler.GetAttendancesHandler)))
	http.Handle("/attendance/", middlewares.JWTAuth(http.HandlerFunc(attendanceHandler.GetAttendanceByIDHandler)))

	leaveRepo := repositories.NewLeaveRepository()
	leaveService := services.NewLeaveService(leaveRepo, userRepo)
	leaveHandler := handlers.NewLeaveHandler(leaveService)
	http.Handle("/leaves", middlewares.JWTAuth(http.HandlerFunc(leaveHandler.LeavesHandler)))
	http.Handle("/leaves/approve/", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr", "manager_outlet")(http.HandlerFunc(leaveHandler.ApproveLeaveHandler))))
	http.Handle("/leaves/reject/", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr", "manager_outlet")(http.HandlerFunc(leaveHandler.RejectLeaveHandler))))
	http.Handle("/leaves/", middlewares.JWTAuth(http.HandlerFunc(leaveHandler.LeaveByIDHandler)))

	payrollRepo := repositories.NewPayrollRepository()
	payrollService := services.NewPayrollService(payrollRepo, attendanceRepo, userRepo)
	payslipPDFService := services.NewPayslipPDFService()
	payrollHandler := handlers.NewPayrollHandlerWithPDF(payrollService, payslipPDFService, userRepo)
	http.Handle("/payrolls", middlewares.JWTAuth(http.HandlerFunc(payrollHandler.PayrollHandler)))
	http.Handle("/payrolls/history", middlewares.JWTAuth(http.HandlerFunc(payrollHandler.PayrollHistoryHandler)))
	http.Handle("/payrolls/download", middlewares.JWTAuth(http.HandlerFunc(payrollHandler.DownloadPayslipPDFHandler)))

	assetService := services.NewAssetService(assetRepo)
	assetHandler := handlers.NewAssetHandler(assetService)
	http.Handle("/assets", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr")(http.HandlerFunc(assetHandler.AssetsHandler))))
	http.Handle("/assets/", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr")(http.HandlerFunc(assetHandler.AssetByIDHandler))))
	http.Handle("/assets/borrow/", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr")(http.HandlerFunc(assetHandler.BorrowAssetHandler))))
	http.Handle("/assets/return/", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr")(http.HandlerFunc(assetHandler.ReturnAssetHandler))))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	reportRepo := repositories.NewReportRepository()
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)
	http.Handle("/api/reports/labour-cost-leakage", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr")(http.HandlerFunc(reportHandler.LabourCostLeakageHandler))))
	http.Handle("/api/reports/attendance-risk", middlewares.JWTAuth(middlewares.RoleAuthorize("admin_hr")(http.HandlerFunc(reportHandler.AttendanceRiskHandler))))

	fmt.Println("Aplikasi siap digunakan! Server berjalan di http://localhost:8080")
	fmt.Println("Dokumentasi Swagger: http://localhost:8080/swagger/index.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Penting: Di net/http, "/" bertindak sebagai catch-all (menangkap semua rute yang tidak terdaftar).
        // Pengecekan ini memastikan pesan hanya muncul jika URL benar-benar "/", bukan "/url-asal-asalan"
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Welcome to HRIS Monitoring API!", "status": "running"}`))
    })
	log.Fatal(http.ListenAndServe(":8080", nil))
}
