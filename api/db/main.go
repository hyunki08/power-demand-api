package db

func Run() {
	connectDB()
	defer disconnectDB()

	var pd = powerDemand{}
	pd.setMetaCollection()
}
