package helper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/AKSHAYK0UL/Email_Auth/model"
	"github.com/AKSHAYK0UL/Email_Auth/smtphost"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// Check if the email entered by the user is allowed or not
func IfEmailIsAllowed(useremail string) (string, error) {
	if val := smtphost.Checkthedomain(useremail); val != "try another email provider" {
		return val, nil

	}
	return "", errors.New("try another email provider")

}

// Generate Verification code
func generateVcode() string {
	rand.New(rand.NewSource(time.Now().UnixMicro()))
	vcode := rand.Intn(90000) + 10000
	return strconv.Itoa(vcode)
}

// send email to the user with verification code
func SendEmail(domailValue string, requestData model.RequestModel) (model.SignUpResponse, error) {
	godotenv.Load()
	filteremail := bson.D{{Key: "useremail", Value: requestData.UserEmail}}
	filtername := bson.D{{Key: "username", Value: requestData.UserName}}
	filterphone_no := bson.D{{Key: "phone", Value: requestData.Phone}}
	var emailExist bson.M
	var nameExist bson.M
	var phoneExist bson.M
	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filteremail).Decode(&emailExist)
	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filtername).Decode(&nameExist)
	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filterphone_no).Decode(&phoneExist)

	if emailExist == nil && nameExist == nil && phoneExist == nil {

		fromemail := os.Getenv("Email")
		apppassword := os.Getenv("AppPassword")
		host := smtphost.Getsmtphost(domailValue)
		verificationcode := generateVcode()
		htmlval := model.HtmlVar{UserName: requestData.UserName, Vcode: verificationcode}

		t, err := template.ParseFiles("html/Mailhtml.html")
		if err != nil {
			return model.SignUpResponse{}, err
		}

		buff := new(bytes.Buffer)
		if err := t.Execute(buff, htmlval); err != nil {

			return model.SignUpResponse{}, err
		}
		m := gomail.NewMessage()
		m.SetHeader("From", fromemail)
		m.SetHeader("To", requestData.UserEmail)
		m.SetHeader("Subject", "Koul Network")
		m.SetBody("text/html", buff.String())

		d := gomail.NewDialer(host, 587, fromemail, apppassword)
		//store the response data
		requestDataForresponse := &model.SignUpResponse{}
		// Send the email
		if err := d.DialAndSend(m); err != nil {
			return model.SignUpResponse{}, err
		} else {
			fmt.Println("Email Sent Successfully!")
			requestData.SendAt = time.Now().Format("1504")
			requestData.Vcode = htmlval.Vcode
			if requestData.Password != "" {
				hashpassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
				requestData.Password = string(hashpassword)
				if err != nil {
					return model.SignUpResponse{}, err
				}
			}

			hashVcode, err := bcrypt.GenerateFromPassword([]byte(verificationcode), bcrypt.DefaultCost)
			requestData.Vcode = string(hashVcode)
			if err != nil {
				return model.SignUpResponse{}, err
			}
			insertedId, err := model.MongoInstance.Mdatabase.Collection("Verification code").InsertOne(context.Background(), requestData)
			if err != nil {
				return model.SignUpResponse{}, err
			}

			filter := bson.D{{Key: "_id", Value: insertedId.InsertedID}}
			responsevalues := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
			if err := responsevalues.Decode(requestDataForresponse); err != nil {
				return model.SignUpResponse{}, err
			} else {
				requestDataForresponse.Status = "202"
			}

		}
		return *requestDataForresponse, nil
	}
	if emailExist != nil && nameExist != nil && phoneExist == nil {
		return model.SignUpResponse{}, errors.New("email, username and phone no. already exist")
	} else if emailExist != nil && nameExist != nil {
		return model.SignUpResponse{}, errors.New("email and username already exist")

	} else if emailExist != nil && phoneExist != nil {
		return model.SignUpResponse{}, errors.New("email and phone no. already exist")

	} else if phoneExist != nil && nameExist != nil {
		return model.SignUpResponse{}, errors.New("phone no. and username already exist")

	} else if emailExist != nil {
		return model.SignUpResponse{}, errors.New("email already exist")
	} else if phoneExist != nil {
		return model.SignUpResponse{}, errors.New("phone no. already exist")

	}
	return model.SignUpResponse{}, errors.New("username already exist")
}

// verify Code
func VerifyCode(userid string, vcode string) (model.UserAccount, error) {
	id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return model.UserAccount{}, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	response := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
	reqresponseData := &model.SignUpResponse{}
	if err := response.Decode(reqresponseData); err != nil {
		return model.UserAccount{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(reqresponseData.Vcode), []byte(vcode))
	if err != nil {
		return model.UserAccount{}, err
	}

	useraccount := model.UserAccountStoreDb{AuthType: "Email Auth", UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Password: reqresponseData.Password, CreateAt: time.Now(), UpdateAt: time.Now(), Phone: reqresponseData.Phone}
	result, err := model.MongoInstance.Mdatabase.Collection("Account").InsertOne(context.Background(), useraccount)
	if err != nil {
		return model.UserAccount{}, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	account := model.UserAccount{AuthType: "Email Auth", UserId: insertedID, UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Phone: reqresponseData.Phone}
	return account, nil
}

// reset password
func ResetpasswordSendEmail(domailValue string, requestData model.RequestModel) (model.SignUpResponse, error) {
	filter := bson.D{{Key: "useremail", Value: requestData.UserEmail}, {Key: "authtype", Value: "Email Auth"}}

	userExist := model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filter)

	user := &model.UserAccount{}
	if err := userExist.Decode(&user); err != nil {
		fmt.Println("+NO ACCOUNT ON DECODE")
		return model.SignUpResponse{}, errors.New("no user found")
	}
	fromemail := os.Getenv("Email")
	apppassword := os.Getenv("AppPassword")
	host := smtphost.Getsmtphost(domailValue)
	verificationcode := generateVcode()
	htmlval := model.HtmlVar{UserName: user.UserName, Vcode: verificationcode}

	t, err := template.ParseFiles("html/resetmail.html")
	if err != nil {
		return model.SignUpResponse{}, err
	}

	buff := new(bytes.Buffer)
	if err := t.Execute(buff, htmlval); err != nil {

		return model.SignUpResponse{}, err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", fromemail)
	m.SetHeader("To", user.UserEmail)
	m.SetHeader("Subject", "Koul Network")
	m.SetBody("text/html", buff.String())

	d := gomail.NewDialer(host, 587, fromemail, apppassword)
	//store the response data
	requestDataForresponse := &model.SignUpResponse{}
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return model.SignUpResponse{}, err
	} else {
		fmt.Println("Email Sent Successfully!")
		requestData.SendAt = time.Now().Format("15:04:05")
		requestData.Vcode = htmlval.Vcode
		hashpassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
		requestData.Password = string(hashpassword)
		if err != nil {
			return model.SignUpResponse{}, err
		}

		hashVcode, err := bcrypt.GenerateFromPassword([]byte(verificationcode), bcrypt.DefaultCost)
		requestData.Vcode = string(hashVcode)
		if err != nil {
			return model.SignUpResponse{}, err
		}
		insertedId, err := model.MongoInstance.Mdatabase.Collection("Verification code").InsertOne(context.Background(), requestData)
		if err != nil {
			return model.SignUpResponse{}, err
		}

		filter := bson.D{{Key: "_id", Value: insertedId.InsertedID}}
		responsevalues := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
		if err := responsevalues.Decode(requestDataForresponse); err != nil {
			return model.SignUpResponse{}, err
		} else {
			requestDataForresponse.Status = "202"
		}

		return *requestDataForresponse, nil
	}

}
func Resetverify(userid string, vcode string) (model.UserAccount, error) {
	id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return model.UserAccount{}, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	response := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
	reqresponseData := &model.SignUpResponse{}
	if err := response.Decode(reqresponseData); err != nil {
		return model.UserAccount{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(reqresponseData.Vcode), []byte(vcode))
	if err != nil {
		return model.UserAccount{}, err
	}

	useraccount := model.UserAccountStoreDb{AuthType: "Email Auth", UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Password: reqresponseData.Password, UpdateAt: time.Now()}
	fmt.Println("USER EMALI $$$ ", useraccount.UserEmail)
	updatefilter := bson.D{{Key: "useremail", Value: useraccount.UserEmail}}
	newData := bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key: "password", Value: useraccount.Password,
		},
			{
				Key: "updated_at", Value: useraccount.UpdateAt,
			},
		},
	},
	}

	_, err = model.MongoInstance.Mdatabase.Collection("Account").UpdateOne(context.Background(), updatefilter, newData)
	if err != nil {

		return model.UserAccount{}, err
	}
	accountFilter := bson.D{{Key: "useremail", Value: reqresponseData.UserEmail}}
	userdata := model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), accountFilter)
	data := &model.UserAccount{}
	if err := userdata.Decode(data); err != nil {
		return model.UserAccount{}, err

	}

	account := model.UserAccount{AuthType: "Email Auth", UserId: data.UserId, UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Phone: reqresponseData.Phone}
	return account, nil
}

// Session verify Code
func SecureVerifyCode(userid string, vcode string) (model.UserAccount, error) {
	id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return model.UserAccount{}, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	response := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
	reqresponseData := &model.SignUpResponse{}
	if err := response.Decode(reqresponseData); err != nil {
		return model.UserAccount{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(reqresponseData.Vcode), []byte(vcode))
	if err != nil {
		return model.UserAccount{}, err
	}

	useraccount := model.UserAccountStoreDb{AuthType: "Secure Auth", UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, CreateAt: time.Now(), UpdateAt: time.Now(), Phone: reqresponseData.Phone}
	result, err := model.MongoInstance.Mdatabase.Collection("Account").InsertOne(context.Background(), useraccount)
	if err != nil {
		return model.UserAccount{}, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	account := model.UserAccount{AuthType: "Secure Auth", UserId: insertedID, UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Phone: reqresponseData.Phone}
	return account, nil
}

func SecureLoginSendEmail(domailValue string, requestData model.RequestModel) (model.SignUpResponse, error) {
	godotenv.Load()

	filter := bson.D{
		{Key: "useremail", Value: requestData.UserEmail},
		{Key: "username", Value: requestData.UserName},
		{Key: "authtype", Value: "Secure Auth"},
	}
	var emailExist bson.M

	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filter).Decode(&emailExist)

	if emailExist != nil {

		fromemail := os.Getenv("Email")
		apppassword := os.Getenv("AppPassword")
		host := smtphost.Getsmtphost(domailValue)
		verificationcode := generateVcode()
		htmlval := model.HtmlVar{UserName: requestData.UserName, Vcode: verificationcode}

		t, err := template.ParseFiles("html/Mailhtml.html")
		if err != nil {
			return model.SignUpResponse{}, err
		}

		buff := new(bytes.Buffer)
		if err := t.Execute(buff, htmlval); err != nil {

			return model.SignUpResponse{}, err
		}
		m := gomail.NewMessage()
		m.SetHeader("From", fromemail)
		m.SetHeader("To", requestData.UserEmail)
		m.SetHeader("Subject", "Koul Network")
		m.SetBody("text/html", buff.String())

		d := gomail.NewDialer(host, 587, fromemail, apppassword)
		//store the response data
		requestDataForresponse := &model.SignUpResponse{}
		// Send the email
		if err := d.DialAndSend(m); err != nil {
			return model.SignUpResponse{}, err
		} else {
			fmt.Println("Email Sent Successfully!")
			requestData.SendAt = time.Now().Format("1504")
			requestData.Vcode = htmlval.Vcode
			if requestData.Password != "" {
				hashpassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
				requestData.Password = string(hashpassword)
				if err != nil {
					return model.SignUpResponse{}, err
				}
			}

			hashVcode, err := bcrypt.GenerateFromPassword([]byte(verificationcode), bcrypt.DefaultCost)
			requestData.Vcode = string(hashVcode)
			if err != nil {
				return model.SignUpResponse{}, err
			}
			insertedId, err := model.MongoInstance.Mdatabase.Collection("Verification code").InsertOne(context.Background(), requestData)
			if err != nil {
				return model.SignUpResponse{}, err
			}

			filter := bson.D{{Key: "_id", Value: insertedId.InsertedID}}
			responsevalues := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
			if err := responsevalues.Decode(requestDataForresponse); err != nil {
				return model.SignUpResponse{}, err
			} else {
				requestDataForresponse.Status = "202"
			}

		}
		return *requestDataForresponse, nil
	}
	if emailExist == nil {
		return model.SignUpResponse{}, errors.New("invalid email")
	}
	return model.SignUpResponse{}, errors.New("username already exist")
}

// session login
func SecureLoginVerifyCode(userid string, vcode string) (model.UserAccount, error) {
	id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return model.UserAccount{}, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	response := model.MongoInstance.Mdatabase.Collection("Verification code").FindOne(context.Background(), filter)
	reqresponseData := &model.SignUpResponse{}
	if err := response.Decode(reqresponseData); err != nil {
		return model.UserAccount{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(reqresponseData.Vcode), []byte(vcode))
	if err != nil {
		return model.UserAccount{}, err
	}
	valuefilter := bson.D{{Key: "username", Value: reqresponseData.UserName}, {Key: "useremail", Value: reqresponseData.UserEmail}, {Key: "authtype", Value: "Secure Auth"}}
	accountvalue := model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), valuefilter)
	valueUserAccount := &model.UserAccount{}
	if err := accountvalue.Decode(valueUserAccount); err != nil {
		return model.UserAccount{}, err
	}
	account := model.UserAccount{AuthType: valueUserAccount.AuthType, UserId: valueUserAccount.UserId, UserName: valueUserAccount.UserName, UserEmail: valueUserAccount.UserEmail, Phone: valueUserAccount.Phone}
	return account, nil
}

// Login
func LoginToAccount(userEmail string, password string) (model.UserAccount, error) {
	filter := bson.D{{Key: "useremail", Value: userEmail}}
	userData := model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filter)
	useraccountobj := &model.UserAccountStoreDb{}
	if err := userData.Decode(&useraccountobj); err != nil {
		return model.UserAccount{}, errors.New("wrong email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(useraccountobj.Password), []byte(password))
	if err != nil {
		return model.UserAccount{}, errors.New("wrong email or password")
	}
	account := model.UserAccount{AuthType: useraccountobj.AuthType, UserId: useraccountobj.UserId, UserName: useraccountobj.UserName, UserEmail: useraccountobj.UserEmail, Phone: useraccountobj.Phone}

	return account, nil

}

func DeleteVcode() {
	currentTime := time.Now().Format("1504")
	converttoINT, _ := strconv.ParseInt(currentTime, 10, 64)
	filter := bson.M{}
	curser, err := model.MongoInstance.Mdatabase.Collection("Verification code").Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)

	}

	for curser.Next(context.Background()) {
		var vcodeDbObject model.RequestModel
		curser.Decode(&vcodeDbObject)

		vcodetime := vcodeDbObject.SendAt

		vcodeTimetoINT, _ := strconv.ParseInt(vcodetime, 10, 64)
		fmt.Println("Add time", converttoINT)
		fmt.Println("vcodeTime :", vcodeTimetoINT)
		if converttoINT-vcodeTimetoINT >= 2 {
			deleteFilter := bson.M{"useremail": vcodeDbObject.UserEmail}
			model.MongoInstance.Mdatabase.Collection("Verification code").DeleteOne(context.Background(), deleteFilter)
		}
	}

}

// Find the account
func UserExist(userid string) (model.UserAccount, error) {
	id, err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return model.UserAccount{}, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	response := model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filter)
	useraccountobj := &model.UserAccount{}
	if err := response.Decode(useraccountobj); err != nil {
		return model.UserAccount{}, err
	}
	account := model.UserAccount{AuthType: useraccountobj.AuthType, UserId: useraccountobj.UserId, UserName: useraccountobj.UserName, UserEmail: useraccountobj.UserEmail, Phone: useraccountobj.Phone}

	return account, nil

}

// Save Google user and also check If user already exist or not [email,name,id]
func SaveGUser(Guser model.UserAccount) (model.UserAccount, error) {
	filteremail := bson.D{{Key: "useremail", Value: Guser.UserEmail}}
	filterphone_no := bson.D{{Key: "phone", Value: Guser.Phone}}
	var emailExist bson.M
	var phoneExist bson.M

	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filteremail).Decode(&emailExist)
	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filterphone_no).Decode(&phoneExist)
	if emailExist == nil && phoneExist == nil {

		G_account := model.UserAccountStoreDb{
			UserName:  Guser.UserName,
			UserEmail: Guser.UserEmail,
			Phone:     Guser.Phone,
			AuthType:  "Google Auth",
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
		}

		insertedrecord, err := model.MongoInstance.Mdatabase.Collection("Account").InsertOne(context.Background(), G_account)
		if err != nil {
			return model.UserAccount{}, errors.New("unable to save Google user")
		}
		insertedid := insertedrecord.InsertedID.(primitive.ObjectID).Hex()

		return model.UserAccount{AuthType: G_account.AuthType, UserId: insertedid, UserName: G_account.UserName, UserEmail: G_account.UserEmail, Phone: G_account.Phone}, nil
	}

	if phoneExist != nil && emailExist != nil {
		return model.UserAccount{}, errors.New("email and phone no. already in use")

	} else if phoneExist != nil {
		return model.UserAccount{}, errors.New("phone no. already in use")

	}
	return model.UserAccount{}, errors.New("email already in use")
}

func GoogleUserExist(email string) (model.UserAccount, error) {
	filter := bson.D{
		{Key: "useremail", Value: email},
		{Key: "authtype", Value: "Google Auth"},
	}

	useraccountobj := &model.UserAccount{}

	response := model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filter)

	if err := response.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.UserAccount{}, errors.New("no account found with this email and authtype")
		}
		return model.UserAccount{}, errors.New("error finding the user")
	}

	if err := response.Decode(useraccountobj); err != nil {
		return model.UserAccount{}, errors.New("error decoding the user")
	}

	// Return the found user
	account := model.UserAccount{
		AuthType:  useraccountobj.AuthType,
		UserId:    useraccountobj.UserId,
		UserName:  useraccountobj.UserName,
		UserEmail: useraccountobj.UserEmail,
		Phone:     useraccountobj.Phone,
	}

	return account, nil
}
