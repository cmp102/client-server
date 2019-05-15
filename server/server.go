package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Crt        string `yaml:"crt"`
	Key        string `yaml:"key"`
	Port       string `yaml:"port"`
	Storage    string `yaml:"storage"`
	AuthTokens string `yaml:"authtokens"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
}

var conf config
var authTokens map[string]string

func main() {
	file, err := os.Open("config.yml")

	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}

	file, err = os.Open("tokens.yml")

	if err != nil {
		log.Fatal(err)
	}

	data, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, &authTokens)

	fmt.Println(authTokens["asdf"])

	http.HandleFunc("/post", myPost)
	http.HandleFunc("/tokengenerator", generateToken)

	http.HandleFunc(
		"/",
		BasicAuth(
			http.FileServer(http.Dir(conf.Storage)).ServeHTTP,
			conf.User,
			conf.Password,
			"Please enter your username and password for this site"))

	fmt.Printf("Listening on port %s...", conf.Port)
	log.Fatal(http.ListenAndServeTLS(":"+conf.Port, conf.Crt, conf.Key, nil))
}

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func myPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.NotFound(w, r)
		return
	}
	authtoken := r.FormValue("authtoken")
	fmt.Printf("AuthToken is: %s\n", authtoken)
	client, ok := authTokens[authtoken]
	os.Mkdir(conf.Storage+"/"+client, 0777)
	if !ok {
		fmt.Printf("Unauthorised %s\n", authtoken)
		return
	}
	infile, header, err := r.FormFile("sendfile")
	if err != nil {
		fmt.Println("Error parsing uploaded file: " + err.Error())
		return
	}
	defer infile.Close()
	outfilename := conf.Storage + "/" + client + "/" + header.Filename

	outfile, err := os.Create(outfilename)
	defer outfile.Close()
	if err != nil {
		fmt.Println("Error creating file in server: " + err.Error())
		return
	}
	_, err = io.Copy(outfile, infile)
	if err != nil {

		fmt.Println("Error saving file in server: " + err.Error())
		return
	}
	fmt.Printf("Image writed\n")
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, `<html>
							<form action="/tokengenerator" method="POST">
								Device name:<br>
								<input type="text" name="devicename">
								<button type="submit">Generate</button>
							</form> 
						</html>`)
		return
	} else if r.Method == "POST" {
		b := make([]byte, 24)
		rand.Read(b)
		token := base64.URLEncoding.EncodeToString(b)
		fmt.Fprintf(w, "The token for %s is:\n%s", r.FormValue("devicename"), token)
	}
	// authtoken := r.FormValue("authtoken")
	// fmt.Printf("AuthToken is: %s\n", authtoken)
	// client, ok := authTokens[authtoken]
	// os.Mkdir(conf.Storage+"/"+client, 0777)
	// if !ok {
	// 	fmt.Printf("Unauthorised %s\n", authtoken)
	// 	return
	// }
	// infile, header, err := r.FormFile("sendfile")
	// if err != nil {
	// 	fmt.Println("Error parsing uploaded file: " + err.Error())
	// 	return
	// }
	// defer infile.Close()
	// outfilename := conf.Storage + "/" + client + "/" + header.Filename

	// outfile, err := os.Create(outfilename)
	// defer outfile.Close()
	// if err != nil {
	// 	fmt.Println("Error creating file in server: " + err.Error())
	// 	return
	// }
	// _, err = io.Copy(outfile, infile)
	// if err != nil {

	// 	fmt.Println("Error saving file in server: " + err.Error())
	// 	return
	// }
	// fmt.Printf("Image writed\n")
}

// func generateToken() (string, error) {

// }
