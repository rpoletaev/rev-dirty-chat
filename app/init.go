package app

import (
	"github.com/cbonello/revel-csrf"
	"github.com/revel/revel"
	"github.com/rpoletaev/rev-dirty-chat/utilities/helper"
	"github.com/rpoletaev/rev-dirty-chat/utilities/mongo"
	"github.com/rpoletaev/rev-dirty-chat/utilities/tracelog"
	"os"
)

//** PRIVATE FUNCTIONS

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		csrf.CSRFFilter,               // CSRP prevention
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}

	// Called to help init the application
	revel.OnAppStart(initApp)
}

// initApp contains all application level initialization
func initApp() {
	// Capture the global email settings
	tracelog.EmailHost = revel.Config.StringDefault("email.host", "")
	tracelog.EmailPort = revel.Config.IntDefault("email.port", 0)
	tracelog.EmailUserName = revel.Config.StringDefault("email.username", "")
	tracelog.EmailPassword = revel.Config.StringDefault("email.password", "")
	tracelog.EmailTo = revel.Config.StringDefault("email.to", "")
	tracelog.EmailAlertSubject = revel.Config.StringDefault("email.alert_subject", "")

	// MongoDB Settings
	helper.MONGO_DATABASE = revel.Config.StringDefault("mgo.use_database", "")

	// Init mongo
	err := mongo.Startup(helper.MAIN_GO_ROUTINE)
	if err != nil {
		tracelog.COMPLETED_ERROR(err, helper.MAIN_GO_ROUTINE, "initApp")
		os.Exit(1)
	}
}
