package globalvars

var (
	Email                 = ""
	Password              = ""
	ProceedingsCheckIndex = 0

	ApplicationJson  = "application/json"
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"
	Origin           = "https://inpol.mazowieckie.pl"
	LoginPageUrl     = "https://inpol.mazowieckie.pl/login"
	HomePageUrl      = "https://inpol.mazowieckie.pl/home"
	HomePageCasesUrl = "https://inpol.mazowieckie.pl/home/cases/%s"

	HtmlAcceptHeader     = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	KeepAliveHeader      = "keep-alive"
	AcceptEncodingHeader = "gzip, deflate, br, zstd"

	LoginRequestUrl                          = "https://inpol.mazowieckie.pl/identity/sign-in"
	GetActiveProceedingsRequestUrl           = "https://inpol.mazowieckie.pl/api/foreigner/active-proceedings"
	GetProceedingRequestUrl                  = "https://inpol.mazowieckie.pl/api/proceedings/%s" // GetProceedingRequestUrl = "https://inpol.mazowieckie.pl/api/proceedings/:proceedingId"
	GetProceedingReservationQueuesRequestUrl = "https://inpol.mazowieckie.pl/api/proceedings/%s/reservationQueues"
	GetReservationQueueDatesRequestUrl       = "https://inpol.mazowieckie.pl/api/reservations/queue/%s/dates"
	GetReservationQueueDateSlotsRequestUrl   = "https://inpol.mazowieckie.pl/api/reservations/queue/%s/%s/slots" // reservations/queue/:queueId/:date/slots
	ReserveAppointmentRequestUrl             = "https://inpol.mazowieckie.pl/api/reservations/queue/%s/reserve"  // reservations/queue/:queueId/reserve

	AppointmentMade = "AppointmentMade"
	Created         = "Created"
)
