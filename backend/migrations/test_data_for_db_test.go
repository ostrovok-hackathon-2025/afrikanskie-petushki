package migrations

import (
	"fmt"
	"testing"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
)

func TestInitDB(t *testing.T) {
	var (
		//user doverlof
		doverlofPassword = "Doverlof12345!"
		//user notblinkyet
		notblinkyetPassword = "Notblinkyet12345!"
		//user smokingElk
		smokingElkPassword = "smokingElk12345!"
		//user sophistik
		sophistikPassword = "sophistik12345!"
		//user chicherin
		chicherinPassword = "Chicherin12345!"
		//user root
		rootPassword = "Root12345!"
	)
	fmt.Println(pkg.HashPassword(doverlofPassword))
	//Got 4b3a8ff5bfd431b2fe755180f65a49c54585001ced55577266e2791856000f5f
	fmt.Println(pkg.HashPassword(notblinkyetPassword))
	//Got 81c3f6c564f43e6a7ae2d22689ea5b254d63582460483db5807b42bba883dd9e
	fmt.Println(pkg.HashPassword(smokingElkPassword))
	//Got 19ba7759f8434288a30f4b433a7b6fc169c9e818467f31ccbc697cc5f6439794
	fmt.Println(pkg.HashPassword(sophistikPassword))
	//Got ffbc6e1ecd0e383f8950e6e0ff7b1226edc1de65c22154b8278fb6b8a8823de0
	fmt.Println(pkg.HashPassword(chicherinPassword))
	//Got be75f8c0378700bb1f4c1c9dbc64d68fd1a65a4c81cb7f469998200acfc691e2
	fmt.Println(pkg.HashPassword(rootPassword))
	//Got ee1a7dc746c6024bd64fe2e2245e6566fbff2d27f617284a002ec1202734d3c7
}
