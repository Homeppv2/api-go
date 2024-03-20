package entity

type ControllersData struct {
	Id_contorller     int `db:"id_contorller"`
	Type_controller   int `db:"type_controller"`
	Number_controller int `db:"number_controller"`
}

type MainMessangesData struct {
	Id_messange               int `db:"id_messange"`
	Id_contorller             int `db:"id_contorller"`
	Status_controller         int `db:"status_controller"`
	Charge_controller         int `db:"charge_controller"`
	Temperature_MK_controller int `db:"temperature_MK_controller"`
}

type ContollersLeackMessangesData struct {
	Id_messange int `db:"id_messange"`
	Leack       int `db:"leack"`
}

type ContollersModuleMessangesData struct {
	Id_messange int `db:"id_messange"`
	Temperature int `db:"temperature"`
	Humidity    int `db:"humidity"`
	Pressure    int `db:"pressure"`
	Gas         int `db:"gas"`
}

type ContollersEnviromentMessangesData struct {
	Id_messange int `db:"id_messange"`
	Temperature int `db:"temperature"`
	Humidity    int `db:"humidity"`
	Pressure    int `db:"pressure"`
	Voc         int `db:"VOC"`
	Gas1        int `db:"gas1"`
	Gas2        int `db:"gas2"`
	Gas3        int `db:"gas3"`
	Pm1         int `db:"pm1"`
	Pm25        int `db:"pm25"`
	Pm10        int `db:"pm10"`
	Fire        int `db:"fire"`
	smoke       int `db:"smoke"`
}

type User struct {
	ID           int    `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	HashPassword string `json:"hashPassword" db:"hash_password"`
}
