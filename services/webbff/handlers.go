package webbff

import (
	"net/http"

	"github.com/TFTPL/AWS-Cost-Calculator/services/coreservices/user"
	"github.com/gorilla/mux"
)

func InitHandlers(router *mux.Router) {
	// login
	//router.PathPrefix("/").HandlerFunc(loginPageHandler)
	router.HandleFunc("/", loginPageHandler).Methods("GET", "POST")
	// signup
	router.HandleFunc("/signup/", signupPageHandler).Methods("GET", "POST")
	// dashboard
	router.HandleFunc("/{uid:[0-9]+}/days/{filter:[0-9]+}", user.VerifyAccess(dashboardHandler)).Methods("GET")
	//list projects
	router.HandleFunc("/{uid:[0-9]+}/projects/", user.VerifyAccess(getProjectByUserIdHandler)).Methods("GET")
	// project dashboard - graphs of resource usage
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/{filter:[0-9]+}", user.VerifyAccess(graphPageHandler)).Methods("GET")

	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/{filter:[0-9]+}/service", user.VerifyAccess(graphPageServiceHandler)).Methods("GET")
	// new project
	// router.HandleFunc("/{uid:[0-9]+}/projects/new/", user.VerifyAccess(newProjectPageHandler)).Methods("GET", "POST")
	// project settings
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/settings/", user.VerifyAccess(projectSettingsPageHandler)).Methods("GET", "POST")
	// export project
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/{filter:[0-9]+}/export", user.VerifyAccess(projectExportPageHandler))
	// export services cost by services
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/{filter:[0-9]+}/export/service/{value:.+}", user.VerifyAccess(projectExportPageHandlerByService))
	// export services cost by fromdate to Todate
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/export/date", user.VerifyAccess(projectExportPageHandlerByDate))
	// export services cost by year
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/export/year", projectExportPageHandlerByYear)
	// logout
	router.HandleFunc("/logout/", logoutPageHandler).Methods("GET")
	//Delete Project
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/delete/", deleteProjectHandler).Methods("GET")
	//delete aws credentials
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/deleteAwsCreds/", deleteAwsCredentialsHandler).Methods("DELETE")
	//UpdatePinned Project
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/pinproject/{pin:.+}", updatePinnedProject).Methods("GET")
	//updatePinnedProjectByProjectName
	router.HandleFunc("/{uid:[0-9]+}/pinproject/{pin:.+}/{pName:.+}", updatePinnedProjectByProjectName).Methods("GET")
	// export services cost by services with custom date
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/export/services/dates/{value:.+}", user.VerifyAccess(projectExportPageHandlerByServiceAndCustomDate))
	// export services cost by services
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/export/services/years/{value:.+}", user.VerifyAccess(projectExportPageHandlerByServiceWithYear))
	// project settings handler
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/project-setting/", projectSettingHandler).Methods("GET", "POST")
	// router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/project/credentials/", projectCredentialsHandler).Methods("GET", "POST")

	//UpdateActive Project
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/activeproject/{active:.+}", updateActiveProject).Methods("GET")
	// add new project
	router.HandleFunc("/{uid:[0-9]+}/projects/new/", newProjectHandler).Methods("GET", "POST")
	// graph page pinned button
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/graphpinned/{pin:.+}", updateByGraphPinned).Methods("GET")
	// for mongo download
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/{filter:[0-9]+}/mongoExport", user.VerifyAccess(mongoProjectExportPageHandler))
	// export mongo services cost by fromdate to Todate
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/mongoExport/date", user.VerifyAccess(mongoProjectExportPageHandlerByDate))
	// export mongo services cost by year
	router.HandleFunc("/{uid:[0-9]+}/projects/{pid:[0-9]+}/mongoExport/year", user.VerifyAccess(mongoProjectExportPageHandlerByYear))

	// serve static content
	serveStaticFiles(router)
	// defaulthandler
	router.HandleFunc("/{path:.*}", defaultHandler)
}

func serveStaticFiles(router *mux.Router) {
	fs := http.FileServer(http.Dir("services/webbff/public"))
	handler := http.StripPrefix("/public/", fs)
	router.PathPrefix("/public/").Handler(handler)
}
