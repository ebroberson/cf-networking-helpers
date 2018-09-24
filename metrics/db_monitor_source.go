package metrics

type Db interface {
	OpenConnections() int
}

func NewDBMonitorSource(Db Db) MetricSource {
	return MetricSource{
		Name: "DBOpenConnections",
		Unit: "",
		Getter: func() (float64, error) {
			return float64(Db.OpenConnections()), nil
		},
	}
}
