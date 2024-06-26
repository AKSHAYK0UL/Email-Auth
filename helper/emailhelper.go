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
	var emailExist bson.M
	var nameExist bson.M
	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filteremail).Decode(&emailExist)
	model.MongoInstance.Mdatabase.Collection("Account").FindOne(context.Background(), filtername).Decode(&nameExist)
	if emailExist == nil && nameExist == nil {

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

		}
		return *requestDataForresponse, nil
	}
	if emailExist != nil && nameExist != nil {
		return model.SignUpResponse{}, errors.New("email and username already exist")
	} else if emailExist != nil {
		return model.SignUpResponse{}, errors.New("email already exist")
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

	useraccount := model.UserAccountStoreDb{UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Password: reqresponseData.Password, CreateAt: time.Now(), UpdateAt: time.Now()}
	result, err := model.MongoInstance.Mdatabase.Collection("Account").InsertOne(context.Background(), useraccount)
	if err != nil {
		return model.UserAccount{}, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()

	account := model.UserAccount{UserId: insertedID, UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail}
	return account, nil
}

// reset password
func ResetpasswordSendEmail(domailValue string, requestData model.RequestModel) (model.SignUpResponse, error) {
	filter := bson.D{{Key: "useremail", Value: requestData.UserEmail}}

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

	useraccount := model.UserAccountStoreDb{UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail, Password: reqresponseData.Password, UpdateAt: time.Now()}
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

	account := model.UserAccount{UserId: data.UserId, UserName: reqresponseData.UserName, UserEmail: reqresponseData.UserEmail}
	return account, nil
}

//Login

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
	account := model.UserAccount{UserId: useraccountobj.UserId, UserName: useraccountobj.UserName, UserEmail: useraccountobj.UserEmail}

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
	account := model.UserAccount{UserId: useraccountobj.UserId, UserName: useraccountobj.UserName, UserEmail: useraccountobj.UserEmail}

	return account, nil

}
