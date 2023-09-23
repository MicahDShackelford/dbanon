package dbanon

import (
	"regexp"
	"strings"
	"math/rand"
	"syreclabs.com/go/faker"
	"github.com/google/uuid"
)

type Provider struct {}

func NewProvider() *Provider {
	p := &Provider{}

	return p
}

type ProviderInterface interface {
	Get(fakeType string, currentValue *string) string
}

func (p Provider) Get(fakeType string, currentValue *string) string {
	if (strings.HasPrefix(fakeType, "dynamic")) {
		return p.processDynamicFake(fakeType, currentValue)
	}

	switch fakeType {
		// Name
		case "first_name":
			return faker.Name().FirstName()
		case "last_name":
			return faker.Name().LastName()
		case "full_name":
			return faker.Name().FirstName() + " " + faker.Name().LastName()
		
		// Company
		case "company_name":
			return faker.Company().Name()

		// Internet
		case "username":
			return faker.Internet().UserName()
		case "password":
			return faker.Internet().Password(8, 14)
		case "ipv4":
			return faker.Internet().IpV4Address()
		case "url":
			return faker.Internet().Url()
		case "linkedin_url":
			return "https://linkedin.com/in/" + faker.Internet().UserName()
		case "md5":
			return faker.Lorem().Characters(32)
		case "uuid":
			return uuid.NewString()
		case "json":
			return p.fakeJson()
		case "query_params":
			return p.fakeQueryParams()

		// Dates
		case "datetime":
			return faker.Date().Birthday(0, 40).Format("2006-01-02 15:04:05")

		// Geo
		case "state":
			return faker.Address().State()
		case "city":
			return faker.Address().City()
		case "postcode":
			return faker.Address().Postcode()
		case "street":
			return faker.Address().StreetAddress()
		case "country_code":
			return faker.Address().CountryCode()

		// Currency
		case "money":
			numberOfDigits := rand.Intn(5 - 1) + 1
			return faker.Number().Number(numberOfDigits)
		case "money_decimal":
			numberOfDigits := rand.Intn(5 - 1) + 1
			return faker.Number().Decimal(numberOfDigits, 2)

		// Contact
		case "email":
			return faker.Internet().Email()
		case "telephone":
			return faker.PhoneNumber().PhoneNumber()
	}

	logger := GetLogger()
	logger.Error(fakeType + " does not match any known type")

	return ""
}

func (p Provider) fakeJson() string {
	return "{\"fake\": true}"
}

func (p Provider) fakeQueryParams() string {
	return "?fake=true"
}

func (p Provider) processDynamicFake(fakeType string, currentValue *string) string {
	// strip dynamic from fakeType
	fakeType = strings.ReplaceAll(fakeType, "dynamic.", "")

	var args []string
	args, fakeType = p.processDynamicArgs(fakeType)

	switch fakeType {
	case "email":
		return p.dynamicEmail(currentValue, args)
	}

	logger := GetLogger()
	logger.Error(fakeType + " does not match any known dynamic type")

	return ""
}

func (p Provider) processDynamicArgs(fakeType string) ([]string, string)  {
	reg := regexp.MustCompile(`\((.*?)\)`)
	args := string(reg.Find([]byte(fakeType)))

	fakeType = strings.ReplaceAll(fakeType, args, "")

	args = strings.TrimPrefix(args, "(")
	args = strings.TrimSuffix(args, ")")

	return strings.Split(args, ","), fakeType
}

func (p Provider) dynamicEmail(currentValue *string, args []string) string {
	if len(args) < 1 || args[0] == "" {
		return faker.Internet().Email()
	}

	// args[0] is email we wish to not fake
	ignoreEmail := args[0]
	val := *currentValue
	if strings.HasSuffix(val, ignoreEmail) {
		return val
	}

	return faker.Internet().Email()
}