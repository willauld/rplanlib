package rplanlib

//import "github.com/pelletier/go-toml"
import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml"
)

var iamKeys = []string{"primary", "age", "retire", "through", "definedContributionPlan"}
var iamReqKeys = []string{"age", "retire", "through"}

var socialSecurityKeys = []string{"amount", "age", "FRA"}
var socialSecurityReqKeys = socialSecurityKeys

var incomeKeys = []string{"amount", "age", "inflation", "tax"}
var incomeReqKeys = []string{"amount", "age"}

var expenseKeys = []string{"amount", "age", "inflation"}
var expenseReqKeys = []string{"amount", "age"}

var assetKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell", "primaryResidence", "rate", "brokerageRate"}
var assetReqKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell"}

var iRAKeys = []string{"bal", "rate", "contrib", "inflation", "period"}
var iRAReqKeys = []string{"bal"}

var rothKeys = iRAKeys
var rothReqKeys = iRAReqKeys

var aftertaxKeys = []string{"bal", "basis", "rate", "contrib", "inflation", "period"}
var aftertaxReqKeys = iRAReqKeys

var minKeys = []string{"amount"}
var minReqKeys = minKeys

var maxKeys = minKeys
var maxReqKeys = minKeys

// TomlStrDefs works with InputStrDefs to Map Toml information to
// rplanlib API static portion
//
// List of toml paths to user supplied information.
// Includes editing codes for correcting as needed at runtime.
// '@' means the corresponding key in InputStrDefs should be used to set
//		iam key names
// '%n' means the nth iam key name is to be substituted here
// '%i' means the current stream key name is to be substituted here
// '#n' means that the resulting string is a range string and the nth
//		number from the range should be extracted and used for assignment
// After converting the toml paths there should not be any of these codes
// remaining.
// The resulting path and the corresponding key from InputStrDefs are used
// together to retreave and set the value in the Input String map.
var TomlStrDefs = []string{ // each entry corresponds to InputStrDefs entries
	"title",
	"retirement_type",
	"@iam.key1",
	"@iam.key2",
	"iam.%0.age",
	"iam.%1.age",
	"iam.%0.retire",
	"iam.%1.retire",
	"iam.%0.through",
	"iam.%1.through",
	"iam.%0.definedContributionPlan#1",
	"iam.%1.definedContributionPlan#1",
	"iam.%0.definedContributionPlan#2",
	"iam.%1.definedContributionPlan#2",
	"SocialSecurity.%0.amount",
	"SocialSecurity.%1.amount",
	"SocialSecurity.%0.age#1",
	"SocialSecurity.%1.age#1",
	"IRA.%0.bal",
	"IRA.%1.bal",
	"IRA.%0.rate",
	"IRA.%1.rate",
	"IRA.%0.contrib",
	"IRA.%1.contrib",
	"IRA.%0.period#1", //ContribStartAge1,
	"IRA.%1.period#1", //ContribStartAge2,
	"IRA.%0.period#2", //ContribEndAge1,
	"IRA.%1.period#2", //ContribEndAge2,
	"IRA.%0.inflation",
	"IRA.%1.inflation",
	"roth.%0.bal",
	"roth.%1.bal",
	"roth.%0.rate",
	"roth.%1.rate",
	"roth.%0.contrib",
	"roth.%1.contrib",
	"roth.%0.period#1", //contribStartAge1
	"roth.%1.period#1", //contribStartAge2,
	"roth.%0.period#2", //contribEndAge1,
	"roth.%1.period#2", // contribEndAge2
	"roth.%0.inflation",
	"roth.%1.inflation",
	"aftertax.bal",
	"aftertax.basis",
	"aftertax.rate",
	"aftertax.contrib",
	"aftertax.period#1", //contribStartAge
	"aftertax.period#2", //contribEndAge,
	"aftertax.inflation",

	"min.income.amount",
	"max.income.amount",

	"inflation",
	"returns",
	"maximize",
	"dollarsInThousands",
}

// TomlStreamStrDefs works with InputStreamStrDefs to Map Toml information to
// rplanlib API dynamic portion (per stream portion)
var TomlStreamStrDefs = []string{ // each entry corresponds to InputStreamStrDefs entries
	"@income",
	"income.%i.amount",
	"income.%i.age#1",
	"income.%i.age#2",
	"income.%i.inflation",
	"income.%i.tax",
	"@expense",
	"expense.%i.amount",
	"expense.%i.age#1",
	"expense.%i.age#2",
	"expense.%i.inflation",
	"expense.%i.tax",
	"@asset",
	"asset.%i.value",
	"asset.%i.ageToSell",
	"asset.%i.costAndImprovements",
	"asset.%i.owedAtAgeToSell",
	"asset.%i.primaryResidence",
	"asset.%i.rate",
	"asset.%i.brokerageRate",
}

// Required to have iam, retirement_type and at least one of IRA, roth or Aftertax
func keyIn(k string, keys []string) bool {
	for _, key := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func categoryMatchAndUnknowns(parent string, keys []string) (bool, []string) {
	var reqKeys []string
	var akeys []string
	unknownPaths := []string{}
	switch parent {
	case "retirement_type": // TODO FIXME change retirement_type to filingStatus
		fallthrough
	case "returns":
		fallthrough
	case "inflation":
		fallthrough
	case "maximize":
		return false, unknownPaths
	case "iam":
		reqKeys = iamReqKeys
		akeys = iamKeys
	case "SocialSecurity":
		reqKeys = socialSecurityReqKeys
		akeys = socialSecurityKeys
	case "IRA":
		reqKeys = iRAReqKeys
		akeys = iRAKeys
	case "roth":
		reqKeys = rothReqKeys
		akeys = rothKeys
	case "aftertax":
		reqKeys = aftertaxReqKeys
		akeys = aftertaxKeys
	case "income":
		reqKeys = incomeReqKeys
		akeys = incomeKeys
	case "asset":
		reqKeys = assetReqKeys
		akeys = assetKeys
	case "expense":
		reqKeys = expenseReqKeys
		akeys = expenseKeys
	case "min":
		reqKeys = minReqKeys
		akeys = minKeys
	case "max":
		reqKeys = maxReqKeys
		akeys = maxKeys
	}
	for _, k := range reqKeys {
		if !keyIn(k, keys) {
			return false, unknownPaths
		}
	}
	for _, k := range keys {
		if !keyIn(k, akeys) {
			//uk := parent + "." + k
			unknownPaths = append(unknownPaths, k)
		}
	}
	return true, unknownPaths
}

func checkNames(golden []string, toevaluate []string) error {
	// toevaluate should be a subset of golden (possibly empty)
	if len(golden) < len(toevaluate) {
		e := fmt.Errorf("checkNames: defined names are (%v) which does not include all of (%v)", golden, toevaluate)
		return e
	}
	for _, v := range toevaluate {
		if !keyIn(v, golden) {
			e := fmt.Errorf("checkNames: name '%v' is not a member of (%v)", v, golden)
			return e
		}
	}
	return nil
}

func getKeys(path string, config *toml.Tree) []string {
	if config.Has(path) {
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
		matched, unknown := categoryMatchAndUnknowns(path, keys)
		if !matched || len(unknown) > 0 {
			if len(unknown) > 0 {
				//fmt.Printf("unknown key list: %#v\n", unknown)
				keys = []string{"nokey"}
				for i := 0; i < len(unknown); i++ {
					keys = append(keys, unknown[i])
				}
				//fmt.Printf("New key list: %#v\n", keys)
			}
			// These should be the unknown 'name' values
			//Need to find the order from within ie which one if Primary
			if path == "iam" {
				L := len(keys)
				if L > 2 {
					//Can only have one or two
					fmt.Printf("TOO MANY IAM NAMES: %#v\n", keys)
					os.Exit(0)
				}
				if L == 1 {
					return keys
				}

				prime := -1
				for i, v := range keys {
					lPath := path + "." + v + "." + "primary"
					if v == "nokey" {
						lPath = path + "." + "primary"
					}
					if config.Has(lPath) {
						lPathobj := config.Get(lPath)
						p := lPathobj.(bool)
						if p == true {
							if prime < 0 {
								prime = i
							} else {
								fmt.Printf("Error Only one 'iam' can be prime\n")
								os.Exit(0)
							}
						}
					}
				}
				if prime < 0 {
					fmt.Printf("Error At least one 'iam' must be prime\n")
					os.Exit(0)
				}
				//lkey[0] = keys[prime]
				//lkey[1] = keys[1]
				if prime == 1 {
					return []string{keys[1], keys[0]}
				}
				// return keys fall through to return keys
			}
			return keys
		}
		return []string{"nokey"}
	}
	return []string{}
}

func getPathStrValue(path string, config *toml.Tree) string {
	var targetVal interface{}
	obj := config.Get(path)
	switch reflect.TypeOf(obj).Name() {
	case "int64":
		targetVal = obj.(int64)
	case "float64":
		targetVal = obj.(float64)
	case "bool":
		targetVal = obj.(bool)
	case "string":
		targetVal = obj.(string)
	}
	s := fmt.Sprintf("%v", targetVal)
	return s
}

func setStringMapValueWithValue(ipsm *map[string]string,
	s string, val string) error {
	//fmt.Printf("Set %s to %s\n", s, val)
	_, ok := (*ipsm)[s]
	if !ok {
		e := fmt.Errorf("setStringMapValue: attempt to set a non-existant parameter: %s", s)
		return e
	}
	(*ipsm)[s] = val
	return nil
}

func setStringMapValue(ipsm *map[string]string,
	s string, path string, config *toml.Tree) error {
	//fmt.Printf("Attempting to set %s\n", s)
	_, ok := (*ipsm)[s]
	if !ok {
		e := fmt.Errorf("setStringMapValue: attempt to set a non-existant parameter: %s", s)
		return e
	}
	//fmt.Printf("checking path: %s\n", path)
	Hval := -1
	indxH := strings.Index(path, "#")
	if indxH >= 0 {
		Hval = int(path[indxH+1] - '0')
		path = path[:indxH]
		//fmt.Printf("now checking path: %s\n", path)
	}
	if config.Has(path) {
		v := getPathStrValue(path, config)
		//fmt.Printf("DOES HAVE and it is: %s\n", v)
		if Hval > 0 {
			svals := strings.Split(v, "-")
			v = svals[Hval-1]
			//fmt.Printf("DOES HAVE and will use: %s\n", v)
		}
		(*ipsm)[s] = v
	} /*else {
		fmt.Printf("DOES not have: %s\n", path)
	}*/
	return nil
}

func GetInputStringsMapFromToml(filename string) (*map[string]string, error) {
	config, err := toml.LoadFile(filename)
	if err != nil {
		e := fmt.Errorf("Error: %s", err)
		return nil, e
	}

	ipsm := NewInputStringsMap()

	// Need to ensure that ONE iam is primary if there are two

	//
	// Get all the unknown keys first
	//
	iamNames := getKeys("iam", config)
	//fmt.Printf("iam names: %#v\n", iamNames)

	ssNames := getKeys("SocialSecurity", config)
	//fmt.Printf("SS names: %#v\n", ssNames)
	err = checkNames(iamNames, ssNames)
	if err != nil {
		e := fmt.Errorf("getInputStringMapFromToml: %s, missing identifiers must be defined in 'iam' section", err)
		return nil, e
	}
	iRANames := getKeys("IRA", config)
	//fmt.Printf("IRA names: %#v\n", iRANames)
	err = checkNames(iamNames, iRANames)
	if err != nil {
		e := fmt.Errorf("getInputStringMapFromToml: %s, missing identifiers must be defined in 'iam' section", err)
		return nil, e
	}
	rothNames := getKeys("roth", config)
	//fmt.Printf("roth names: %#v\n", rothNames)
	err = checkNames(iamNames, rothNames)
	if err != nil {
		e := fmt.Errorf("getInputStringMapFromToml: %s, missing identifiers must be defined in 'iam' section", err)
		return nil, e
	}
	//aftertaxNames := getKeys("aftertax", config)
	//fmt.Printf("aftertax names: %#v\n", aftertaxNames)
	assetNames := getKeys("asset", config)
	//fmt.Printf("Asset names: %#v\n", assetNames)
	incomeNames := getKeys("income", config)
	//fmt.Printf("Income names: %#v\n", incomeNames)
	expenseNames := getKeys("expense", config)
	//fmt.Printf("Expense names: %#v\n", expenseNames)

	//
	// Now we can work our way though setting values in InputStrDefs
	//
	for i, k := range TomlStrDefs {
		if k[0] == '@' {
			// All keys used in TomlStrDefs are iam keys (names)
			indx := int(k[len(k)-1] - '0')
			if len(iamNames) < indx {
				continue
			}
			n := iamNames[indx-1]
			err = setStringMapValueWithValue(&ipsm, InputStrDefs[i], n)
			if err != nil {
				fmt.Printf("getInputStringsMapFromToml: %s\n", err)
			}
			continue
		}
		indxP := strings.Index(k, "%")
		if indxP < 0 {
			err = setStringMapValue(&ipsm, InputStrDefs[i], k, config)
			if err != nil {
				fmt.Printf("getInputStringsMapFromToml: %s\n", err)
			}
			continue
		}
		val := int(k[indxP+1] - '0')
		//fmt.Printf("Index val is %d\n", val)
		var p string
		if len(iamNames) <= val {
			continue
		}
		if iamNames[val] != "nokey" {
			//fmt.Printf("*** iamNames[%d]: %s\n", val, iamNames[val])
			p = strings.Replace(k, "%0", iamNames[val], 1)
			if val == 1 {
				p = strings.Replace(k, "%1", iamNames[val], 1)
			}
		} else {
			// iamNames is "nokey"
			s := strings.Split(k, ".")
			p = s[0] + "." + s[2]
			//fmt.Printf("have nokey using: %s\n", p)
		}
		//fmt.Printf("Will use path: %s\n", p)

		err = setStringMapValue(&ipsm, InputStrDefs[i], p, config)
		if err != nil {
			fmt.Printf("getInputStringsMapFromToml: %s\n", err)
		}
		continue
	}
	//
	// Now we can work our way though setting values in InputStreamStrDefs
	//
	var names []string
	for j := 1; j < MaxStreams+1; j++ {
		for i, k := range TomlStreamStrDefs {
			//fmt.Printf("InputStrDefs[%d]: '%s', TomlStreamStrDefs[%d]: '%s'\n", i, rplanlib.InputStreamStrDefs[i], i, k)
			targetStr := fmt.Sprintf("%s%d", InputStreamStrDefs[i], j)
			if k[0] == '@' {
				switch k[1:] {
				case "income":
					names = incomeNames
				case "expense":
					names = expenseNames
				case "asset":
					names = assetNames
				default:
					fmt.Printf("EEEEEError 232323\n")
				}
				if len(names) < j {
					continue
				}
				n := names[j-1]
				err = setStringMapValueWithValue(&ipsm, targetStr, n)
				if err != nil {
					fmt.Printf("getInputStringsMapFromToml: %s\n", err)
				}
				continue
			} else { // TODO should not need else with the above continue
				indxP := strings.Index(k, "%")
				if indxP < 0 {
					err = setStringMapValue(&ipsm, targetStr, k, config)
					if err != nil {
						fmt.Printf("getInputStringsMapFromToml: %s\n", err)
					}
					continue
				}
				strs := strings.Split(k, ".")
				switch strs[0] {
				case "income":
					names = incomeNames
				case "expense":
					names = expenseNames
				case "asset":
					names = assetNames
				default:
					fmt.Printf("EEEEEError 8989\n")
				}
				if len(names) < j {
					continue
				}
				n := names[j-1]
				var p string
				//fmt.Printf("*** names[%d]: %s\n", j-1, n)
				if n != "nokey" {
					p = strings.Replace(k, "%i", n, 1)
				} else {
					p = strs[0] + "." + strs[2]
					//fmt.Printf("have nokey using: %s\n", p)
				}
				//fmt.Printf("Will use path: %s\n", p)

				err = setStringMapValue(&ipsm, targetStr, p, config)
				if err != nil {
					fmt.Printf("getInputStringsMapFromToml: %s\n", err)
				}
				continue
			}
		}
	}
	// Toml file does NOT have dollars listed in thousands $(000)
	err = setStringMapValueWithValue(&ipsm, "dollarsInThousands", "false")
	if err != nil {
		fmt.Printf("getInputStringsMapFromToml: %s\n", err)
	}
	return &ipsm, nil
}
