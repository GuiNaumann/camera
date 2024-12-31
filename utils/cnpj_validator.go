package utils

import (
	"camera/infrastructure/modules/impl/http_error"
	"strconv"
	"strings"
)

func IsCNPJValid(digits string) (bool, error) {
	return valid(digits)
}

func sanitize(date string) string {
	date = strings.Replace(date, ".", "", -1)
	date = strings.Replace(date, "-", "", -1)
	date = strings.Replace(date, "/", "", -1)
	return date
}

func valid(date string) (bool, error) {
	date = sanitize(date)

	if len(date) != 14 {
		return false, http_error.NewBadRequestError(http_error.CNPJInvalid)
	}

	if strings.Contains(blacklist, date) || !check(date) {
		return false, http_error.NewBadRequestError(http_error.CNPJInvalid)
	}

	return true, nil
}

const blacklist = `00000000000000
11111111111111
22222222222222
33333333333333
44444444444444
55555555555555
66666666666666
77777777777777
88888888888888
99999999999999`

func stringToIntSlice(date string) (res []int) {
	for _, d := range date {
		x, err := strconv.Atoi(string(d))
		if err != nil {
			continue
		}
		res = append(res, x)
	}
	return
}

func check(date string) bool {
	return verify(stringToIntSlice(date), 5, 12) && verify(stringToIntSlice(date), 6, 13)
}

func verify(date []int, j int, n int) bool {

	sum := 0

	for i := 0; i < n; i++ {
		v := date[i]
		sum += v * j

		if j == 2 {
			j = 9
		} else {
			j -= 1
		}
	}

	rest := sum % 11

	v := date[n]
	x := 0

	if rest >= 2 {
		x = 11 - rest
	}

	if v != x {
		return false
	}

	return true
}
