package form

//Register Form
type Register struct {
	UserName string `binding:"Required;AlphaDashDot;MaxSize(50)"`
	Email    string `binding:"Required;Email;MaxSize(254)"`
	Password string `binding:"Required;MaxSize(40)"`
}

type SignIn struct {
	UserName string `binding:"Required;MaxSize(255)"`
	Password string `binding:"Required;MaxSize(255)"`
	Remember bool
}
