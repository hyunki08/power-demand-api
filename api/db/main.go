package db

var PDCollection powerDemand

func Run() {
	connectDB()

	PDCollection = powerDemand{}
	PDCollection.setMetaCollection()
}

func Disconnect() {
	disconnectDB()
}
