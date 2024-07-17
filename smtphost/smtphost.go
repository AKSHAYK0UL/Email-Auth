package smtphost

import "strings"

//List of  Domain names
var listofdomains = [10]string{"@gmail.com", "@yahoo.com", "@outlook.com", "@aol.com", "@icloud.com", "@protonmail.com", "@zoho.com", "@yandex.com", "@ya.ru", "@yandex.ru"}

//List of Smtp server for different Domains
var listofsmtp = map[string]string{"@gmail.com": "smtp.gmail.com", "@yahoo.com": "smtp.mail.yahoo.com", "@outlook.com": "smtp.office365.com", "@aol.com": "smtp.aol.com", "@icloud.com": "smtp.mail.me.com", "@protonmail.com": "127.0.0.1", "@zoho.com": "smtp.zoho.com", "@yandex.com": "smtp.yandex.com", "@yandex.ru": "smtp.yandex.ru", "@ya.ru": "smtp.yandex.ru"}

//Take input as a Domain like gmail.com and return smtp.gmail.com
//it returns the smtp host
func Checkthedomain(useremail string) string {
	for _, domainvalue := range listofdomains {
		if strings.Contains(useremail, domainvalue) {
			return domainvalue
		}
	}
	return "try another email provider"
}
func Getsmtphost(domain string) string {
	if smtpvalue := listofsmtp[domain]; smtpvalue != "" {
		return smtpvalue
	}
	return "try another email provider"
}
