package models

//User ...

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

//Login ...

func (m UserModel) Login(login string) (token Token, err error) {

	//Generate the JWT auth token

	tokenDetails, err := authModel.CreateToken(login)

	saveErr := authModel.CreateAuth(login, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	return token, nil
}

//Register ...
//func (m UserModel) Register(form forms.RegisterForm) (user User, err error) {
//	getDb := db.GetDB()
//
//	//Check if the user exists in database
//	checkUser, err := getDb.SelectInt("SELECT count(id) FROM public.user WHERE mail=LOWER($1) LIMIT 1", form.Email)
//
//	if err != nil {
//		return user, err
//	}
//
//	if checkUser > 0 {
//		return user, errors.New("User already exists")
//	}
//
//	//Create the user and return back the user ID
//	err = getDb.QueryRow("INSERT INTO public.user(mail, ao) VALUES($1, $2) RETURNING id", form.Email, form.Ao).Scan(&user.ID)
//
//	user.AO = form.Ao
//	user.Email = form.Email
//
//	return user, err
//}
