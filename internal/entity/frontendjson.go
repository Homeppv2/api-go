package entity

type MessageFromFrontendJSON struct {
	Id int `json:"id"`
	/*

		id - сообщений которые приходят с frontend
		600 - запроос на получение сообщений из архива с id - from
		601 - отправить запрос на контролеры TODO

		id - сообщений которые приходят с backend
		801 - отправить на фронтент несколько сообщений о статусе контролеров

	*/
	Rng  IdMessages             `json:"rng,omitempty"`
	Msgs []MessangeTypeZiroJson `json:"msgs,omitempty"`
}

type IdMessages struct {
	From  int `json:"from"`
	To    int `json:"to"`
	Count int `json:"count"`
}
