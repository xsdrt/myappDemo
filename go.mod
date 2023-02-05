module myappDemo

go 1.19

replace github.com/xsdrt/hiSpeed => ../hiSpeed //whenever  get a request to go get this pkg; go up one folder and use whats in hiSpeed... so instead of going to github will just take the contents of this folder instead (hiSpeed)...

require github.com/xsdrt/hiSpeed v0.0.0-00010101000000-000000000000

require github.com/alexedwards/scs/v2 v2.5.0 // indirect

require (
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet/v6 v6.2.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/joho/godotenv v1.4.0 // indirect
)
