package utils

import (
	"fmt"
	"github.com/leekchan/accounting"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// CapitalizeWords Format text first letters capitalized
func CapitalizeWords(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}
	return strings.Join(words, " ")
}

func ConvertToTime(value interface{}) (*time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return &v, nil
	case []uint8: // Se for uma string no banco de dados
		str := string(v)
		parsedTime, err := time.Parse("15:04:05", str) // Ajuste o formato conforme necessário
		if err != nil {
			return nil, err
		}
		return &parsedTime, nil
	default:
		return nil, fmt.Errorf("unsupported type %T for time conversion", value)
	}
}

func FormatCPF(cpf string) string {
	if len(cpf) != 11 {
		return cpf
	}
	return fmt.Sprintf("%s.%s.%s-%s", cpf[:3], cpf[3:6], cpf[6:9], cpf[9:])
}

var stopWords = []string{",", ".", ";", ":", "/", "|", "-", "_", "!", "'", "@", "#", "$", "%", "&", "*", "de", "a", "o", "e", "que", "do", "da", "em", "um", "para", "com", "não", "uma", "os", "no", "se", "na", "por", "mais", "as", "dos", "como", "mas", "foi", "ao", "ele", "das", "tem", "à", "seu", "sua", "ou", "ser", "quando", "muito", "há", "nos", "já", "está", "eu", "também", "só", "pelo", "pela", "até", "isso", "ela", "entre", "era", "depois", "sem", "mesmo", "aos", "ter", "seus", "quem", "nas", "me", "esse", "eles", "estão", "você", "tinha", "foram", "essa", "num", "nem", "suas", "meu", "às", "minha", "têm", "numa", "pelos", "elas", "havia", "seja", "qual", "será", "nós", "tenho", "lhe", "deles", "essas", "esses", "pelas", "este", "fosse", "dele", "tu", "te", "vocês", "vos", "lhes", "meus", "minhas", "teu", "tua", "teus", "tuas", "nosso", "nossa", "nossos", "nossas", "dela", "delas", "esta", "estes", "estas", "aquele", "aquela", "aqueles", "aquelas", "isto", "aquilo"}

func RemoveStopWords(input string) string {
	words := strings.Fields(input)
	// Se houver apenas uma palavra, retorna a entrada original
	if len(words) == 1 {
		return input
	}

	var result []string
	for _, word := range words {
		skip := false
		for _, stopWord := range stopWords {
			if strings.EqualFold(word, stopWord) {
				skip = true
				break
			}
		}
		if !skip {
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

func FormatDecoin(decoins int64) string {
	var decoinConfig = accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}

	return decoinConfig.FormatMoney(decoins)
}

func DifferenceSlices(slice1 []int64, slice2 []int64) []int64 {
	var diff []int64
	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}

// ParseDateTimeFromQueryParam Try parse an optional parameter from url to datetime, if it can't try parse to datetime
// will try parse to date and return an error if it can't
// In case you don't receive the optional parameter, return the real date and time
func ParseDateTimeFromQueryParam(r *http.Request, key string) (*time.Time, error) {
	date := r.URL.Query().Get(key)

	var parsedDate time.Time
	var err error

	if date != "" {
		parsedDate, err = time.Parse("20060102150405", date)
		if err != nil {
			parsedDate, err = time.Parse("20060102", date)
		}
		if err != nil {
			return nil, err
		}
	} else {
		return NewTimeNowUTC(), nil
	}
	return &parsedDate, nil
}

// RemoveAccents String with accents count as two when use len
// To count quantity of characters we need replace chars with accents
// The following character is counting as 2 characters: [¨,°].
func RemoveAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	return output, nil
}

// ToPointer used to return the pointer of any type of value
func ToPointer[T any](v T, e error) (*T, error) {
	return &v, e
}

// Character's blocked -[".","[","]","|","(",")","*","$","?","+","\","{","}"]
const cleanRegexp = "\\{.*\\}\\..*\\[.*].*\\|\\(\\)\\*.*\\$\\?\\+\\\\"

// CleanMySQLRegexp Will add a double slash before each special character to when we use on mysql query doesn't scape
func CleanMySQLRegexp(regexp string) string {
	var newRegex string
	for _, it := range regexp {
		if strings.ContainsRune(cleanRegexp, it) {
			newRegex += "\\" + string(it)
		} else {
			newRegex += string(it)
		}
	}

	return newRegex
}

func FormatHourMinuteSecond(timeToFormat int64) string {
	hours := math.Floor(float64(timeToFormat) / 60 / 60)
	seconds := timeToFormat % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = timeToFormat % 60

	//15m
	if hours == 0 && seconds == 0 {
		return fmt.Sprintf("%.0fm", minutes)
	}

	//1h 15m
	if hours > 0 && seconds > 0 {
		return fmt.Sprintf("%.0fh %.0fm", hours, minutes)
	}

	//2h 30s
	if hours > 0 && minutes == 0 && seconds > 0 {
		return fmt.Sprintf("%.0fh %ds", hours, seconds)
	}

	//1h
	if hours > 0 && minutes == 0 && seconds == 0 {
		return fmt.Sprintf("%.0fh", hours)
	}

	//30s
	if hours == 0 && minutes == 0 && seconds > 0 {
		return fmt.Sprintf("%ds", seconds)
	}

	//15m 30s
	if hours == 0 {
		return fmt.Sprintf("%.0fm %ds", minutes, seconds)
	}

	//1h 15m 30s
	return fmt.Sprintf("%.0fh %.0fm %ds", hours, minutes, seconds)
}

func IndexToLetter(index int) string {
	return string(rune('A' + index))
}

// GetPublicRoutes return all routes that can be accessed by users without login.

func AccentRegex(term string) string {
	replacements := map[string]string{
		"a": "[aáàâã]",
		"á": "[aáàâã]",
		"à": "[aáàâã]",
		"â": "[aáàâã]",
		"ã": "[aáàâã]",
		"e": "[eéèê]",
		"é": "[eéèê]",
		"è": "[eéèê]",
		"ê": "[eéèê]",
		"i": "[iíìî]",
		"í": "[iíìî]",
		"ì": "[iíìî]",
		"î": "[iíìî]",
		"o": "[oóòôõ]",
		"ó": "[oóòôõ]",
		"ò": "[oóòôõ]",
		"ô": "[oóòôõ]",
		"õ": "[oóòôõ]",
		"u": "[uúùû]",
		"ú": "[uúùû]",
		"ù": "[uúùû]",
		"û": "[uúùû]",
		"c": "[cç]",
		"ç": "[cç]",
	}
	term = strings.ToLower(term)

	for k, v := range replacements {
		term = regexp.MustCompile(k).ReplaceAllString(term, v)
	}

	return term
}
